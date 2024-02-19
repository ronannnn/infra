package apirecord

import "github.com/ronannnn/infra/models"

type Service interface {
	Create(model *models.ApiRecord) error
	Update(partialUpdatedModel *models.ApiRecord) (models.ApiRecord, error)
	DeleteById(id uint) error
	DeleteByIds(ids []uint) error
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

func (srv *ServiceImpl) Create(model *models.ApiRecord) error {
	return srv.store.create(model)
}

func (srv *ServiceImpl) Update(partialUpdatedModel *models.ApiRecord) (updatedModel models.ApiRecord, err error) {
	return srv.store.update(partialUpdatedModel)
}

func (srv *ServiceImpl) DeleteById(id uint) error {
	return srv.store.deleteById(id)
}

func (srv *ServiceImpl) DeleteByIds(ids []uint) error {
	return srv.store.deleteByIds(ids)
}
