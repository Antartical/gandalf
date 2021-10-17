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
AppCreateError -> this error will be returned on user creation
failure.
*/
type AppCreateError struct {
	raisedFrom error
}

func (e AppCreateError) Error() string {
	return "App cannot be created"
}

/*
AppNotFoundError -> this error will be returned on user not found
exception
*/
type AppNotFoundError struct {
	raisedFrom error
}

func (e AppNotFoundError) Error() string {
	return "App not found"
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
