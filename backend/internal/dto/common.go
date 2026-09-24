package dto

// PaginationQuery 通用分页查询参数。
type PaginationQuery struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"pageSize"`
}

// PageResult 统一分页返回结构。
type PageResult struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}
