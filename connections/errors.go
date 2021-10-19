package connections

import "fmt"

// Represents an error that happens during the database connection
type DatabaseConnectionError struct {
	addr string
}

func (e DatabaseConnectionError) Error() string {
	return fmt.Sprintf("Error connecting database `addr`:%s", e.addr)
}
