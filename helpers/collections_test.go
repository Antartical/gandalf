package helpers

import (
	"fmt"
	"gandalf/models"
	"gandalf/tests"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestPqStringArrayContains(t *testing.T) {
	assert := require.New(t)

	t.Run("Test contains", func(t *testing.T) {
		element := "Me"
		pqArray := pq.StringArray{"wowowo", element}

		assert.True(PqStringArrayContains(pqArray, element))
	})

	t.Run("Test not contains", func(t *testing.T) {
		element := "Me"
		pqArray := pq.StringArray{"wowowo"}
		assert.False(PqStringArrayContains(pqArray, element))
	})
}

func TestDBPaginate(t *testing.T) {
	assert := require.New(t)

	t.Run("Test db paginate", func(t *testing.T) {
		var users []models.User
		page := 0
		pageSize := 30
		db := tests.NewTestDatabase(false)
		tx := db.Scopes(DBPaginate(page, pageSize)).Find(&users)
		raw := fmt.Sprint(tx.Statement.Clauses["LIMIT"].Expression)
		assert.Contains(raw, fmt.Sprint(page))
		assert.Contains(raw, fmt.Sprint(pageSize))
	})
}
