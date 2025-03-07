package db

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Comvoca-AI/comvoca-admin-back/config"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/entity"
	"github.com/Comvoca-AI/comvoca-admin-back/internal/logger"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.AppConfig.Database.Host, config.AppConfig.Database.Port, config.AppConfig.Database.Username, config.AppConfig.Database.Password, config.AppConfig.Database.DBName)

	var err error
	// Create a GORM DB instance using the *sql.DB
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &logger.CustomLogger{}, // Attach our custom logger
	})
	if err != nil {
		panic("Failed to create GORM DB instance")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}

	// Configure connection pool settings
	sqlDB.SetMaxIdleConns(10)                  // Maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)                 // Maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour)        // Maximum lifetime of a connection
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // Maximum idle time for a connection

	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		panic("Failed to enable uuid-ossp extension")
	}

	//migrations.RunMigrations(sqlDB)
	if config.AppConfig.Database.RunMigrations {
		fmt.Println("Starting migrations...")
		err := db.AutoMigrate(
			&entity.User{},
			&entity.Organization{},
			&entity.DailySchedule{},
			&entity.Speciality{},
		)

		if err == nil && db.Migrator().HasTable(&entity.Speciality{}) {
			if err := db.First(&entity.Speciality{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
				SeedSpecialities(db)
			}
		}
		if err != nil {
			fmt.Printf("Migration failed: %v\n", err)
			panic(err)
		}
		fmt.Println("Migrations completed successfully.")
	}

	return db
}
