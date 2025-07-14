package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create migration driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal("Failed to create migration driver:", err)
	}

	// Create migrator
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal("Failed to create migrator:", err)
	}

	// Get command line argument
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go [up|down|version|create <name>]")
	}

	command := os.Args[1]

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to run migrations:", err)
		}
		fmt.Println("Migrations completed successfully")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal("Failed to rollback migrations:", err)
		}
		fmt.Println("Migrations rolled back successfully")

	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatal("Failed to get version:", err)
		}
		fmt.Printf("Current version: %d, Dirty: %v\n", version, dirty)

	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Usage: go run main.go create <migration_name>")
		}
		createMigration(os.Args[2])

	default:
		log.Fatal("Unknown command. Use: up, down, version, or create")
	}
}

func createMigration(name string) {
	// Create migrations directory if it doesn't exist
	if err := os.MkdirAll("migrations", 0755); err != nil {
		log.Fatal("Failed to create migrations directory:", err)
	}

	// Generate timestamp
	timestamp := fmt.Sprintf("%d", getCurrentTimestamp())

	// Create migration files
	upFile := fmt.Sprintf("migrations/%s_%s.up.sql", timestamp, name)
	downFile := fmt.Sprintf("migrations/%s_%s.down.sql", timestamp, name)

	// Create up migration file
	if err := os.WriteFile(upFile, []byte("-- Migration up\n"), 0644); err != nil {
		log.Fatal("Failed to create up migration file:", err)
	}

	// Create down migration file
	if err := os.WriteFile(downFile, []byte("-- Migration down\n"), 0644); err != nil {
		log.Fatal("Failed to create down migration file:", err)
	}

	fmt.Printf("Created migration files:\n- %s\n- %s\n", upFile, downFile)
}

func getCurrentTimestamp() int64 {
	return 20240101000000 // This should be replaced with actual timestamp generation
}