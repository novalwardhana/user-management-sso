package pagination

import "math"

type PaginationTable struct {
	Page        int         `json:"page"`
	DataPerPage int         `json:"data_per_page"`
	TotalData   int         `json:"total_data"`
	TotalPage   int         `json:"total_page"`
	Data        interface{} `json:"data"`
}

func (p *PaginationTable) PaginationTotalPage() {
	p.TotalPage = int(math.Ceil(float64(p.TotalData) / float64(p.DataPerPage)))
}
