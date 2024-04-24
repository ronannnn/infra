package refreshtoken

import (
	"context"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/utils/useragent"
	"gorm.io/gorm"
)

type Store interface {
	Save(context.Context, *gorm.DB, *RefreshToken) (bool, error)
	Update(context.Context, *gorm.DB, *RefreshToken) (*RefreshToken, error)
	Delete(ctx context.Context, tx *gorm.DB, userId uint, loginDeviceType useragent.DeviceType) error
	Get(ctx context.Context, tx *gorm.DB, userId uint, loginDeviceType useragent.DeviceType) (string, error)
}

func ProvideStore() Store {
	return &StoreImpl{}
}

type StoreImpl struct {
}

// Save 保存RefreshToken
//
//	 检查是否存在相同的UserId，LoginDeviceType
//	 1. 如果存在，更新RefreshToken和deviceId
//			1.1. 旧的DeviceId和本次更新的DeviceId不同，说明是重复登录，dupLogin返回true
//	    1.2. 旧的DeviceId和本次更新的DeviceId相同，dupLogin返回false
//	 2. 如果不存在，则插入新的RefreshToken
func (s *StoreImpl) Save(ctx context.Context, tx *gorm.DB, refreshToken *RefreshToken) (dupLogin bool, err error) {
	var dbRefreshToken *RefreshToken
	if err = tx.
		Where(RefreshToken{UserId: refreshToken.UserId, LoginDeviceType: refreshToken.LoginDeviceType}).
		First(&dbRefreshToken).
		Error; err == gorm.ErrRecordNotFound {
		err = tx.Create(refreshToken).Error
		return
	} else if err != nil {
		return
	}
	// 两个都没有deviceId，或者两个都有deviceId，但是不相等，都是重复登录
	if (dbRefreshToken.DeviceId == nil && refreshToken.DeviceId == nil) ||
		(dbRefreshToken.DeviceId != nil && refreshToken.DeviceId != nil && *dbRefreshToken.DeviceId != *refreshToken.DeviceId) {
		dupLogin = true
	}
	err = tx.Model(&dbRefreshToken).Updates(RefreshToken{
		RefreshToken: refreshToken.RefreshToken,
		DeviceId:     refreshToken.DeviceId,
	}).Error
	return
}

func (s *StoreImpl) Update(ctx context.Context, tx *gorm.DB, partialUpdatedModel *RefreshToken) (updatedModel *RefreshToken, err error) {
	if partialUpdatedModel.Id == 0 {
		return updatedModel, models.ErrUpdatedId
	}
	result := tx.Updates(partialUpdatedModel)
	if result.Error != nil {
		return updatedModel, result.Error
	}
	if result.RowsAffected == 0 {
		return updatedModel, models.ErrModified("RefreshToken")
	}
	err = tx.First(&updatedModel, "id = ?", partialUpdatedModel.Id).Error
	return
}

func (s *StoreImpl) Delete(ctx context.Context, tx *gorm.DB, userId uint, loginDeviceType useragent.DeviceType) error {
	return tx.
		Where(&RefreshToken{
			UserId:          &userId,
			LoginDeviceType: &loginDeviceType,
		}).
		Delete(&RefreshToken{}).
		Error
}

func (s *StoreImpl) Get(ctx context.Context, tx *gorm.DB, userId uint, loginDeviceType useragent.DeviceType) (refreshToken string, err error) {
	var model RefreshToken
	if err = tx.
		Where(&RefreshToken{
			UserId:          &userId,
			LoginDeviceType: &loginDeviceType,
		}).
		First(&model).
		Error; err != nil {
		return
	}
	return *model.RefreshToken, nil
}
