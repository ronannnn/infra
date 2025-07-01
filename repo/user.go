package repo

import (
	"context"

	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/model/request/query"
	"github.com/ronannnn/infra/model/response"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"github.com/ronannnn/infra/service"
	"gorm.io/gorm"
)

func NewUserRepo(
	menuRepo service.MenuRepo,
	roleRepo service.RoleRepo,
) service.UserRepo {
	return &userRepo{
		menuRepo: menuRepo,
		roleRepo: roleRepo,
		preloads: []string{
			"Menus",
			"Roles",
			"Roles.Menus",
			"JobTitle",
			"JobGrade",
			"Department",
		},
	}
}

type userRepo struct {
	menuRepo service.MenuRepo
	roleRepo service.RoleRepo
	preloads []string // 用于gorm的Preload
}

func (r userRepo) Create(ctx context.Context, tx *gorm.DB, model *model.User) error {
	if model == nil {
		return msg.NewError(reason.DbModelCreatedError).WithStack()
	}
	return tx.WithContext(ctx).Create(model).Error
}

func (r userRepo) CreateWithScopes(ctx context.Context, tx *gorm.DB, model *model.User) (updatedModel *model.User, err error) {
	if model == nil {
		err = msg.NewError(reason.DbModelCreatedError).WithStack()
		return
	}
	if err = tx.WithContext(ctx).Create(&model).Error; err != nil {
		return
	}
	return r.GetById(ctx, tx, (*model).GetId())
}

func (r userRepo) Update(ctx context.Context, tx *gorm.DB, partialUpdatedModel *model.User) (updatedModel *model.User, err error) {
	if partialUpdatedModel.Id == 0 {
		return nil, msg.NewError(reason.DbModelUpdatedIdCannotBeZero).WithStack()
	}
	err = tx.Transaction(func(tx2 *gorm.DB) (err error) {
		if partialUpdatedModel.Roles != nil {
			if err = tx2.Model(partialUpdatedModel).Association("Roles").Unscoped().Replace(partialUpdatedModel.Roles); err != nil {
				return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
			}
			for _, item := range partialUpdatedModel.Roles {
				if _, err = r.roleRepo.Update(ctx, tx2, item); err != nil {
					return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
				}
			}
			partialUpdatedModel.Roles = nil
		}
		if partialUpdatedModel.Menus != nil {
			if err = tx2.Model(partialUpdatedModel).Association("Menus").Unscoped().Replace(partialUpdatedModel.Menus); err != nil {
				return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
			}
			for _, item := range partialUpdatedModel.Menus {
				if _, err = r.menuRepo.Update(ctx, tx2, item); err != nil {
					return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
				}
			}
			partialUpdatedModel.Menus = nil
		}
		result := tx2.Updates(partialUpdatedModel)
		if result.Error != nil {
			return msg.NewError(reason.DatabaseError).WithError(result.Error).WithStack()
		}
		if result.RowsAffected == 0 {
			return msg.NewError(reason.DbModelAlreadyUpdatedByOthers).WithStack()
		}
		updatedModel, err = r.GetById(ctx, tx2, partialUpdatedModel.Id)
		return
	})
	return
}

func (r userRepo) DeleteById(ctx context.Context, tx *gorm.DB, id uint) error {
	return tx.WithContext(ctx).Transaction(func(tx2 *gorm.DB) (err error) {
		var t model.User
		if err = tx2.Delete(&t, "id = ?", id).Error; err != nil {
			return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
		}
		return
	})
}

func (r userRepo) DeleteByIds(ctx context.Context, tx *gorm.DB, ids []uint) (err error) {
	var t model.User
	if err = tx.WithContext(ctx).WithContext(ctx).Delete(&t, "id IN ?", ids).Error; err != nil {
		return msg.NewError(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (r userRepo) List(ctx context.Context, tx *gorm.DB, apiQuery query.Query) (result *response.PageResult, err error) {
	txWithCtx := tx.WithContext(ctx)
	var t model.User
	var total int64
	var list []*model.User
	if err = txWithCtx.Model(&t).Count(&total).Error; err != nil {
		return
	}
	queryScope, err := query.MakeConditionFromQuery(apiQuery, t)
	if err != nil {
		return nil, msg.NewError(reason.DatabaseError).WithError(err).WithStack()
	}
	if err = txWithCtx.
		Scopes(r.Preload()).
		Scopes(queryScope).
		Scopes(query.Paginate(apiQuery.Pagination)).
		Find(&list).Error; err != nil {
		return nil, msg.NewError(reason.DatabaseError).WithError(err).WithStack()
	}
	result = &response.PageResult{
		List:     list,
		Total:    total,
		PageNum:  apiQuery.Pagination.PageNum,
		PageSize: apiQuery.Pagination.PageSize,
	}
	return
}

func (r userRepo) GetById(ctx context.Context, tx *gorm.DB, id uint) (model *model.User, err error) {
	if err = tx.
		WithContext(ctx).
		Scopes(r.Preload()).
		First(&model, "id = ?", id).
		Error; err == gorm.ErrRecordNotFound {
		err = msg.NewError(reason.DbModelReadIdNotExists).WithReasonTemplateData(reason.IdTd{Id: id}).WithStack()
		return
	} else if err != nil {
		err = msg.NewError(reason.DatabaseError).WithError(err).WithStack()
		return
	}
	return
}

func (r userRepo) Preload() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, preload := range r.preloads {
			db = db.Preload(preload)
		}
		return db
	}
}

func (r userRepo) GetByUsername(ctx context.Context, tx *gorm.DB, username string) (user *model.User, err error) {
	err = tx.WithContext(ctx).Scopes(r.Preload()).First(&user, &model.User{Username: &username}).Error
	return
}

func (r userRepo) GetByNickname(ctx context.Context, tx *gorm.DB, nickname string) (user *model.User, err error) {
	err = tx.WithContext(ctx).Scopes(r.Preload()).First(&user, &model.User{Nickname: &nickname}).Error
	return
}

func (r userRepo) ChangePwd(ctx context.Context, tx *gorm.DB, userId uint, newPwd string) error {
	return tx.WithContext(ctx).Model(&model.User{}).Where("id = ?", userId).Update("password", newPwd).Error
}
