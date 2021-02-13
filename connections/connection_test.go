package connections

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestConnection(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		config := NewPostrgresConnectionConfig()
		assert.Equal(config.Host, os.Getenv("POSTGRES_HOST"))
		assert.Equal(config.Port, os.Getenv("POSTGRES_PORT"))
		assert.Equal(config.User, os.Getenv("POSTGRES_USER"))
		assert.Equal(config.Password, os.Getenv("POSTGRES_PASSWORD"))
		assert.Equal(config.Name, os.Getenv("POSTGRES_DB"))
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

		config := ConnectionConfig{
			Host:     host,
			Port:     port,
			User:     user,
			Password: password,
			Name:     name,
			open:     gorm.Open,
		}

		assert.Equal(config.PostgresDSN(), expectedDSN)
	})
}
