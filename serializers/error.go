package serializers

// Struct for http error serialization
type HTTPErrorSerializer struct {
	Code  int    `json:"code" example:"400"`
	Error string `json:"error" example:"status bad request"`
}

// Creates a new http error serializer
func NewHTTPErrorSerializer(status int, err error) HTTPErrorSerializer {
	return HTTPErrorSerializer{
		Code:  status,
		Error: err.Error(),
	}
}
