package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"bem_be/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the database connection
var DB *gorm.DB

// Initialize connects to the database and creates tables if they don't exist
func Initialize() {
	var err error

	// Get database connection details from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)

	// Create MySQL DSN: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbname)

	// Connect to MySQL database
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Get the underlying SQL DB to configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error getting underlying SQL DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to MySQL database successfully")

	log.Println("Starting database migration...")

	// Example migration sequence (sama seperti versi PostgreSQL)
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Error auto-migrating User model: %v\n", err)
	}
	log.Println("User table migrated successfully")

	err = DB.AutoMigrate(&models.Student{})
	if err != nil {
		log.Fatalf("Error auto-migrating Student model: %v\n", err)
	}
	log.Println("Student table migrated successfully")


	// Lanjutkan migrasi model lain seperti Course, StudentGroup, dll sesuai urutan versi PostgreSQL

	log.Println("Database schema migrated successfully")
}

// Close closes the database connection
func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			log.Printf("Error getting underlying SQL DB: %v", err)
			return
		}
		sqlDB.Close()
	}
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}
