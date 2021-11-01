package bindings

import (
	"encoding/json"
	"fmt"
	"gandalf/security"
	"strings"
)

// Not found scope error
type ScopeNotFoundError struct {
	Scope string
}

func (e ScopeNotFoundError) Error() string {
	return fmt.Sprintf("%s, does not exists", e.Scope)
}

var validScopes = map[string]bool{
	security.ScopeUserAuthorizeApp:   true,
	security.ScopeUserVerify:         true,
	security.ScopeUserChangePassword: true,
	security.ScopeUserRead:           true,
	security.ScopeUserWrite:          true,
	security.ScopeUserDelete:         true,
	security.ScopeAppRead:            true,
	security.ScopeAppWrite:           true,
}

// Scope binding
type Scope string

// Implement Unmarshaler interface
func (scope *Scope) UnmarshalJSON(b []byte) error {
	*scope = Scope(strings.Replace(string(b), "\"", "", -1))
	if !validScopes[scope.ToString()] {
		return ScopeNotFoundError{
			Scope: scope.ToString(),
		}
	}
	return nil
}

// Implement Marshaler interface
func (scope Scope) MarshalJSON() ([]byte, error) {
	return json.Marshal(scope.ToString())
}

// To string
func (scope Scope) ToString() string {
	return strings.ToLower(string(scope))
}

// Convert Scope array to String array
func ScopeArrayToStringArray(array []Scope) []string {
	var scopes []string
	for _, value := range array {
		scopes = append(scopes, value.ToString())
	}
	return scopes
}
