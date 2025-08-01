package handler

import (
	"net/http"

	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/model/request"
	"github.com/ronannnn/infra/model/request/query"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"github.com/ronannnn/infra/service"
)

type CrudHandler[T model.Crudable] interface {
	Create(w http.ResponseWriter, r *http.Request)
	BatchCreate(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	BatchUpdate(w http.ResponseWriter, r *http.Request)
	DeleteById(w http.ResponseWriter, r *http.Request)
	BatchDelete(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	GetById(w http.ResponseWriter, r *http.Request)
}

type DefaultCrudHandler[T model.Crudable] struct {
	H   HttpHandler
	Srv service.CrudSrv[T]
}

func (c *DefaultCrudHandler[T]) Create(w http.ResponseWriter, r *http.Request) {
	var newModel T
	if c.H.BindAndCheck(w, r, &newModel) {
		return
	}
	newModel = newModel.WithOprFromReq(r).(T)
	if err := c.Srv.Create(r.Context(), &newModel); err != nil {
		c.H.Fail(w, r, err, nil)
	} else {
		c.H.Success(w, r, msg.New(reason.SuccessToCreate), newModel)
	}
}

func (c *DefaultCrudHandler[T]) BatchCreate(w http.ResponseWriter, r *http.Request) {
	var cmd request.BatchSaveCommand[T]
	if c.H.BindAndCheck(w, r, &cmd) {
		return
	}
	for i := range cmd.Items {
		itemWithOpr := (*cmd.Items[i]).WithOprFromReq(r).(T)
		cmd.Items[i] = &itemWithOpr
	}
	if err := c.Srv.BatchCreate(r.Context(), cmd.Items); err != nil {
		c.H.Fail(w, r, err, nil)
	} else {
		c.H.Success(w, r, msg.New(reason.SuccessToCreate), cmd.Items)
	}
}

func (c *DefaultCrudHandler[T]) Update(w http.ResponseWriter, r *http.Request) {
	var partialUpdatedModel T
	if c.H.BindAndCheckPartial(w, r, &partialUpdatedModel) {
		return
	}
	partialUpdatedModel = partialUpdatedModel.WithUpdaterFromReq(r).(T)
	if updatedModel, err := c.Srv.Update(r.Context(), &partialUpdatedModel); err != nil {
		c.H.Fail(w, r, err, nil)
	} else {
		c.H.Success(w, r, msg.New(reason.SuccessToUpdate), updatedModel)
	}
}

func (c *DefaultCrudHandler[T]) BatchUpdate(w http.ResponseWriter, r *http.Request) {
	var cmd request.BatchSaveCommand[T]
	if c.H.BindAndCheck(w, r, &cmd) {
		return
	}
	for i := range cmd.Items {
		itemWithOpr := (*cmd.Items[i]).WithUpdaterFromReq(r).(T)
		cmd.Items[i] = &itemWithOpr
	}
	var err error
	if cmd.Items, err = c.Srv.BatchUpdate(r.Context(), cmd.Items); err != nil {
		c.H.Fail(w, r, err, nil)
	} else {
		c.H.Success(w, r, msg.New(reason.SuccessToUpdate), cmd.Items)
	}
}

func (c *DefaultCrudHandler[T]) DeleteById(w http.ResponseWriter, r *http.Request) {
	var id uint64
	if c.H.BindUint64Param(w, r, "id", &id) {
		return
	}
	if err := c.Srv.DeleteById(r.Context(), uint(id)); err != nil {
		c.H.Fail(w, r, err, nil)
	} else {
		c.H.Success(w, r, msg.New(reason.SuccessToDelete), nil)
	}
}

func (c *DefaultCrudHandler[T]) BatchDelete(w http.ResponseWriter, r *http.Request) {
	var idsReq request.BatchDeleteCommand
	if c.H.BindAndCheck(w, r, &idsReq) {
		return
	}
	if err := c.Srv.DeleteByIds(r.Context(), idsReq.Ids); err != nil {
		c.H.Fail(w, r, err, nil)
	} else {
		c.H.Success(w, r, msg.New(reason.SuccessToDelete), nil)
	}
}

func (c *DefaultCrudHandler[T]) List(w http.ResponseWriter, r *http.Request) {
	var query query.Query
	if c.H.BindAndCheck(w, r, &query) {
		return
	}
	if pageResult, err := c.Srv.List(r.Context(), query); err != nil {
		c.H.Fail(w, r, err, nil)
	} else {
		c.H.SuccessSilently(w, r, nil, pageResult)
	}
}

func (c *DefaultCrudHandler[T]) GetById(w http.ResponseWriter, r *http.Request) {
	var id uint64
	if c.H.BindUint64Param(w, r, "id", &id) {
		return
	}
	if model, err := c.Srv.GetById(r.Context(), uint(id)); err != nil {
		c.H.Fail(w, r, err, nil)
	} else {
		c.H.SuccessSilently(w, r, nil, model)
	}
}
