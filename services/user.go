package services

import (
	"gandalf/models"
	"gandalf/validators"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

/*
IUserService -> interface for user service
*/
type IUserService interface {
	Create(db *gorm.DB, userData validators.UserCreateData) (models.User, error)
	Read(db *gorm.DB, uuid uuid.UUID) (models.User, error)
	Update(db *gorm.DB, uuid uuid.UUID, userData validators.UserUpdateData) (models.User, error)
	Delete(db *gorm.DB, uuid uuid.UUID) error
}
