package biz

import (
	"net/http"

	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/model/request"
	"github.com/ronannnn/infra/model/request/query"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"github.com/ronannnn/infra/service"
)

func NewUserHandler(
	h handler.HttpHandler,
	srv *service.UserService,
) *UserHandler {
	return &UserHandler{
		h:   h,
		srv: srv,
	}
}

type UserHandler struct {
	h   handler.HttpHandler
	srv *service.UserService
}

func (c *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var newModel model.User
	if c.h.BindAndCheck(w, r, &newModel) {
		return
	}
	newModel.GetOprFromReq(r)
	if err := c.srv.Create(r.Context(), &newModel); err != nil {
		c.h.Fail(w, r, err, nil)
	} else {
		c.h.Success(w, r, msg.New(reason.SuccessToCreate), newModel)
	}
}

func (c *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var partialUpdatedModel model.User
	if c.h.BindAndCheckPartial(w, r, &partialUpdatedModel) {
		return
	}
	partialUpdatedModel.GetUpdaterFromReq(r)
	if updatedModel, err := c.srv.Update(r.Context(), &partialUpdatedModel); err != nil {
		c.h.Fail(w, r, err, nil)
	} else {
		c.h.Success(w, r, msg.New(reason.SuccessToUpdate), updatedModel)
	}
}

func (c *UserHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	var id uint64
	if c.h.BindUint64Param(w, r, "id", &id) {
		return
	}
	if err := c.srv.DeleteById(r.Context(), uint(id)); err != nil {
		c.h.Fail(w, r, err, nil)
	} else {
		c.h.Success(w, r, msg.New(reason.SuccessToDelete), nil)
	}
}

func (c *UserHandler) BatchDelete(w http.ResponseWriter, r *http.Request) {
	var idsReq request.BatchDeleteCommand
	if c.h.BindAndCheck(w, r, &idsReq) {
		return
	}
	if err := c.srv.DeleteByIds(r.Context(), idsReq.Ids); err != nil {
		c.h.Fail(w, r, err, nil)
	} else {
		c.h.Success(w, r, msg.New(reason.SuccessToDelete), nil)
	}
}

func (c *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	var query query.Query
	if c.h.BindAndCheck(w, r, &query) {
		return
	}
	if pageResult, err := c.srv.List(r.Context(), query); err != nil {
		c.h.Fail(w, r, err, nil)
	} else {
		c.h.SuccessSilently(w, r, nil, pageResult)
	}
}

func (c *UserHandler) GetById(w http.ResponseWriter, r *http.Request) {
	var id uint64
	if c.h.BindUint64Param(w, r, "id", &id) {
		return
	}
	if model, err := c.srv.GetById(r.Context(), uint(id)); err != nil {
		c.h.Fail(w, r, err, nil)
	} else {
		c.h.SuccessSilently(w, r, nil, model)
	}
}
