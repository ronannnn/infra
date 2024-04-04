package response

type StatsItem struct {
	Field  string `json:"field"`
	Result any    `json:"result"`
}

type PageResult struct {
	List     any         `json:"list"`
	Total    int64       `json:"total"` // total rows
	PageNum  int         `json:"pageNum"`
	PageSize int         `json:"pageSize"`
	Stats    []StatsItem `json:"stats"` // totals of other fields, like container count, weight, volume, etc.
}
