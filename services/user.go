package services

import (
	"gandalf/models"
	"gandalf/validators"

	"github.com/gofrs/uuid"
)

/*
IUserService -> interface for user service
*/
type IUserService interface {
	Create(userData validators.UserCreateData) (models.User, error)
	Read(uuid uuid.UUID) (models.User, error)
	Update(uuid uuid.UUID, userData validators.UserUpdateData) (models.User, error)
	Delete(uuid uuid.UUID) error
}
