package pagination

import "math"

type PaginationTable struct {
	Page     int
	Total    int
	PerPage  int
	LastPage int
	Data     interface{}
}

func (p *PaginationTable) PaginationLastPage(limit int) {
	p.LastPage = int(math.Ceil(float64(p.Total) / float64(limit)))
}
