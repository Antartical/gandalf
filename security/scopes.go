package security

// Security scopes
const (
	ScopeUserAuthorizationCode = "user:authorization-code"
	ScopeUserAuthorizeApp      = "user:authorized-app"
	ScopeUserVerify            = "user:verify"
	ScopeUserChangePassword    = "user:change-password"
	ScopeUserRead              = "user:read"
	ScopeUserWrite             = "user:write"
	ScopeUserDelete            = "user:delete"

	ScopeAppRead = "app:read"
)

// Group scopes
var (
	GroupUserSelf          = []string{ScopeUserAuthorizeApp, ScopeUserRead, ScopeUserWrite, ScopeUserDelete}
	GroupUserOauth2Request = []string{ScopeUserAuthorizeApp, ScopeUserRead, ScopeAppRead}
)
