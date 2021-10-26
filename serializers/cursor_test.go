package serializers

import (
	"gandalf/helpers"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCursorSerializer(t *testing.T) {
	assert := require.New(t)

	t.Run("Test cursor serializer", func(t *testing.T) {
		page := 2
		pageSize := 3
		cursor := helpers.NewCursor(page, pageSize)

		serializedCursor := NewCursorSerializer(cursor)
		assert.Equal(page, serializedCursor.Data.ActualPage)
		assert.Equal(pageSize, serializedCursor.Data.PageSize)
		assert.Equal(0, serializedCursor.Data.TotalPages)
		assert.Equal(0, serializedCursor.Data.TotalObjects)
		assert.Zero(*serializedCursor.Data.PreviousPage)
		assert.Zero(*serializedCursor.Data.NextPage)
	})
}
