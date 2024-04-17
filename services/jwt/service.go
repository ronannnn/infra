package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/services/jwt/accesstoken"
	"github.com/ronannnn/infra/services/jwt/refreshtoken"
	"gorm.io/gorm"
)

var (
	ErrInvalidTokenUserIdNotFound = fmt.Errorf("invalid token, user id not found")
)

type Service interface {
	// generate access token and refresh token， used for login
	GenerateTokens(ctx context.Context, claims refreshtoken.BaseClaims) (refreshToken string, accessToken string, err error)
	// update access token and refresh token
	UpdateTokens(ctx context.Context, refreshToken string) (newRefreshToken string, accessToken string, err error)
	// disable refresh token
	DisableRefreshToken(ctx context.Context, tokenWithPrefix string) error
	// disable refresh token by user id
	DisableRefreshTokenByUserId(ctx context.Context, userId uint) error
}

func ProvideService(
	cfg *cfg.Auth,
	db *gorm.DB,
	refreshtokenService refreshtoken.Service,
	refreshtokenStore refreshtoken.Store,
	accesstokenService accesstoken.Service,
) Service {
	return &ServiceImpl{
		cfg:                 cfg,
		db:                  db,
		refreshtokenService: refreshtokenService,
		refreshtokenStore:   refreshtokenStore,
		accesstokenService:  accesstokenService,
	}
}

type ServiceImpl struct {
	cfg                 *cfg.Auth
	db                  *gorm.DB
	refreshtokenService refreshtoken.Service
	refreshtokenStore   refreshtoken.Store
	accesstokenService  accesstoken.Service
}

func (srv *ServiceImpl) GenerateTokens(ctx context.Context, claims refreshtoken.BaseClaims) (refreshToken string, accessToken string, err error) {
	err = srv.db.Transaction(func(tx *gorm.DB) (err error) {
		refreshTokenClaims := claims.ToMap()
		jwtauth.SetExpiryIn(refreshTokenClaims, time.Duration(srv.cfg.RefreshTokenMinuteDuration)*time.Minute)
		if _, refreshToken, err = srv.refreshtokenService.GetJwtAuth().Encode(refreshTokenClaims); err != nil {
			return
		}
		if err = srv.refreshtokenStore.SaveTokenByUserId(ctx, tx, claims.UserId, refreshToken); err != nil {
			return
		}
		accessTokenClaims := claims.ToMap()
		jwtauth.SetExpiryIn(accessTokenClaims, time.Duration(srv.cfg.AccessTokenMinuteDuration)*time.Minute)
		_, accessToken, err = srv.accesstokenService.GetJwtAuth().Encode(accessTokenClaims)
		return
	})
	return
}

func (srv *ServiceImpl) UpdateTokens(ctx context.Context, refreshToken string) (newRefreshToken string, accessToken string, err error) {
	var token jwt.Token
	// validate refresh token
	if token, err = jwtauth.VerifyToken(srv.refreshtokenService.GetJwtAuth(), refreshToken); err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}
	// get user id from refresh token
	username, _ := token.Get("username")
	userId, exists := token.Get("userId")
	if !exists {
		return "", "", fmt.Errorf("invalid refresh token: missing userId")
	}
	// compare with it in db
	var tokenInDb string
	if tokenInDb, err = srv.refreshtokenStore.GetTokenByUserId(ctx, srv.db, uint(userId.(float64))); err != nil {
		return
	}
	if tokenInDb != refreshToken {
		return "", "", fmt.Errorf("incorrect refresh token")
	}
	return srv.GenerateTokens(ctx, refreshtoken.BaseClaims{
		UserId:   uint(userId.(float64)),
		Username: username.(string),
	})
}

func (srv *ServiceImpl) DisableRefreshToken(ctx context.Context, tokenString string) error {
	var err error
	var token jwt.Token
	if token, err = srv.refreshtokenService.GetJwtAuth().Decode(tokenString); err != nil {
		return err
	}
	userId, exists := token.Get("userId")
	if !exists {
		return ErrInvalidTokenUserIdNotFound
	}
	return srv.DisableRefreshTokenByUserId(ctx, uint(userId.(float64)))
}

// disable refresh token by user id
func (srv *ServiceImpl) DisableRefreshTokenByUserId(ctx context.Context, userId uint) error {
	return srv.refreshtokenStore.DeleteByUserId(ctx, srv.db, userId)
}
