package dbcore

import (
	"fmt"
	"sort"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MigrationFunc func(db *gorm.DB) error

// Migration представляет собой миграцию сервиса
type Migration struct {
	Version     int
	Description string
	Up          func(db *gorm.DB) error
	Down        func(db *gorm.DB) error
}

func Migrations(migrations ...Migration) []Migration {
	return migrations
}

// MigrationRecord представляет запись о примененной миграции в базе данных
type MigrationRecord struct {
	Version int  `gorm:"primaryKey"`
	Applied bool `gorm:"column:applied;default:false"`
}

// Migrator управляет миграциями для конкретного сервиса
type Migrator struct {
	DB          *gorm.DB
	Logger      *zap.Logger
	ServiceName string
	Migrations  []Migration
}

// New создает новый менеджер миграций для сервиса
func NewMigrator(db *gorm.DB, logger *zap.Logger, serviceName string, migrations []Migration) *Migrator {
	return &Migrator{
		DB:          db,
		Logger:      logger,
		ServiceName: serviceName,
		Migrations:  migrations,
	}
}

func (m *Migrator) TableName() string {
	return fmt.Sprintf("%s_migrations", m.ServiceName)
}

func (m *Migrator) AddMigration(description string, up MigrationFunc, down MigrationFunc) {
	m.Migrations = append(m.Migrations, Migration{
		Version:     len(m.Migrations) + 1,
		Description: description,
		Up:          up,
		Down:        down,
	})
}

func (m *Migrator) AddUpMigration(description string, up MigrationFunc) {
	m.AddMigration(description, up, m.EmptyFunc)
}

func (m *Migrator) AddDownMigration(description string, down MigrationFunc) {
	m.AddMigration(description, m.EmptyFunc, down)
}

// Run запускает все миграции для сервиса
func (m *Migrator) Run() error {
	err := m.DB.Exec(`
		CREATE TABLE IF NOT EXISTS ` + m.TableName() + ` (
			version INT PRIMARY KEY,
			applied BOOLEAN NOT NULL DEFAULT FALSE
		)
	`).Error
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Получаем все применённые миграции
	var appliedMigrations []MigrationRecord
	if err := m.DB.Table(m.TableName()).Find(&appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Создаем map для быстрой проверки
	appliedVersions := make(map[int]bool)
	for _, record := range appliedMigrations {
		appliedVersions[record.Version] = record.Applied
	}

	// Сортируем миграции по версии
	sort.Slice(m.Migrations, func(i, j int) bool {
		return m.Migrations[i].Version < m.Migrations[j].Version
	})

	// Применяем миграции
	for _, migration := range m.Migrations {
		// Если миграция уже применена, пропускаем
		if applied, exists := appliedVersions[migration.Version]; exists && applied {
			m.Logger.Debug("Migration already applied",
				zap.String("service", m.ServiceName),
				zap.Int("version", migration.Version),
				zap.Bool("applied", applied))
			continue
		}

		m.Logger.Info("Applying migration",
			zap.String("service", m.ServiceName),
			zap.Int("version", migration.Version),
			zap.String("description", migration.Description))

		// Применяем миграцию в транзакции
		err := m.DB.Transaction(func(tx *gorm.DB) error {
			// Запускаем Up миграцию
			if err := migration.Up(tx); err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
			}

			// Записываем в таблицу миграций
			return tx.Table(m.TableName()).Create(&MigrationRecord{
				Version: migration.Version,
				Applied: true,
			}).Error
		})

		if err != nil {
			return err
		}

		m.Logger.Info("Migration applied successfully",
			zap.String("service", m.ServiceName),
			zap.Int("version", migration.Version))
	}

	return nil
}

func (m *Migrator) EmptyFunc(db *gorm.DB) error {
	return nil
}
