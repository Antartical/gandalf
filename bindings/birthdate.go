package bindings

import (
	"encoding/json"
	"strings"
	"time"
)

// Birthdate JSON binding
type BirthDate time.Time

// Implement Unmarshaler interface
func (j *BirthDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = BirthDate(t)
	return nil
}

// Implement Marshaler interface
func (j BirthDate) MarshalJSON() ([]byte, error) {
	formattedDate := j.Format("2006-01-02")
	return json.Marshal(formattedDate)
}

// Format function for printing your date
func (j BirthDate) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}
