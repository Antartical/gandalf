package services

import (
	"gandalf/models"
	"gandalf/validators"

	"gorm.io/gorm"
)

/*
IUserService -> interface for user service
*/
type IUserService interface {
	Create(db *gorm.DB, userData validators.UserCreate) error
	Read(db *gorm.DB, id int) models.User
	Update(db *gorm.DB, id int, userData validators.UserUpdate) error
	Delete(db *gorm.DB, id int) error
}
