package role

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
)

type Service interface {
	Create(model *models.Role) error
	Update(partialUpdatedModel *models.Role) (models.Role, error)
	DeleteById(id uint) error
	DeleteByIds(ids []uint) error
	List(query query.RoleQuery) (response.PageResult, error)
	GetById(id uint) (models.Role, error)
}

func ProvideService(
	store Store,
) Service {
	return &ServiceImpl{
		store: store,
	}
}

type ServiceImpl struct {
	store Store
}

func (srv *ServiceImpl) Create(model *models.Role) (err error) {
	return srv.store.create(model)
}

func (srv *ServiceImpl) Update(partialUpdatedModel *models.Role) (models.Role, error) {
	return srv.store.update(partialUpdatedModel)
}

func (srv *ServiceImpl) DeleteById(id uint) error {
	return srv.store.deleteById(id)
}

func (srv *ServiceImpl) DeleteByIds(ids []uint) error {
	return srv.store.deleteByIds(ids)
}

func (srv *ServiceImpl) List(query query.RoleQuery) (response.PageResult, error) {
	return srv.store.list(query)
}

func (srv *ServiceImpl) GetById(id uint) (models.Role, error) {
	return srv.store.getById(id)
}
