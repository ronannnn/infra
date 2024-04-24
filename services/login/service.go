package login

import (
	"context"

	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/services/jwt"
	"github.com/ronannnn/infra/services/jwt/refreshtoken"
	"github.com/ronannnn/infra/services/user"
	"github.com/ronannnn/infra/utils/useragent"
	"gorm.io/gorm"
)

type Service interface {
	LoginByUsername(ctx context.Context, cmd UsernameCmd) (*Result, error)
	Logout(ctx context.Context, userId uint, userAgent string) error
	ChangePwd(ctx context.Context, cmd ChangeUserPwdCmd) error
}

func ProvideService(
	db *gorm.DB,
	store user.Store,
	jwtService jwt.Service,
) Service {
	return &ServiceImpl{
		db:         db,
		store:      store,
		jwtService: jwtService,
	}
}

type ServiceImpl struct {
	db         *gorm.DB
	store      user.Store
	jwtService jwt.Service
}

func (srv *ServiceImpl) LoginByUsername(ctx context.Context, cmd UsernameCmd) (resp *Result, err error) {
	var user models.User
	if user, err = srv.store.GetByUsername(srv.db, cmd.Username); err == gorm.ErrRecordNotFound {
		return nil, models.ErrWrongUsernameOrPassword
	} else if err != nil {
		return
	}
	if !CheckPassword(*user.Password, cmd.Password) {
		return nil, models.ErrWrongUsernameOrPassword
	}
	var refreshToken, accessToken string
	var dupLogin bool
	if refreshToken, accessToken, dupLogin, err = srv.jwtService.GenerateTokens(ctx, refreshtoken.BaseClaims{
		UserId:   user.Id,
		Username: *user.Username,
	}, cmd.UserAgent, cmd.DeviceId); err != nil {
		return
	}
	return &Result{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		DupLogin:     dupLogin,
	}, err
}

func (srv *ServiceImpl) Logout(ctx context.Context, userId uint, userAgent string) (err error) {
	return srv.jwtService.DeleteTokenByUserIdAndLoginDeviceType(ctx, userId, useragent.Parse(userAgent).DeviceType())
}

func (srv *ServiceImpl) ChangePwd(ctx context.Context, cmd ChangeUserPwdCmd) (err error) {
	var user models.User
	if user, err = srv.store.GetById(srv.db, cmd.UserId); err != nil {
		return
	}
	if !CheckPassword(*user.Password, cmd.OldPwd) {
		return models.ErrWrongUsernameOrPassword
	}
	var hashedNewPwd string
	if hashedNewPwd, err = HashPassword(cmd.NewPwd); err != nil {
		return
	}
	return srv.store.ChangePwd(srv.db, cmd.UserId, hashedNewPwd)
}
