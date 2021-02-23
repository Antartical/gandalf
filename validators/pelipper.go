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
