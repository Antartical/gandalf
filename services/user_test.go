package services

import (
	"gandalf/models"
	"gandalf/tests"
	"gandalf/validators"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func userFactory() models.User {
	userData := validators.UserCreateData{
		Email:    "test@test.com",
		Password: "testestestestest",
		Name:     "test",
		Surname:  "test",
		Birthday: time.Now(),
	}
	return models.NewUser(
		userData.Email,
		userData.Password,
		userData.Name,
		userData.Surname,
		userData.Birthday,
		userData.Phone,
	)
}

func TestUserServiceCreate(t *testing.T) {
	assert := require.New(t)

	t.Run("Test user create successfully", func(t *testing.T) {
		service := UserService{tests.NewTestDatabase(true)}
		userData := validators.UserCreateData{
			Email:    "test@test.com",
			Password: "testestestestest",
			Name:     "test",
			Surname:  "test",
			Birthday: time.Now(),
		}

		user, err := service.Create(userData)

		assert.NoError(err)
		assert.Equal(user.Email, userData.Email)
		assert.Equal(user.Name, userData.Name)
		assert.Equal(user.Surname, userData.Surname)
		assert.Equal(user.Birthday, userData.Birthday)
	})

	t.Run("Test user create database error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := UserService{db}
		userData := validators.UserCreateData{
			Email:    "test@test.com",
			Password: "testestestestest",
			Name:     "test",
			Surname:  "test",
			Birthday: time.Now(),
		}

		user, _ := service.Create(userData)
		_, err := service.Create(userData)

		assert.Error(err, UserCreateError{nil}.Error())
		db.Unscoped().Delete(&user)
	})
}

func TestUserServiceRead(t *testing.T) {
	assert := require.New(t)

	t.Run("Test read user successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := UserService{db}

		user := userFactory()
		db.Create(&user)

		readUser, err := service.Read(user.UUID)

		assert.NoError(err)
		assert.Equal(user.ID, readUser.ID)

		db.Unscoped().Delete(&user)
	})

	t.Run("Test read user not found error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := UserService{db}
		user := userFactory()
		_, err := service.Read(user.UUID)

		assert.Error(err, UserNotFoundError{nil}.Error())
	})

}

// func TestUserServiceUpdate(t *testing.T) {
// 	assert := require.New(t)

// 	t.Run("Test update user successfully", func(t *testing.T) {

// 	})

// 	t.Run("Test update user not found error", func(t *testing.T) {

// 	})

// }

// func TestUserServiceDelete(t *testing.T) {
// 	assert := require.New(t)

// 	t.Run("Test delete user successfully", func(t *testing.T) {

// 	})

// }

// func TestUserServiceSoftDelete(t *testing.T) {
// 	assert := require.New(t)

// 	t.Run("Test soft delete user successfully", func(t *testing.T) {

// 	})

// }
