package validators

/*
PelipperUserVerifyEmail -> pelipper user verify email struct
*/
type PelipperUserVerifyEmail struct {
	Email            string `binding:"required,email"`
	Name             string `binding:"required"`
	Subject          string `binding:"required"`
	VerificationLink string `binding:"required"`
}

/*
PelipperUserChangePassword -> pelipper user change password email struct
*/
type PelipperUserChangePassword struct {
	Email              string `binding:"required,email"`
	Name               string `binding:"required"`
	Subject            string `binding:"required"`
	ChangePasswordLink string `binding:"required"`
}
