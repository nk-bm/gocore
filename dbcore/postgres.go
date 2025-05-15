package dbcore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PostgresConfig struct {
	Host        string `env:"POSTGRES_HOST; default:localhost"`
	Port        int    `env:"POSTGRES_PORT; default:5432"`
	User        string `env:"POSTGRES_USER; default:postgres"`
	Password    string `env:"POSTGRES_PASSWORD; default:postgres"`
	DBName      string `env:"POSTGRES_DB; default:postgres"`
	TablePrefix string `env:"POSTGRES_TABLE_PREFIX; default:"`
}

func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName,
	)
}

type PostgresClient struct {
	config PostgresConfig

	gormDB *gorm.DB
	sqlDB  *sql.DB
}

func NewPostgresClient(config PostgresConfig) *PostgresClient {
	return &PostgresClient{
		config: config,
	}
}

func (c *PostgresClient) Connect() error {
	db, err := sql.Open("postgres", c.config.DSN())
	if err != nil {
		return fmt.Errorf("open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	c.sqlDB = db
	c.gormDB, err = gorm.Open(postgres.Open(c.config.DSN()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: c.config.TablePrefix,
		},
	})
	if err != nil {
		return fmt.Errorf("open gorm database connection: %w", err)
	}

	return nil
}

func (c *PostgresClient) Close() error {
	return c.sqlDB.Close()
}

func (c *PostgresClient) GormDB() *gorm.DB {
	return c.gormDB
}

func (c *PostgresClient) SqlDB() *sql.DB {
	return c.sqlDB
}
