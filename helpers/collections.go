package helpers

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Check if the given element is present on the given array
func PqStringArrayContains(pqArray pq.StringArray, element interface{}) bool {
	for _, accessedElement := range pqArray {
		if accessedElement == element {
			return true
		}
	}
	return false
}

// Paginates the database results
func DBPaginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
