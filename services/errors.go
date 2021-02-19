package services

/*
AuthenticationError -> this error will be returned on user authentication
failure
*/
type AuthenticationError struct {
	raisedFrom error
}

func (e AuthenticationError) Error() string {
	return "User cannot be authenticate"
}

/*
AuthorizationError -> this error will be returned on user authorization
failure
*/
type AuthorizationError struct {
	raisedFrom error
}

func (e AuthorizationError) Error() string {
	return "User cannot be authorized"
}

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

/*
UserNotFoundError -> this error will be returned on user not found
exception
*/
type UserNotFoundError struct {
	raisedFrom error
}

func (e UserNotFoundError) Error() string {
	return "User not found"
}
