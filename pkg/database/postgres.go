package database

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/srcndev/message-service/config"
	"github.com/srcndev/message-service/internal/domain"
	applogger "github.com/srcndev/message-service/pkg/logger"
)

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, ErrDatabaseConnectionFailed.WithError(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, ErrDatabaseInstanceFailed.WithError(err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		return nil, ErrDatabasePingFailed.WithError(err)
	}

	applogger.Info("Database connection established")
	return db, nil
}

// AutoMigrate runs database migrations for all models
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.Message{}); err != nil {
		return ErrDatabaseMigrationFailed.WithError(err)
	}

	applogger.Info("Database auto-migration completed")
	return nil
}
