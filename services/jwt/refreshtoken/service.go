package refreshtoken

import (
	"github.com/go-chi/jwtauth/v5"
)

type Service interface {
	GetJwtAuth() *jwtauth.JWTAuth
}

func ProvideService(
	cfg *Cfg,
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
