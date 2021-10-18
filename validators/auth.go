package validators

/*
Credentials -> User credentials through the one he can autenticates himself.
*/
type Credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=10"`
}

/*
AuthTokens -> user tokens to be refreshed
*/
type AuthTokens struct {
	AcessToken   string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}
