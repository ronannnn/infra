package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
)

type CrudRouter[T model.Crudable] interface {
	Register(r chi.Router)
}

type DefaultCrudRouter[T model.Crudable] struct {
	Handler                  handler.CrudHandler[T]
	BasePath                 string
	ExtraRoutesUnderBasePath func(r chi.Router) // 允许自定义路由
}

func (c *DefaultCrudRouter[T]) Register(r chi.Router) {
	r.Route(c.BasePath, func(r chi.Router) {
		// batch operations
		r.Post("/batch-delete", c.Handler.BatchDelete)
		r.Post("/list", c.Handler.List)

		// single item operations
		r.Post("/", c.Handler.Create)
		r.Put("/", c.Handler.Update)

		// single item operations by ID
		r.Route("/{id}", func(r chi.Router) {
			r.Delete("/", c.Handler.DeleteById)
			r.Get("/", c.Handler.GetById)
		})

		c.ExtraRoutesUnderBasePath(r)
	})
}
