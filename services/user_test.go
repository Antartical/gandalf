package services

import (
	"gandalf/connections"
	"gandalf/validators"
	"testing"

	"time"

	"github.com/stretchr/testify/require"
)

func TestUserServiceCreate(t *testing.T) {
	assert := require.New(t)

	t.Run("Test user create successfully", func(t *testing.T) {
		db := connections.NewGormPostgresConnection().Connect()
		service := UserService{db}
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

	})
}
