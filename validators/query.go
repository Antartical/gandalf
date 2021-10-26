package validators

// validator for pagination query
type PaginationQuery struct {
	Page     int `form:"page" binding:"omitempty,min=0"`
	PageSize int `form:"limit" binding:"omitempty,min=5,max=50"`
}
