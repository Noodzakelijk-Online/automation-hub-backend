package infra

import (
	"automation-hub-backend/internal/config"
	"automation-hub-backend/internal/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDatabase(user, password, dbName, dbHost string, dbPort int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		dbHost, user, password, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func GetDefaultDB() (*gorm.DB, error) {
	db, err := NewPostgresDatabase(config.AppConfig.DbUser, config.AppConfig.DbPassword,
		config.AppConfig.DbName, config.AppConfig.DbHost, config.AppConfig.DbPort)
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func RunMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Automation{}); err != nil {
		return err
	}
	return nil
}
