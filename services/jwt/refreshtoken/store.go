package refreshtoken

import (
	"context"

	"github.com/ronannnn/infra/models"
	"gorm.io/gorm"
)

type Store interface {
	Create(context.Context, *gorm.DB, *RefreshToken) error
	Update(context.Context, *gorm.DB, *RefreshToken) (*RefreshToken, error)
	SaveTokenByUserId(ctx context.Context, tx *gorm.DB, userId uint, newRefreshToken string) error
	DeleteByUserId(ctx context.Context, tx *gorm.DB, userId uint) error
	GetTokenByUserId(ctx context.Context, tx *gorm.DB, userId uint) (string, error)
}

func ProvideStore() Store {
	return &StoreImpl{}
}

type StoreImpl struct {
}

func (s *StoreImpl) Create(ctx context.Context, tx *gorm.DB, model *RefreshToken) error {
	return tx.Create(model).Error
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

func (s *StoreImpl) SaveTokenByUserId(ctx context.Context, tx *gorm.DB, userId uint, newRefreshToken string) error {
	return tx.
		Where(RefreshToken{UserId: &userId}).
		Assign(RefreshToken{RefreshToken: &newRefreshToken}).
		FirstOrCreate(&RefreshToken{
			UserId:       &userId,
			RefreshToken: &newRefreshToken,
		}).Error
}

func (s *StoreImpl) DeleteByUserId(ctx context.Context, tx *gorm.DB, id uint) error {
	return tx.Delete(&RefreshToken{}, "user_id = ?", id).Error
}

func (s *StoreImpl) GetTokenByUserId(ctx context.Context, tx *gorm.DB, userId uint) (refreshToken string, err error) {
	var model RefreshToken
	if err = tx.First(&model, "user_id = ?", userId).Error; err != nil {
		return
	}
	return *model.RefreshToken, nil
}
