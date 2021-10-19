package services

import (
	"gandalf/bindings"
	"gandalf/tests"
	"gandalf/validators"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func TestUserServiceConstructor(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		service := NewUserService(db)

		assert.Equal(service.db, db)
	})
}

func TestUserServiceCreate(t *testing.T) {
	assert := require.New(t)

	t.Run("Test user create successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := UserService{db}
		userData := validators.UserCreateData{
			Email:    "test@test.com",
			Password: "testestestestest",
			Name:     "test",
			Surname:  "test",
			Birthday: bindings.BirthDate(time.Now()),
		}

		user, err := service.Create(userData)

		assert.NoError(err)
		assert.Equal(user.Email, userData.Email)
		assert.Equal(user.Name, userData.Name)
		assert.Equal(user.Surname, userData.Surname)
		assert.Equal(user.Birthday, userData.Birthday)
		db.Unscoped().Delete(&user)
	})

	t.Run("Test user create database error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := UserService{db}
		userData := validators.UserCreateData{
			Email:    "test@test.com",
			Password: "testestestestest",
			Name:     "test",
			Surname:  "test",
			Birthday: bindings.BirthDate(time.Now()),
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

		user := tests.UserFactory()
		db.Create(&user)

		readUser, err := service.Read(user.UUID)

		assert.NoError(err)
		assert.Equal(user.ID, readUser.ID)

		db.Unscoped().Delete(&user)
	})

	t.Run("Test read user not found error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := UserService{db}
		uuid, _ := uuid.NewV4()

		_, err := service.Read(uuid)
		assert.Error(err, UserNotFoundError{nil}.Error())
	})

}

func TestUserServiceReadByEmail(t *testing.T) {
	assert := require.New(t)

	t.Run("Test read user by email successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := UserService{db}

		user := tests.UserFactory()
		db.Create(&user)

		readUser, err := service.ReadByEmail(user.Email)

		assert.NoError(err)
		assert.Equal(user.ID, readUser.ID)

		db.Unscoped().Delete(&user)
	})

	t.Run("Test read user by email not found error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := UserService{db}
		user := tests.UserFactory()
		_, err := service.ReadByEmail(user.Email)

		assert.Error(err, UserNotFoundError{nil}.Error())
	})

}

func TestUserServiceUpdate(t *testing.T) {
	assert := require.New(t)

	t.Run("Test update user successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := UserService{db}

		user := tests.UserFactory()
		db.Create(&user)

		password := "NewPassword"
		phone := "+34666666666"
		userData := validators.UserUpdateData{
			Password: password,
			Phone:    phone,
		}

		updatedUser, err := service.Update(user.UUID, userData)

		assert.NoError(err)
		assert.Equal(updatedUser.Phone, phone)
		assert.True(updatedUser.VerifyPassword(password))

		db.Unscoped().Delete(&user)
	})

	t.Run("Test update user not found error", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := UserService{db}

		uuid, _ := uuid.NewV4()
		userData := validators.UserUpdateData{
			Password: "NewPassword",
			Phone:    "+34666666666",
		}

		_, err := service.Update(uuid, userData)

		assert.Error(err, UserNotFoundError{nil}.Error())
	})

}

func TestUserServiceDelete(t *testing.T) {
	assert := require.New(t)

	t.Run("Test delete user successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := UserService{db}

		user := tests.UserFactory()
		db.Create(&user)

		err := service.Delete(user.UUID)
		assert.NoError(err)
	})

	t.Run("Test delete user error not found", func(t *testing.T) {
		db := tests.NewTestDatabase(false)
		service := UserService{db}

		user := tests.UserFactory()

		err := service.Delete(user.UUID)
		assert.Error(err, UserNotFoundError{nil}.Error())
	})

}

func TestUserServiceVerificate(t *testing.T) {
	assert := require.New(t)

	t.Run("Test verify user successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		service := UserService{db}
		user := tests.UserFactory()

		service.Verificate(&user)

		assert.True(user.Verified)
	})

}

func TestUserServiceResetPassword(t *testing.T) {
	assert := require.New(t)

	t.Run("Test verify user successfully", func(t *testing.T) {
		db := tests.NewTestDatabase(true)
		newPassword := "wowowowowowoow"
		service := UserService{db}
		user := tests.UserFactory()

		service.ResetPassword(&user, newPassword)

		assert.True(user.VerifyPassword(newPassword))
	})

}
