package pagination

type Pagination struct {
	CurrentPage  int64 `json:"current_page"`
	PerPage      int64 `json:"per_page"`
	TotalPages   int64 `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
}
