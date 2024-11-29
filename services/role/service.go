package role

import (
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"gorm.io/gorm"
)

type Service interface {
	Create(model *models.Role) error
	Update(partialUpdatedModel *models.Role) (models.Role, error)
	DeleteById(id uint) error
	DeleteByIds(ids []uint) error
	List(query query.Query) (response.PageResult, error)
	GetById(id uint) (models.Role, error)
}

func ProvideService(
	db *gorm.DB,
	store Store,
) Service {
	return &ServiceImpl{
		db:    db,
		store: store,
	}
}

type ServiceImpl struct {
	db    *gorm.DB
	store Store
}

func (srv *ServiceImpl) Create(model *models.Role) (err error) {
	return srv.store.Create(srv.db, model)
}

func (srv *ServiceImpl) Update(partialUpdatedModel *models.Role) (models.Role, error) {
	return srv.store.Update(srv.db, partialUpdatedModel)
}

func (srv *ServiceImpl) DeleteById(id uint) error {
	return srv.store.DeleteById(srv.db, id)
}

func (srv *ServiceImpl) DeleteByIds(ids []uint) error {
	return srv.store.DeleteByIds(srv.db, ids)
}

func (srv *ServiceImpl) List(query query.Query) (response.PageResult, error) {
	return srv.store.List(srv.db, query)
}

func (srv *ServiceImpl) GetById(id uint) (models.Role, error) {
	return srv.store.GetById(srv.db, id)
}
