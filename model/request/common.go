package request

import "github.com/ronannnn/infra/model"

type BatchSaveCommand[T model.Crudable] struct {
	Items []*T `json:"items"`
}

type BatchDeleteCommand struct {
	Ids []uint `json:"ids"`
}
