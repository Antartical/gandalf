package connections

import (
	"errors"
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

func TestPostgresConnection(t *testing.T) {
	assert := require.New(t)

	t.Run("Setup connection successfully", func(t *testing.T) {
		mockedOpen, recorder := mockConfigOpen(nil)
		mockedConfig := ConnectionConfig{
			Host:     "Test",
			Port:     "Test",
			User:     "Test",
			Password: "Test",
			Name:     "Test",
			open:     mockedOpen,
		}

		PostgresConnection(mockedConfig)
		assert.Equal(recorder.config, &gorm.Config{})
		assert.Equal(recorder.dialector, postgres.Open(mockedConfig.PostgresDSN()))
	})

	t.Run("Setup connection error", func(t *testing.T) {
		mockedOpen, _ := mockConfigOpen(errors.New("connection error"))
		mockedConfig := ConnectionConfig{
			Host:     "Test",
			Port:     "Test",
			User:     "Test",
			Password: "Test",
			Name:     "Test",
			open:     mockedOpen,
		}

		assert.PanicsWithValue(&PostgresConnectionError{}, func() { PostgresConnection(mockedConfig) })
	})
}
