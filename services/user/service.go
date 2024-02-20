package user

import (
	"context"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/models/response"
	"github.com/ronannnn/infra/services/auth"
	"gorm.io/gorm"
)

type Service interface {
	Create(context.Context, *models.User) error
	Update(ctx context.Context, partialUpdatedModel *models.User) (models.User, error)
	DeleteById(ctx context.Context, id uint) error
	DeleteByIds(ctx context.Context, ids []uint) error
	List(ctx context.Context, query query.UserQuery) (response.PageResult, error)
	GetById(ctx context.Context, id uint) (models.User, error)
	LoginByUsername(ctx context.Context, username, password string) (*AuthResult, error)
	ChangePwd(ctx context.Context, cmd ChangeUserLoginPwdCommand) error
}

func ProvideService(
	defaultHashedPassword *string,
	store Store,
	authService auth.Service,
) Service {
	return &ServiceImpl{
		defaultHashedPassword: defaultHashedPassword,
		store:                 store,
		authService:           authService,
	}
}

type ServiceImpl struct {
	defaultHashedPassword *string
	store                 Store
	authService           auth.Service
}

func (srv *ServiceImpl) Create(ctx context.Context, model *models.User) (err error) {
	if model.Password == nil {
		model.Password = srv.defaultHashedPassword
	}
	return srv.store.Create(model)
}

func (srv *ServiceImpl) Update(ctx context.Context, partialUpdatedModel *models.User) (models.User, error) {
	return srv.store.Update(partialUpdatedModel)
}

func (srv *ServiceImpl) DeleteById(ctx context.Context, id uint) error {
	return srv.store.DeleteById(id)
}

func (srv *ServiceImpl) DeleteByIds(ctx context.Context, ids []uint) error {
	return srv.store.DeleteByIds(ids)
}

func (srv *ServiceImpl) List(ctx context.Context, query query.UserQuery) (response.PageResult, error) {
	return srv.store.List(query)
}

func (srv *ServiceImpl) GetById(ctx context.Context, id uint) (models.User, error) {
	return srv.store.GetById(id)
}

func (srv *ServiceImpl) LoginByUsername(ctx context.Context, username, password string) (resp *AuthResult, err error) {
	var user models.User
	if user, err = srv.store.GetByUsername(username); err == gorm.ErrRecordNotFound {
		return nil, models.ErrWrongUsernameOrPassword
	} else if err != nil {
		return
	}
	if !CheckPassword(*user.Password, password) {
		return nil, models.ErrWrongUsernameOrPassword
	}
	var refreshToken, accessToken string
	if refreshToken, accessToken, err = srv.authService.GenerateTokens(ctx, models.BaseClaims{
		UserId:   user.Id,
		Username: *user.Username,
	}); err != nil {
		return
	}
	return &AuthResult{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, err
}

func (srv *ServiceImpl) ChangePwd(ctx context.Context, cmd ChangeUserLoginPwdCommand) (err error) {
	var user models.User
	if user, err = srv.GetById(ctx, cmd.UserId); err != nil {
		return
	}
	if !CheckPassword(*user.Password, cmd.OldPwd) {
		return models.ErrWrongUsernameOrPassword
	}
	var hashedNewPwd string
	if hashedNewPwd, err = HashPassword(cmd.NewPwd); err != nil {
		return
	}
	return srv.store.ChangePwd(cmd.UserId, hashedNewPwd)
}
