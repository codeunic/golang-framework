package framework

type Pagination struct {
	Page     int64 `form:"page"`
	PageSize int64 `form:"page_size"`
}

func (p Pagination) GetLimit() int64 {
	return p.GetOffset() + p.PageSize
}

func (p Pagination) GetOffset() int64 {
	return (p.Page - 1) * p.PageSize
}

func NewPagination(page int64, pageSize int64) *Pagination {
	return &Pagination{Page: page, PageSize: pageSize}
}

type ResponsePagination struct {
	Pagination
}
