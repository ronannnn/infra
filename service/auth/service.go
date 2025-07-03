package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/ronannnn/infra/constant"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
	"github.com/ronannnn/infra/service/jwt"
	"github.com/ronannnn/infra/service/jwt/accesstoken"
	"github.com/ronannnn/infra/service/jwt/refreshtoken"
	"github.com/ronannnn/infra/service/login"
	"go.uber.org/zap"
)

func NewService(
	log *zap.SugaredLogger,
	jwtService *jwt.Service,
	accessTokenService *accesstoken.Service,
	userService *service.UserService,
) *Service {
	return &Service{
		log:                log,
		jwtService:         jwtService,
		accessTokenService: accessTokenService,
		userService:        userService,
	}
}

type Service struct {
	log                *zap.SugaredLogger
	jwtService         *jwt.Service
	accessTokenService *accesstoken.Service
	userService        *service.UserService
}

func (srv *Service) RefreshToken(ctx context.Context, oldRefreshToken string) (*login.Result, error) {
	ua := ctx.Value(constant.CtxKeyUa).(string)
	deviceId := ctx.Value(constant.CtxKeyDeviceId).(string)
	refreshToken, accessToken, dupLogin, err := srv.jwtService.UpdateTokens(ctx, oldRefreshToken, ua, deviceId)
	if err != nil {
		return nil, err
	}

	return &login.Result{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		DupLogin:     dupLogin,
	}, nil
}

func (srv *Service) GenerateLongTermAccessToken(ctx context.Context, userId uint) (accessToken string, err error) {
	if userId == 0 {
		err = fmt.Errorf("用户ID不能为空")
		return
	}
	var user *model.User
	if user, err = srv.userService.GetById(ctx, userId); err != nil {
		return
	}
	if user.Username == nil {
		err = fmt.Errorf("该ID的用户名为空")
		return
	}
	claims := refreshtoken.BaseClaims{
		UserId:   userId,
		Username: *user.Username,
	}
	mappedClaims := claims.ToMap()
	jwtauth.SetExpiryIn(mappedClaims, time.Duration(24*365*10)*time.Hour)
	_, accessToken, err = srv.accessTokenService.GetJwtAuth().Encode(mappedClaims)
	return
}
