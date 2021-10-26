package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCursor(t *testing.T) {
	assert := require.New(t)

	t.Run("Test Cursor update", func(t *testing.T) {
		page := 3
		pageSize := 5
		expectedPrevious := 2
		expectedNext := 4
		expectedTotal := 100
		expectedNumberOfPages := 20
		cursor := NewCursor(page, pageSize)

		cursor.Update(expectedTotal)

		assert.Equal(page, cursor.Page)
		assert.Equal(pageSize, cursor.PageSize)
		assert.Equal(expectedPrevious, *cursor.PreviousPage)
		assert.Equal(expectedNext, *cursor.NextPage)
		assert.Equal(expectedNumberOfPages, cursor.TotalPages)
	})
}
