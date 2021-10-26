package models

import (
	"gandalf/security"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"syreclabs.com/go/faker"
)

func TestClaimModel(t *testing.T) {
	assert := require.New(t)

	t.Run("Test constructor success", func(t *testing.T) {
		redirectUrl := faker.Internet().Url()
		authorizationCode := faker.RandomString(10)
		scopes := []string{security.ScopeUserAuthorizationCode}
		user := User{Name: faker.Name().FirstName()}
		app := App{Name: faker.Company().Name()}

		claim := NewClaim(redirectUrl, authorizationCode, scopes, user, app)
		assert.Equal(redirectUrl, claim.RedirectUrl)
		assert.Equal(authorizationCode, claim.AuthorizationCode)
		assert.Equal(pq.StringArray(scopes), claim.Scopes)
		assert.Equal(user.Name, claim.User.Name)
		assert.Equal(app.Name, claim.App.Name)
	})

}
