package connections

import (
	"gorm.io/gorm"
)

// Gorm database connection interface
type GormDBConnection interface {
	Connect() *gorm.DB
}
