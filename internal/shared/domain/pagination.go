package domain

type Pagination struct {
	Limit  int64
	Page   int64
	Offset int64
}

func NewPagination(page, limit int64) *Pagination {

	if limit <= 0 {
		limit = 10
	}

	if limit >= 100 {
		limit = 100
	}

	if page <= 0 {
		page = 1
	}

	return &Pagination{
		Page:   page,
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}
