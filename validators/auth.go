package validators

/*
Credentials -> User credentials through the one he can autenticates himself.
*/
type Credentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,alphanumunicode,min=10"`
}
