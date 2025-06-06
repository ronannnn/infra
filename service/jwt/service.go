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

func NewService(
	accessTokenCfg *accesstoken.Cfg,
	refreshtokenCfg *refreshtoken.Cfg,
	db *gorm.DB,
	refreshtokenService refreshtoken.Service,
	refreshtokenRepo refreshtoken.Repo,
	accesstokenService accesstoken.Service,
) *Service {
	return &Service{
		accessTokenCfg:      accessTokenCfg,
		refreshtokenCfg:     refreshtokenCfg,
		db:                  db,
		refreshtokenService: refreshtokenService,
		refreshtokenRepo:    refreshtokenRepo,
		accesstokenService:  accesstokenService,
	}
}

type Service struct {
	accessTokenCfg      *accesstoken.Cfg
	refreshtokenCfg     *refreshtoken.Cfg
	db                  *gorm.DB
	refreshtokenService refreshtoken.Service
	refreshtokenRepo    refreshtoken.Repo
	accesstokenService  accesstoken.Service
}

// generate access token and refresh token， used for login
func (srv *Service) GenerateTokens(ctx context.Context, claims refreshtoken.BaseClaims, userAgent string, deviceId string) (refreshToken string, accessToken string, dupLogin bool, err error) {
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

// update access token and refresh token
func (srv *Service) UpdateTokens(ctx context.Context, refreshToken string, userAgent string, deviceId string) (newRefreshToken string, accessToken string, dupLogin bool, err error) {
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
func (srv *Service) DeleteTokenByUserIdAndLoginDeviceType(ctx context.Context, userId uint, loginDeviceType useragent.DeviceType) error {
	return srv.refreshtokenRepo.Delete(ctx, srv.db, userId, loginDeviceType)
}
