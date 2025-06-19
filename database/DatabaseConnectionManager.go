package database

import (
	"fmt"
	"os"
	"time"

	zlogger "github.com/siva-chegondi/go-utils/logger" // customized zerolog
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var log = zlogger.DefaultLogger

const (
	MAX_IDLE_CONNS    = 10
	MAX_OPEN_CONNS    = 100
	CONN_MAX_LIFETIME = time.Hour
)

func InitDB() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}

	var err error
	// Configure GORM logger
	newLogger := logger.New(
		log,
		logger.Config{
			SlowThreshold:             time.Second, // Log SQL queries taking longer than 1 second
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Database connection string
	dbHost, err := getEnv("DB_HOST")
	if err != nil {
		return nil, err
	}
	dbPasswd, err := getEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}
	dbUser, err := getEnv("DB_USER")
	if err != nil {
		return nil, err
	}
	dbName, err := getEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		dbHost, dbUser, dbPasswd, dbName,
		getEnvOrDefault("DB_PORT", "5432"),
		getEnvOrDefault("DB_SSLMODE", "disable"),
	)

	// Open database connection
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(MAX_IDLE_CONNS)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(MAX_OPEN_CONNS)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(CONN_MAX_LIFETIME)

	log.Info().Msg("Database connection established")
	return DB, nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Helper function to get environment variable with fallback
func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnv(key string) (string, error) {
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", fmt.Errorf("environment variable %s is mandatory", key)
}
