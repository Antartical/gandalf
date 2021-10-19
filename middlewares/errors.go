package middlewares

// This error will be returned when a controller  tries to get authorization
// user before calling the middleware handler function.
type AuthBearerMiddlewareNotCalledError struct{}

func (e AuthBearerMiddlewareNotCalledError) Error() string {
	return "You must set the `Authorize` middleware handler"
}
