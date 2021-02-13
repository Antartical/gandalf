package connections

import (
	"fmt"
	"os"

	"gorm.io/gorm"
)

/*
ConnectionConfig -> contains connection data
*/
type ConnectionConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	open     func(gorm.Dialector, *gorm.Config) (*gorm.DB, error)
}

/*
PostgresDSN -> get DSN for postgres configuration
*/
func (config ConnectionConfig) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Name,
		config.Password,
	)
}

/*
NewPostrgresConnectionConfig -> returns new postgres connection
configuration
*/
func NewPostrgresConnectionConfig() ConnectionConfig {
	return ConnectionConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DB"),
		open:     gorm.Open,
	}
}
