package services

/*
UserCreateError -> this error will be returned on user creation
failure.
*/
type UserCreateError struct {
	raisedFrom error
}

func (e UserCreateError) Error() string {
	return "User email already registered"
}
