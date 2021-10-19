package validators

// Validator for send verification email with pelipper
type PelipperUserVerifyEmail struct {
	Email            string `binding:"required,email"`
	Name             string `binding:"required"`
	Subject          string `binding:"required"`
	VerificationLink string `binding:"required"`
}

// Validator for send change password email with pelipper
type PelipperUserChangePassword struct {
	Email              string `binding:"required,email"`
	Name               string `binding:"required"`
	Subject            string `binding:"required"`
	ChangePasswordLink string `binding:"required"`
}
