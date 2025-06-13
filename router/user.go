package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/ronannnn/infra/handler/biz"
)

func NewUserRouter(
	handler *biz.UserHandler,
) *UserRouter {
	return &UserRouter{
		handler:  handler,
		basePath: "/users",
	}
}

type UserRouter struct {
	handler  *biz.UserHandler
	basePath string
}

func (c *UserRouter) Register(r chi.Router) {
	r.Route(c.basePath, func(r chi.Router) {
		// batch operations
		r.Post("/batch-delete", c.handler.BatchDelete)
		r.Post("/list", c.handler.List)

		// single item operations
		r.Post("/", c.handler.Create)
		r.Put("/", c.handler.Update)

		// single item operations by ID
		r.Route("/{id}", func(r chi.Router) {
			r.Delete("/", c.handler.DeleteById)
			r.Get("/", c.handler.GetById)
		})
	})
}
