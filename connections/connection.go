package connections

import (
	"gorm.io/gorm"
)

/*
GormDBConnection -> Represents a gorm database connection
*/
type GormDBConnection interface {
	Connect() *gorm.DB
}
