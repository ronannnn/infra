package login

import (
	"context"

	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
	"github.com/ronannnn/infra/service/jwt"
	"github.com/ronannnn/infra/service/jwt/refreshtoken"
	"github.com/ronannnn/infra/utils/useragent"
	"gorm.io/gorm"
)

type Service interface {
	LoginByUsername(ctx context.Context, cmd UsernameCmd) (*Result, error)
	Logout(ctx context.Context, userId uint, userAgent string) error
	ChangePwd(ctx context.Context, cmd ChangeUserPwdCmd) error
}

func NewService(
	db *gorm.DB,
	userRepo service.UserRepo,
	jwtService jwt.Service,
) Service {
	return &ServiceImpl{
		db:         db,
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

type ServiceImpl struct {
	db         *gorm.DB
	userRepo   service.UserRepo
	jwtService jwt.Service
}

func (srv *ServiceImpl) LoginByUsername(ctx context.Context, cmd UsernameCmd) (resp *Result, err error) {
	var user *model.User
	if user, err = srv.userRepo.GetByUsername(ctx, srv.db, cmd.Username); err == gorm.ErrRecordNotFound {
		return nil, model.ErrWrongUsernameOrPassword
	} else if err != nil {
		return
	}
	if !CheckPassword(*user.Password, cmd.Password) {
		return nil, model.ErrWrongUsernameOrPassword
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
	var user *model.User
	if user, err = srv.userRepo.GetById(ctx, srv.db, cmd.UserId); err != nil {
		return
	}
	if !CheckPassword(*user.Password, cmd.OldPwd) {
		return model.ErrWrongUsernameOrPassword
	}
	var hashedNewPwd string
	if hashedNewPwd, err = HashPassword(cmd.NewPwd); err != nil {
		return
	}
	return srv.userRepo.ChangePwd(ctx, srv.db, cmd.UserId, hashedNewPwd)
}
