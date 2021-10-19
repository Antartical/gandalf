package connections

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Struct with the necessary data to create a new gorm connection
// with a postgres database
type GormPostgresConnection struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	OpenDb   func(gorm.Dialector, ...gorm.Option) (*gorm.DB, error)
}

// Connects to the database by using the addr
func (connection GormPostgresConnection) Connect() *gorm.DB {
	addr := connection.getPostgresDSN()
	db, err := connection.OpenDb(postgres.Open(addr), &gorm.Config{})
	if err != nil {
		panic(&DatabaseConnectionError{addr})
	}
	return db
}

// Creates the postgres uri for the driver in order to connect with it
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

// Creates a new postgres connection
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
