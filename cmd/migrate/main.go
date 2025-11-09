package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/srcndev/message-service/config"
	"github.com/srcndev/message-service/internal/domain"
	"github.com/srcndev/message-service/pkg/database"
	"github.com/srcndev/message-service/pkg/logger"
	"github.com/srcndev/message-service/seed"
)

func main() {
	// Define CLI flags
	seedFlag := flag.Bool("seed", false, "Run database seeding after migration")
	helpFlag := flag.Bool("help", false, "Show help message")
	flag.Parse()

	if *helpFlag {
		printHelp()
		os.Exit(0)
	}

	logger.Info("Starting database migration tool...")

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatal("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	// Run migrations
	logger.Info("Running database migrations...")
	if err := db.AutoMigrate(&domain.Message{}); err != nil {
		logger.Fatal("Migration failed: %v", err)
	}
	logger.Info("✓ Migrations completed successfully")

	// Run seeding if flag is set
	if *seedFlag {
		logger.Info("Running database seeding...")
		seeder := seed.NewSeeder(db)
		if err := seeder.Run(); err != nil {
			logger.Fatal("Seeding failed: %v", err)
		}
		logger.Info("✓ Database seeded successfully")
	}

	logger.Info("Done!")
}

func printHelp() {
	fmt.Println("Database Migration Tool")
	fmt.Println("\nUsage:")
	fmt.Println("  migrate [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -seed    Run database seeding after migration")
	fmt.Println("  -help    Show this help message")
	fmt.Println("\nExamples:")
	fmt.Println("  # Run migrations only")
	fmt.Println("  ./bin/migrate")
	fmt.Println("")
	fmt.Println("  # Run migrations and seed data")
	fmt.Println("  ./bin/migrate -seed")
}
