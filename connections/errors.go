package connections

import "fmt"

/*
DatabaseConnectionError -> this error will be returned on connection failure.
*/
type DatabaseConnectionError struct {
	addr string
}

func (e DatabaseConnectionError) Error() string {
	return fmt.Sprintf("Error connecting database `addr`:%s", e.addr)
}
