package seed

import (
	"log"
	"os"

	config "github.com/vnnyx/employee-management/config/api"
	"github.com/vnnyx/employee-management/pkg/database"
)

func Run() {
	log.Println("Seeding database...")

	env := os.Getenv("config")
	if env == "" {
		env = "local"
	}

	log.Println("Environment: ", env)

	cfg, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("Error loading config file: %s", err)
	}

	db, err := database.GetPostgresConnectionSeed(database.PostgresConfig{
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
		log.Fatalf("Error connecting to database: %s", err)
	}

	log.Println("Database connection established")

	// Run .sql files in the seed directory
	seedFiles, err := os.ReadDir("database/seed")
	if err != nil {
		log.Fatalf("Error reading seed directory: %s", err)
	}

	for _, file := range seedFiles {
		if file.IsDir() || file.Name()[len(file.Name())-4:] != ".sql" {
			continue // Skip directories and non-SQL files
		}

		filePath := "database/seed/" + file.Name()
		log.Printf("Executing seed file: %s\n", filePath)

		if err := database.ExecuteSQLFile(db, filePath); err != nil {
			log.Fatalf("Error executing seed file %s: %s", filePath, err)
		}
	}
}
