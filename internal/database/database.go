package database

import (
	"log"
	"os"
	"time"

	"github.com/MiguelCVA/mhook-backend/internal/database/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func StartDB() {
	str := os.Getenv("DB_URL")

	database, err := gorm.Open(postgres.Open(str), &gorm.Config{})
	if err != nil {
		log.Fatal("Error: " + err.Error())
	}
	db = database

	dbConfig, _ := db.DB()

	dbConfig.SetConnMaxIdleTime(10)
	dbConfig.SetMaxOpenConns(100)
	dbConfig.SetConnMaxLifetime(time.Hour)

	migrations.RunMigrations(db)
}

func CloseConn() error {
	dbConfig, err := db.DB()
	if err != nil {
		return err
	}

	err = dbConfig.Close()
	if err != nil {
		return err
	}

	return nil
}

func GetDatabase() *gorm.DB {
	return db
}
