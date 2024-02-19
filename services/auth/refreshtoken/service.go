package refreshtoken

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/ronannnn/infra/cfg"
)

type Service interface {
	GetJwtAuth() *jwtauth.JWTAuth
}

func ProvideService(
	cfg *cfg.Jwt,
) Service {
	return &ServiceImpl{
		jwtAuth: jwtauth.New("HS256", []byte(cfg.RefreshTokenSecret), nil),
	}
}

type ServiceImpl struct {
	jwtAuth *jwtauth.JWTAuth
}

func (srv *ServiceImpl) GetJwtAuth() *jwtauth.JWTAuth {
	return srv.jwtAuth
}
