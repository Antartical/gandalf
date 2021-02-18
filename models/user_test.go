package models

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type generatePasswordRecorder struct {
	password string
}

type verifyPasswordRecorder struct {
	hashedPassword string
	plainPassword  string
}

type mockedHasher struct {
	generatePasswordError    error
	verifyPasswordError      error
	generatePasswordRecorder *generatePasswordRecorder
	verifyPasswordRecorder   *verifyPasswordRecorder
}

func (hasher *mockedHasher) GeneratePassword(password string) ([]byte, error) {
	hasher.generatePasswordRecorder = &generatePasswordRecorder{password}
	return []byte(password), hasher.generatePasswordError
}

func (hasher *mockedHasher) VerifyPassword(hashedPassword string, plainPassword string) error {
	hasher.verifyPasswordRecorder = &verifyPasswordRecorder{hashedPassword, plainPassword}
	return hasher.verifyPasswordError
}

func newMockedHasher(generatePasswordError error, verifyPasswordError error) *mockedHasher {
	return &mockedHasher{generatePasswordError, verifyPasswordError, new(generatePasswordRecorder), new(verifyPasswordRecorder)}
}

func TestUserModel(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		email := "test@example.com"
		name := "test"
		surname := "test test"
		birthday := time.Now()
		phone := "+34666"

		user := NewUser(email, "test", name, surname, birthday, phone)

		assert.Equal(user.Email, email)
		assert.Equal(user.Name, name)
		assert.Equal(user.Surname, surname)
		assert.Equal(user.Birthday, birthday)
		assert.Equal(user.Phone, phone)
	})

	t.Run("Test AfterFind gorm hook", func(t *testing.T) {
		user := User{}
		err := user.AfterFind(nil)

		assert.NoError(err)
		assert.NotNil(user.hasher)
	})

	t.Run("Test SetPassword successfully", func(t *testing.T) {
		plainPassword := "test"
		hasher := newMockedHasher(nil, nil)
		user := User{hasher: hasher}

		user.SetPassword(plainPassword)
		assert.Equal(hasher.generatePasswordRecorder.password, plainPassword)
	})

	t.Run("Test SetPassowrd panics", func(t *testing.T) {
		expectedError := errors.New("woohps")
		plainPassword := "test"
		hasher := newMockedHasher(expectedError, nil)
		user := User{hasher: hasher}

		assert.PanicsWithError(expectedError.Error(), func() { user.SetPassword(plainPassword) })
	})

	t.Run("Test VerifyPassword successfully", func(t *testing.T) {
		plainPassword := "test"
		hasher := newMockedHasher(nil, nil)
		user := User{hasher: hasher}

		user.SetPassword(plainPassword)
		match := user.VerifyPassword(plainPassword)

		assert.True(match)
		assert.Equal(hasher.verifyPasswordRecorder.hashedPassword, user.Password)
		assert.Equal(hasher.verifyPasswordRecorder.plainPassword, plainPassword)
	})

	t.Run("Test VerifyPassword wrongly", func(t *testing.T) {
		plainPassword := "test"
		hasher := newMockedHasher(nil, errors.New("not equals"))
		user := User{hasher: hasher}

		user.SetPassword(plainPassword)
		match := user.VerifyPassword(plainPassword)

		assert.False(match)
	})
}
