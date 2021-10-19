package serializers

import (
	"gandalf/tests"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserSerializer(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor", func(t *testing.T) {
		user := tests.UserFactory()
		userSerializer := NewUserSerializer(user)

		assert.Equal(userSerializer.Data.Name, user.Name)
		assert.Equal(userSerializer.Data.Email, user.Email)
		assert.Equal(userSerializer.Data.Surname, user.Surname)
		assert.Equal(userSerializer.Data.Birthday, user.Birthday)
		assert.Equal(userSerializer.Data.Phone, user.Phone)
	})
}
