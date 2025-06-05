package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/ronannnn/infra/service/jwt/accesstoken"
	"github.com/ronannnn/infra/service/jwt/refreshtoken"
	"github.com/ronannnn/infra/utils/useragent"
	"gorm.io/gorm"
)

var (
	ErrInvalidTokenUserIdNotFound = fmt.Errorf("invalid token, user id not found")
)

type Service interface {
	// generate access token and refresh token， used for login
	GenerateTokens(ctx context.Context, claims refreshtoken.BaseClaims, userAgent string, deviceId string) (refreshToken string, accessToken string, dupLogin bool, err error)
	// update access token and refresh token
	UpdateTokens(ctx context.Context, refreshToken string, userAgent string, deviceId string) (newRefreshToken string, accessToken string, dupLogin bool, err error)
	// disable refresh token
	DeleteTokenByUserIdAndLoginDeviceType(ctx context.Context, userId uint, loginDeviceType useragent.DeviceType) error
}

func ProvideService(
	accessTokenCfg *accesstoken.Cfg,
	refreshtokenCfg *refreshtoken.Cfg,
	db *gorm.DB,
	refreshtokenService refreshtoken.Service,
	refreshtokenRepo refreshtoken.Repo,
	accesstokenService accesstoken.Service,
) Service {
	return &ServiceImpl{
		accessTokenCfg:      accessTokenCfg,
		refreshtokenCfg:     refreshtokenCfg,
		db:                  db,
		refreshtokenService: refreshtokenService,
		refreshtokenRepo:    refreshtokenRepo,
		accesstokenService:  accesstokenService,
	}
}

type ServiceImpl struct {
	accessTokenCfg      *accesstoken.Cfg
	refreshtokenCfg     *refreshtoken.Cfg
	db                  *gorm.DB
	refreshtokenService refreshtoken.Service
	refreshtokenRepo    refreshtoken.Repo
	accesstokenService  accesstoken.Service
}

func (srv *ServiceImpl) GenerateTokens(ctx context.Context, claims refreshtoken.BaseClaims, userAgent string, deviceId string) (refreshToken string, accessToken string, dupLogin bool, err error) {
	err = srv.db.Transaction(func(tx *gorm.DB) (err error) {
		// get refresh token
		refreshTokenClaims := claims.ToMap()
		jwtauth.SetExpiryIn(refreshTokenClaims, time.Duration(srv.refreshtokenCfg.RefreshTokenMinuteDuration)*time.Minute)
		if _, refreshToken, err = srv.refreshtokenService.GetJwtAuth().Encode(refreshTokenClaims); err != nil {
			return
		}
		// get login device type
		loginDeviceType := useragent.Parse(userAgent).DeviceType()
		if dupLogin, err = srv.refreshtokenRepo.Save(ctx, tx, &refreshtoken.RefreshToken{
			UserId:          &claims.UserId,
			RefreshToken:    &refreshToken,
			LoginDeviceType: &loginDeviceType,
			DeviceId:        &deviceId,
		}); err != nil {
			return
		}
		accessTokenClaims := claims.ToMap()
		jwtauth.SetExpiryIn(accessTokenClaims, time.Duration(srv.accessTokenCfg.AccessTokenMinuteDuration)*time.Minute)
		_, accessToken, err = srv.accesstokenService.GetJwtAuth().Encode(accessTokenClaims)
		return
	})
	return
}

func (srv *ServiceImpl) UpdateTokens(ctx context.Context, refreshToken string, userAgent string, deviceId string) (newRefreshToken string, accessToken string, dupLogin bool, err error) {
	var token jwt.Token
	// validate refresh token
	if token, err = jwtauth.VerifyToken(srv.refreshtokenService.GetJwtAuth(), refreshToken); err != nil {
		return "", "", false, fmt.Errorf("invalid refresh token: %w", err)
	}
	// get user id from refresh token
	username, _ := token.Get("username")
	userId, exists := token.Get("userId")
	if !exists {
		return "", "", false, fmt.Errorf("invalid refresh token: missing userId")
	}
	// compare with it in db
	var tokenInDb string
	if tokenInDb, err = srv.refreshtokenRepo.Get(ctx, srv.db, uint(userId.(float64)), useragent.Parse(userAgent).DeviceType()); err != nil {
		return
	}
	if tokenInDb != refreshToken {
		return "", "", false, fmt.Errorf("incorrect refresh token")
	}
	return srv.GenerateTokens(ctx, refreshtoken.BaseClaims{
		UserId:   uint(userId.(float64)),
		Username: username.(string),
	}, userAgent, deviceId)
}

// disable refresh token by user id
func (srv *ServiceImpl) DeleteTokenByUserIdAndLoginDeviceType(ctx context.Context, userId uint, loginDeviceType useragent.DeviceType) error {
	return srv.refreshtokenRepo.Delete(ctx, srv.db, userId, loginDeviceType)
}
