package gocore

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nk-bm/gocore/dbcore"
	"github.com/nk-bm/gocore/env"
	"github.com/nk-bm/gocore/gincore"
	"go.uber.org/zap"
)

type AppOptions struct {
	Logger              *zap.Logger
	DisableGlobalLogger bool
	DisableMigrations   bool
}

type AppConfig struct {
	IsProd         bool `env:"PRODUCTION; default:false"`
	PostgresConfig dbcore.PostgresConfig
	GinConfig      gincore.Config
	Options        AppOptions
}

type App struct {
	Name      string
	GinServer *gincore.Server
	Postgres  *dbcore.PostgresClient
	Migrator  *dbcore.Migrator

	L *zap.Logger
}

func NewDefaultApp(name string) (*App, error) {
	if err := godotenv.Load(); err != nil {
		// Ignore error if .env file doesn't exist
		if !strings.Contains(err.Error(), "no such file") {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	var config AppConfig
	if err := env.LoadEnv(&config); err != nil {
		return nil, err
	}

	return NewApp(name, config, []dbcore.Migration{})
}

func NewApp(appName string, config AppConfig, migrations []dbcore.Migration) (*App, error) {
	if strings.Contains(appName, " ") {
		return nil, fmt.Errorf("service name cannot contain spaces")
	}

	if !config.Options.DisableGlobalLogger {
		if config.Options.Logger == nil {
			InitGlobalLogger(config.IsProd)
			config.Options.Logger = zap.L()
		} else {
			SetGlobalLogger(config.Options.Logger)
		}
	}

	L := config.Options.Logger
	L.Info("Initializing core components...", zap.String("app_name", appName))

	defer func() {
		if r := recover(); r != nil {
			L.Error("Core components initialization failed", zap.Any("error", r))
		}
	}()

	postgres := dbcore.NewPostgresClient(config.PostgresConfig)
	if err := postgres.Connect(); err != nil {
		return nil, err
	}
	var migrator *dbcore.Migrator
	if !config.Options.DisableMigrations {
		migrator = dbcore.NewMigrator(postgres.GormDB(), config.Options.Logger, appName, migrations)
		if err := migrator.Run(); err != nil {
			return nil, err
		}
	}

	ginServer := gincore.NewServer(config.GinConfig, config.Options.Logger)

	L.Info("Core components initialized", zap.String("app_name", appName))
	return &App{
		Name:      appName,
		GinServer: ginServer,
		Postgres:  postgres,
		Migrator:  migrator,
		L:         config.Options.Logger,
	}, nil
}

func (s *App) Start() error {
	s.L.Info("Starting core components...")
	s.L.Info("Running migrations...")
	if err := s.Migrator.Run(); err != nil {
		return err
	}

	s.L.Info("Migrations completed")

	s.L.Info("Starting Gin server...")
	return s.GinServer.Start()
}
