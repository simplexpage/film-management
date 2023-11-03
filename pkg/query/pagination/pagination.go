package pagination

type Pagination struct {
	Page       int `json:"page" example:"1"`
	TotalPages int `json:"total_pages" example:"10"`
	PageSize   int `json:"page_size" example:"20"`
	TotalCount int `json:"total_count" example:"200"`
}

func NewPagination(totalCount, limit, offset int) Pagination {
	return Pagination{
		Page:       offset/limit + 1,
		PageSize:   limit,
		TotalCount: totalCount,
		TotalPages: totalCount/limit + 1,
	}
}
