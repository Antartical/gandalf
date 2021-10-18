package security

/*
Authentication scopes
*/
const (
	ScopeUserAuthorizeApp   = "user:authorized-app"
	ScopeUserRead           = "user:read"
	ScopeUserVerify         = "user:verify"
	ScopeUserChangePassword = "user:change-password"
	ScopeUserWrite          = "user:write"
	ScopeUserDelete         = "user:delete"
)

/*
Group scopes
*/
var (
	GroupUserAll = []string{ScopeUserAuthorizeApp, ScopeUserRead, ScopeUserWrite, ScopeUserDelete}
)
