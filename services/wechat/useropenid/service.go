package useropenid

import "gorm.io/gorm"

type Service interface {
	Create(model *UserOpenId) error
	UpdateOpenIdByUserIdAndAppId(partialUpdatedModel *UserOpenId) error
	SaveIfNotExists(model *UserOpenId) error
	GetByUserIdAndAppId(userId uint, appId string) (UserOpenId, error)
}

func ProvideService(
	store Store,
	db *gorm.DB,
) Service {
	return &ServiceImpl{
		store: store,
		db:    db,
	}
}

type ServiceImpl struct {
	store Store
	db    *gorm.DB
}

func (srv *ServiceImpl) Create(model *UserOpenId) (err error) {
	return srv.store.create(srv.db, model)
}

func (srv *ServiceImpl) UpdateOpenIdByUserIdAndAppId(partialUpdatedModel *UserOpenId) error {
	return srv.store.updateOpenIdByUserIdAndAppId(srv.db, partialUpdatedModel)
}

func (srv *ServiceImpl) SaveIfNotExists(partialUpdatedModel *UserOpenId) (err error) {
	if _, err = srv.GetByUserIdAndAppId(partialUpdatedModel.UserId, partialUpdatedModel.OfficialAccountAppId); err != gorm.ErrRecordNotFound {
		return srv.Create(partialUpdatedModel)
	} else if err != nil {
		return err
	}
	return srv.store.updateOpenIdByUserIdAndAppId(srv.db, partialUpdatedModel)
}

func (srv *ServiceImpl) GetByUserIdAndAppId(userId uint, appId string) (UserOpenId, error) {
	return srv.store.getByUserIdAndAppId(srv.db, userId, appId)
}
