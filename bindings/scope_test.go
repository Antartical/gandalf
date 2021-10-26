package bindings

import (
	"fmt"
	"gandalf/security"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScope(t *testing.T) {
	assert := require.New(t)

	t.Run("Test binding interface", func(t *testing.T) {
		var scope Scope
		scope.UnmarshalJSON([]byte(security.ScopeAppRead))

		marshaledScope, _ := scope.MarshalJSON()
		expectedMarshaledScope := "\"app:read\""

		assert.Equal(security.ScopeAppRead, scope.ToString())
		assert.Equal(expectedMarshaledScope, string(marshaledScope))
	})

	t.Run("Test binding error", func(t *testing.T) {
		var scope Scope
		fakeScope := "fake scope"
		expectedError := ScopeNotFoundError{
			Scope: fakeScope,
		}
		expectedMsg := fmt.Sprintf("%s, does not exists", fakeScope)
		err := scope.UnmarshalJSON([]byte(fakeScope))

		assert.Error(err, expectedError)
		assert.Equal(expectedError.Error(), expectedMsg)
	})

	t.Run("Test ScopeArrayToStringArray", func(t *testing.T) {
		scopes := []Scope{security.ScopeAppRead, security.ScopeUserAuthorizationCode}
		stringScopes := ScopeArrayToStringArray(scopes)

		expectedStringScopes := []string{security.ScopeAppRead, security.ScopeUserAuthorizationCode}
		assert.Equal(expectedStringScopes, stringScopes)
	})
}
