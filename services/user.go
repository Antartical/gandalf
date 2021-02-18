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
	SoftDelete(uuid uuid.UUID) error
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

/*
Read -> read user from database by his UUID
*/
func (service UserService) Read(uuid uuid.UUID) (*models.User, error) {
	var user models.User
	if err := service.db.Where(&models.User{UUID: uuid}).First(&user).Error; err != nil {
		return nil, UserNotFoundError{err}
	}
	return &user, nil
}

/*
Update -> updates the user which belongs to the given ID according to
the given user data
*/
func (service UserService) Update(uuid uuid.UUID, userData validators.UserUpdateData) (*models.User, error) {
	user, err := service.Read(uuid)
	if err != nil {
		return nil, err
	}

	if userData.Password != "" {
		user.SetPassword(userData.Password)
	}

	if userData.Phone != "" {
		user.Phone = userData.Phone
	}

	service.db.Save(user)
	return user, nil
}

/*
Delete -> removes the user which belongs to the given id from the database.
*/
func (service UserService) Delete(uuid uuid.UUID) error {
	if err := service.db.Unscoped().Where(&models.User{UUID: uuid}).Delete(&models.User{}).Error; err != nil {
		return UserNotFoundError{err}
	}
	return nil
}

/*
SoftDelete -> set the field `deletes_at` of the user but it will sitll alive
in database. Soft deleted users will not appear as result of any query that
not includes `unscoped`
*/
func (service UserService) SoftDelete(uuid uuid.UUID) error {
	if err := service.db.Unscoped().Where(&models.User{UUID: uuid}).Delete(&models.User{}).Error; err != nil {
		return UserNotFoundError{err}
	}
	return nil
}
