package validators

// Validator for user credentials data
type Credentials struct {
	Email    string `json:"email" binding:"required,email" example:"johndoe@example.com"`
	Password string `json:"password" binding:"required,min=10" example:"My@appPassw0rd"`
}

// Validator for user access tokens data
type AuthTokens struct {
	AcessToken   string `json:"access_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	RefreshToken string `json:"refresh_token" binding:"required" example:"kpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyf"`
}
