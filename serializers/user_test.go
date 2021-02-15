package serializers

import (
	"gandalf/models"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestUserSerializer(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		email := "test@example.com"
		name := "test"
		surname := "test test"
		birthday := time.Now()
		phone := "+34666"

		user := models.NewUser(email, "test", name, surname, birthday, phone)
		userSerializer := NewUserSerializer(user)

		assert.Equal(userSerializer.Name, name)
		assert.Equal(userSerializer.Email, email)
		assert.Equal(userSerializer.Surname, surname)
		assert.Equal(userSerializer.Birthday, birthday)
		assert.Equal(userSerializer.Phone, phone)
	})
}
