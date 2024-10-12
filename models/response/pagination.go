package response

type StatsItem struct {
	Field  string `json:"field"`
	Result any    `json:"result"`
}

type Stats struct {
	StatsItems []StatsItem `json:"statsItems"`
	Rmk        string      `json:"rmk"`
}

type PageResult struct {
	List      any     `json:"list"`
	Total     int64   `json:"total"` // total rows
	PageNum   int     `json:"pageNum"`
	PageSize  int     `json:"pageSize"`
	StatsList []Stats `json:"statsList"`
}
