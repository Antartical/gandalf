package connections

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
GormPostgresConnection -> postgres gorm connection
*/
type GormPostgresConnection struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	OpenDb   func(gorm.Dialector, *gorm.Config) (*gorm.DB, error)
}

/*
Connect -> return the database connection
*/
func (connection GormPostgresConnection) Connect() *gorm.DB {
	addr := connection.getPostgresDSN()
	db, err := connection.OpenDb(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		panic(&DatabaseConnectionError{addr})
	}
	return db
}

func (connection GormPostgresConnection) getPostgresDSN() string {
	return fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
		connection.Host,
		connection.Port,
		connection.User,
		connection.Name,
		connection.Password,
	)
}

/*
NewGormPostgresConnection -> returns new postgres connection
*/
func NewGormPostgresConnection() GormPostgresConnection {
	return GormPostgresConnection{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DB"),
		OpenDb:   gorm.Open,
	}
}
