package page

import "math"

const (
	DefaultPageSize = 10
	MaxPageSize     = 100
	DefaultPage     = 1
)

type Pagination struct {
	CurrentPage int   `json:"current_page"` // 当前页码
	PerPage     int   `json:"per_page"`     // 当前页行数
	TotalPage   int   `json:"total_page"`   // 总页码
	Total       int64 `json:"total"`        // 总行数
}

func GetPagination(page, pageSize int, total int64) Pagination {
	return Pagination{
		Total:       total,
		CurrentPage: page,
		PerPage:     pageSize,
		TotalPage:   int(math.Ceil(float64(total) / float64(pageSize))),
	}
}
