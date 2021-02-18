package services

import (
	"gandalf/tests"
	"gandalf/validators"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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
