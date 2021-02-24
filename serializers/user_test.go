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

		assert.Equal(userSerializer.Data.Name, name)
		assert.Equal(userSerializer.Data.Email, email)
		assert.Equal(userSerializer.Data.Surname, surname)
		assert.Equal(userSerializer.Data.Birthday, birthday)
		assert.Equal(userSerializer.Data.Phone, phone)
	})
}
