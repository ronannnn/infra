package refreshtoken

import (
	"context"

	"github.com/go-chi/jwtauth/v5"
	"github.com/ronannnn/infra/utils/useragent"
	"gorm.io/gorm"
)

type Repo interface {
	Save(context.Context, *gorm.DB, *RefreshToken) (bool, error)
	Update(context.Context, *gorm.DB, *RefreshToken) (*RefreshToken, error)
	Delete(ctx context.Context, tx *gorm.DB, userId uint, loginDeviceType useragent.DeviceType) error
	Get(ctx context.Context, tx *gorm.DB, userId uint, loginDeviceType useragent.DeviceType) (string, error)
}

func NewService(
	cfg *Cfg,
) *Service {
	return &Service{
		jwtAuth: jwtauth.New("HS256", []byte(cfg.RefreshTokenSecret), nil),
	}
}

type Service struct {
	jwtAuth *jwtauth.JWTAuth
}

func (srv *Service) GetJwtAuth() *jwtauth.JWTAuth {
	return srv.jwtAuth
}
