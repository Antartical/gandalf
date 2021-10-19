package tests

import (
	"gandalf/connections"
	"gandalf/models"
	"os"

	"gorm.io/gorm"
)

// Creates a test database connection
func NewTestDatabase(dryRun bool) *gorm.DB {
	connection := connections.GormPostgresConnection{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DB_TEST"),
		OpenDb:   gorm.Open,
	}

	db := connection.Connect()
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.App{})

	return db.Session(&gorm.Session{DryRun: dryRun})
}
