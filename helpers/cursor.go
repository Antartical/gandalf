package helpers

// Represents a cursor over a query
type Cursor struct {
	Total        int
	Page         int
	PageSize     int
	TotalPages   int
	NextPage     *int
	PreviousPage *int
}

// Updates the cursor with the total number of objects
func (cursor *Cursor) Update(count int) {
	cursor.Total = count
	cursor.TotalPages = count / cursor.PageSize
	if cursor.Page > 0 {
		cursor.PreviousPage = new(int)
		*cursor.PreviousPage = cursor.Page - 1
	}
	if cursor.Page < cursor.TotalPages {
		cursor.NextPage = new(int)
		*cursor.NextPage = cursor.Page + 1
	}
}

// Creates a new cursor
func NewCursor(page int, pageSize int) Cursor {
	if pageSize <= 0 {
		pageSize = 5
	}
	return Cursor{
		Page:         page,
		PageSize:     pageSize,
		NextPage:     nil,
		PreviousPage: nil,
	}
}
