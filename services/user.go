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
	Create(userData validators.UserCreateData) (*models.User, error)
	Read(uuid uuid.UUID) (*models.User, error)
	Update(uuid uuid.UUID, userData validators.UserUpdateData) (*models.User, error)
	Delete(uuid uuid.UUID) error
}

/*
UserService -> user's service
*/
type UserService struct {
	db *gorm.DB
}

/*
Create -> creates a new user
*/
func (service UserService) Create(userData validators.UserCreateData) (*models.User, error) {
	user := models.NewUser(
		userData.Email,
		userData.Password,
		userData.Name,
		userData.Surname,
		userData.Birthday,
		userData.Phone,
	)

	if err := service.db.Create(&user).Error; err != nil {
		return nil, UserCreateError{err}
	}

	return &user, nil
}
