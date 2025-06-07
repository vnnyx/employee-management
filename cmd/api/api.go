package api

import (
	"log"
	"os"

	configapi "github.com/vnnyx/employee-management/config/api"
	"github.com/vnnyx/employee-management/internal/server"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/logger"
)

func Run() {
	log.Println("Starting app")

	env := os.Getenv("config")
	if env == "" {
		env = "local"
	}

	log.Println("Environment: ", env)

	cfg, err := configapi.LoadConfig(env)
	if err != nil {
		log.Fatalf("Error on loading config file: %s", err)
	}

	appLogger, err := logger.InitLogger(logger.LoggerConfig{
		Mode:  cfg.Logger.Mode,
		Level: cfg.Logger.Level,
	})
	if err != nil {
		log.Fatalf("Error on initializing logger: %s", err)
	}
	defer func() { _ = appLogger.Sync() }()

	db, err := database.GetPostgresConnection(database.PostgresConfig{
		Host:                cfg.Postgres.Host,
		Port:                cfg.Postgres.Port,
		User:                cfg.Postgres.User,
		Password:            cfg.Postgres.Password,
		DBName:              cfg.Postgres.DBName,
		SSLMode:             cfg.Postgres.SSLMode,
		MaxConn:             cfg.Postgres.MaxConn,
		EnableObservability: *cfg.Observability.Enable,
	})
	if err != nil {
		appLogger.Fatalf("Error on connecting to database: %s", err)
	}

	server := server.NewServer(cfg, appLogger, db)
	if err = server.Run(); err != nil {
		log.Fatal(err)
	}

	appLogger.Info("App stopped")
}
