package services

import "fmt"

// This error will be returned on user authentication failure
type AuthenticationError struct {
	raisedFrom error
}

func (e AuthenticationError) Error() string {
	return "User cannot be authenticate"
}

// This error will be returned on user authorization failure
type AuthorizationError struct {
	raisedFrom error
}

func (e AuthorizationError) Error() string {
	return "User cannot be authorized"
}

// This error will be returned on user creation failure
type UserCreateError struct {
	raisedFrom error
}

func (e UserCreateError) Error() string {
	return "User email already registered"
}

// This error will be returned on user creation failure.
type AppCreateError struct {
	raisedFrom error
}

func (e AppCreateError) Error() string {
	return "App cannot be created"
}

// This error will be returned on app not found exception
type AppNotFoundError struct {
	raisedFrom error
}

func (e AppNotFoundError) Error() string {
	return "App not found"
}

// This error will be returned on user not found exception
type UserNotFoundError struct {
	raisedFrom error
}

func (e UserNotFoundError) Error() string {
	return "User not found"
}

// Error for app authorization on unknown redirect uir
type RedirectUriDoesNotMatch struct {
	raisedFrom  error
	redirectUri string
}

func (e RedirectUriDoesNotMatch) Error() string {
	return fmt.Sprintf("Redirect uri is not registered for the app, %s", e.redirectUri)
}
