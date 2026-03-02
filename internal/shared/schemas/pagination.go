package schemas

type ResponsePagination struct {
	Page      int64 `json:"page"`
	PageTotal int64 `json:"page_total"`
	Limit     int64 `json:"limit"`
	Total     int64 `json:"total"`
}

type PaginationQuery struct {
	Limit int64 `form:"limit"`
	Page  int64 `form:"page"`
}

func (p *PaginationQuery) NormalizePagination() {

	if p.Limit <= 0 {
		p.Limit = 10
	}

	if p.Limit >= 100 {
		p.Limit = 100
	}

	if p.Page <= 0 {
		p.Page = 1
	}
}
