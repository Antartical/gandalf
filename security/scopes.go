package security

// Security scopes
const (
	ScopeUserAuthorizationCode = "user:me:authorization-code"
	ScopeUserAuthorizeApp      = "user:me:authorized-app"
	ScopeUserVerify            = "user:me:verify"
	ScopeUserChangePassword    = "user:me:change-password"

	ScopeUserRead   = "user:me:read"
	ScopeUserWrite  = "user:me:write"
	ScopeUserDelete = "user:me:delete"

	ScopeAppRead = "app:read"
)

// Group scopes
var (
	GroupUserSelf          = []string{ScopeUserAuthorizeApp, ScopeUserRead, ScopeUserWrite, ScopeUserDelete}
	GroupUserOauth2Request = []string{ScopeUserAuthorizeApp, ScopeUserRead, ScopeAppRead}
	GroupAdmin             = []string{ScopeUserRead, ScopeUserWrite, ScopeUserDelete, ScopeAppRead}
)
