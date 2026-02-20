package schemas

type PaginationQuery struct {
	Limit int64 `form:"limit"`
	Page  int64 `form:"page"`
}

func (p *PaginationQuery) Normalize() {

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
