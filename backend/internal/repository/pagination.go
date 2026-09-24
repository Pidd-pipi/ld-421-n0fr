package repository

// Pagination 通用分页参数，page 从 1 开始。
type Pagination struct {
	Page     int
	PageSize int
}

// Normalize 规整分页参数，返回安全值。
func (p *Pagination) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// Offset 返回查询偏移量。
func (p Pagination) Offset() int { return (p.Page - 1) * p.PageSize }
