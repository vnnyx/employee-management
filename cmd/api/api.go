package api

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	configapi "github.com/vnnyx/employee-management/config/api"
	_ "github.com/vnnyx/employee-management/docs"
	"github.com/vnnyx/employee-management/internal/server"
	"github.com/vnnyx/employee-management/pkg/database"
	"github.com/vnnyx/employee-management/pkg/logger"
	"github.com/vnnyx/employee-management/pkg/redis"
)

// @title						Employee Management Service
// @version					1.0
// @description				Employee Management Service API Docs
// @Schemes					http
// @host						localhost:9000
// @BasePath					/external/api
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
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

	redisClient, err := redis.GetRedisConnection(redis.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		Database: cfg.Redis.Database,
	})
	if err != nil {
		log.Fatalf("Error on connecting to Redis: %s", err)
	}

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		redisClient.Close()
	}()

	appLogger.Info("App stopped")
}
