package connections

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func mockConfigOpen(err error) (func(gorm.Dialector, *gorm.Config) (*gorm.DB, error), *configOpenRecorder) {
	recorder := new(configOpenRecorder)
	return func(dialector gorm.Dialector, config *gorm.Config) (*gorm.DB, error) {
		*recorder = configOpenRecorder{dialector, config}
		return nil, err
	}, recorder
}

type configOpenRecorder struct {
	dialector gorm.Dialector
	config    *gorm.Config
}

func TestGormPostgresConnection(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		connection := NewGormPostgresConnection()
		assert.Equal(connection.Host, os.Getenv("POSTGRES_HOST"))
		assert.Equal(connection.Port, os.Getenv("POSTGRES_PORT"))
		assert.Equal(connection.User, os.Getenv("POSTGRES_USER"))
		assert.Equal(connection.Password, os.Getenv("POSTGRES_PASSWORD"))
		assert.Equal(connection.Name, os.Getenv("POSTGRES_DB"))
	})

	t.Run("Test generate postgres DSN", func(t *testing.T) {
		host := "test"
		port := "test"
		user := "test"
		password := "test"
		name := "test"
		expectedDSN := fmt.Sprintf(
			"host=%v port=%v user=%v dbname=%v password=%v sslmode=disable",
			host,
			port,
			user,
			name,
			password,
		)

		mockedOpen, _ := mockConfigOpen(nil)
		mockedConnection := GormPostgresConnection{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			Name:     name,
			open:     mockedOpen,
		}

		assert.Equal(mockedConnection.getPostgresDSN(), expectedDSN)
	})

	t.Run("Setup connection successfully", func(t *testing.T) {
		mockedOpen, recorder := mockConfigOpen(nil)
		mockedConnection := GormPostgresConnection{
			Host:     "Test",
			Port:     "Test",
			User:     "Test",
			Password: "Test",
			Name:     "Test",
			open:     mockedOpen,
		}

		mockedConnection.Connect()

		assert.Equal(recorder.config, &gorm.Config{})
		assert.Equal(recorder.dialector, postgres.Open(mockedConnection.getPostgresDSN()))
	})

	t.Run("Setup connection error", func(t *testing.T) {
		mockedOpen, _ := mockConfigOpen(errors.New("connection error"))
		mockedConnection := GormPostgresConnection{
			Host:     "Test",
			Port:     "Test",
			User:     "Test",
			Password: "Test",
			Name:     "Test",
			open:     mockedOpen,
		}

		expectedError := DatabaseConnectionError{mockedConnection.getPostgresDSN()}
		assert.PanicsWithError(expectedError.Error(), func() { mockedConnection.Connect() })
	})
}
