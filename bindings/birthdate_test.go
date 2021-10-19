package bindings

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBirthdate(t *testing.T) {
	assert := require.New(t)

	t.Run("Test binding interface", func(t *testing.T) {
		var birthdate BirthDate
		date := "1997-12-21"
		birthdate.UnmarshalJSON([]byte(date))

		expectedMarshaledDate := "\"1997-12-21\""
		marshaledDate, _ := birthdate.MarshalJSON()

		assert.Equal(date, birthdate.Format("2006-01-02"))
		assert.Equal(expectedMarshaledDate, string(marshaledDate))
	})

	t.Run("Test binding error", func(t *testing.T) {
		var birthdate BirthDate
		date := "fake format"
		err := birthdate.UnmarshalJSON([]byte(date))

		assert.Error(err)
	})
}
