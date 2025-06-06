package router

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
)

type CrudRouter[T model.Crudable] interface {
	Register(r *chi.Mux)
	GetBasePath() string
}

type DefaultCrudRouter[T model.Crudable] struct {
	Handler handler.CrudHandler[T]
}

func (c *DefaultCrudRouter[T]) Register(r *chi.Mux) {
	basePath := c.GetBasePath()
	// batch operations
	r.Post(fmt.Sprintf("%s/batch-delete", basePath), c.Handler.BatchDelete)
	r.Post(fmt.Sprintf("%s/list", basePath), c.Handler.List)
	// single item operations
	r.Post(basePath, c.Handler.Create)
	r.Put(basePath, c.Handler.Update)
	r.Delete(fmt.Sprintf("%s/{id}", basePath), c.Handler.DeleteById)
	r.Get(fmt.Sprintf("%s/{id}", basePath), c.Handler.GetById)
}

func (c *DefaultCrudRouter[T]) GetBasePath() string {
	return ""
}
