package response

type TotalResult struct {
	Field string  `json:"field"`
	Total float64 `json:"total"`
}

type PageResult struct {
	List         any           `json:"list"`
	Total        int64         `json:"total"` // total rows
	PageNum      int           `json:"pageNum"`
	PageSize     int           `json:"pageSize"`
	TotalResults []TotalResult `json:"totalResults"` // totals of other fields, like container count, weight, volume, etc.
}
