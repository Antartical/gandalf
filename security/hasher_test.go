package security

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type generateFromPasswordRecorder struct {
	password []byte
	cost     int
}

type compareHashAndPasswordRecorder struct {
	hashedPassword []byte
	plainPassword  []byte
}

func mockedGenerateFromPassword(err error) (func(password []byte, cost int) ([]byte, error), *generateFromPasswordRecorder) {
	recorder := new(generateFromPasswordRecorder)
	return func(password []byte, cost int) ([]byte, error) {
		*recorder = generateFromPasswordRecorder{password, cost}
		return password, err
	}, recorder
}

func mockedCompareHashAndPassword(err error) (func(hashedPassword []byte, password []byte) error, *compareHashAndPasswordRecorder) {
	recorder := new(compareHashAndPasswordRecorder)
	return func(hashedPassword []byte, password []byte) error {
		*recorder = compareHashAndPasswordRecorder{hashedPassword, password}
		return err
	}, recorder
}

func TestBcryptHasher(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		plainPassword := "Test"
		hasher := NewBcryptHasher()

		hashedPassword, _ := hasher.GeneratePassword(plainPassword)
		err := hasher.VerifyPassword(string(hashedPassword), plainPassword)

		assert.Nil(err)
	})

	t.Run("Test GeneratePassword", func(t *testing.T) {
		password := "test"
		GenerateFromPassword, recorder := mockedGenerateFromPassword(nil)
		hasher := BcryptHasher{
			generateFromPassword: GenerateFromPassword,
		}

		hasher.GeneratePassword(password)
		assert.Equal(recorder.password, []byte(password))
		assert.Equal(recorder.cost, bcrypt.DefaultCost)
	})

	t.Run("Test VerifyPassword", func(t *testing.T) {
		password := "test"
		CompareHashAndPassword, recorder := mockedCompareHashAndPassword(nil)
		hasher := BcryptHasher{
			compareHashAndPassword: CompareHashAndPassword,
		}

		hasher.VerifyPassword(password, password)

		assert.Equal(recorder.hashedPassword, []byte(password))
		assert.Equal(recorder.plainPassword, []byte(password))
	})
}
