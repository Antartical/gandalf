package connections

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
PostgresConnection -> initializes postgres conection
*/
func PostgresConnection(config ConnectionConfig) *gorm.DB {
	db, err := config.open(postgres.Open(config.PostgresDSN()), &gorm.Config{})
	if err != nil {
		panic(&PostgresConnectionError{})
	}
	fmt.Println("-> Connected to postgres <-")
	return db
}
