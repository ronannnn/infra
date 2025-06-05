package request

type BatchDeleteCommand struct {
	Ids []uint `json:"ids"`
}
