package accesstoken

import (
	"github.com/go-chi/jwtauth/v5"
	"github.com/ronannnn/infra/cfg"
)

type Service interface {
	GetJwtAuth() *jwtauth.JWTAuth
}

func ProvideService(
	cfg *cfg.Auth,
) Service {
	return &ServiceImpl{
		jwtAuth: jwtauth.New("HS256", []byte(cfg.AccessTokenSecret), nil),
	}
}

type ServiceImpl struct {
	jwtAuth *jwtauth.JWTAuth
}

func (srv *ServiceImpl) GetJwtAuth() *jwtauth.JWTAuth {
	return srv.jwtAuth
}
