package validators

import "time"

/*
UserResendEmail -> user data for resend email notification
*/
type UserResendEmail struct {
	Email           string `json:"email" binding:"required,email"`
	VerificationURL string `json:"verification_url" binding:"required"`
}

/*
UserCreateData -> user data for creation
*/
type UserCreateData struct {
	Email           string    `json:"email" binding:"required,email"`
	Password        string    `json:"password" binding:"required,min=10"`
	Name            string    `json:"name" binding:"required"`
	Surname         string    `json:"surname" binding:"required"`
	Birthday        time.Time `json:"birthday" binding:"required"`
	Phone           string    `json:"phone" binding:"omitempty,e164"`
	VerificationURL string    `json:"verification_url" binding:"required"`
}

/*
UserUpdateData -> user data for update
*/
type UserUpdateData struct {
	Password string `json:"password" binding:"omitempty,min=10"`
	Phone    string `json:"phone" binding:"omitempty,e164"`
}
