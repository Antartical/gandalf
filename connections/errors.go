package connections

/*
PostgresConnectionError -> this error will be returned on postgres
connection fail.
*/
type PostgresConnectionError struct{}

func (e *PostgresConnectionError) Error() string {
	return "Error connecting postgres"
}
