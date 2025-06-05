package accesstoken

import (
	"github.com/go-chi/jwtauth/v5"
)

func NewService(
	cfg *Cfg,
) *Service {
	return &Service{
		jwtAuth: jwtauth.New("HS256", []byte(cfg.AccessTokenSecret), nil),
	}
}

type Service struct {
	jwtAuth *jwtauth.JWTAuth
}

func (srv *Service) GetJwtAuth() *jwtauth.JWTAuth {
	return srv.jwtAuth
}
