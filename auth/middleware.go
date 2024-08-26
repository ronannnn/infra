package auth

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/ronannnn/infra/cfg"
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"github.com/ronannnn/infra/services/jwt/accesstoken"
)

type Middleware interface {
	AuthHandlers() []func(http.Handler) http.Handler
	Verifier(http.Handler) http.Handler
	Authenticator(http.Handler) http.Handler
	AuthInfoSetter(next http.Handler) http.Handler
}

func ProvideMiddleware(
	// auth
	authCfg *cfg.Auth,
	accessTokenService accesstoken.Service,
	// handler
	httpHandler handler.HttpHandler,
) Middleware {
	return &MiddlewareImpl{
		authCfg:            authCfg,
		accessTokenService: accessTokenService,
		httpHandler:        httpHandler,
	}
}

type MiddlewareImpl struct {
	authCfg            *cfg.Auth
	accessTokenService accesstoken.Service
	httpHandler        handler.HttpHandler
}

func (m *MiddlewareImpl) AuthHandlers() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		m.Verifier,
		m.Authenticator,
		m.AuthInfoSetter,
	}
}

func (m *MiddlewareImpl) Verifier(next http.Handler) http.Handler {
	if !m.authCfg.Enabled {
		return next
	}
	return jwtauth.Verifier(m.accessTokenService.GetJwtAuth())(next)
}

// Authenticator override chi.Authenticator
func (m *MiddlewareImpl) Authenticator(next http.Handler) http.Handler {
	if !m.authCfg.Enabled {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			m.httpHandler.FailWithCode(w, r, msg.NewError(reason.UnauthorizedError).WithError(err), nil, handler.AccessTokenErrorCode)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			m.httpHandler.FailWithCode(w, r, msg.NewError(reason.UnauthorizedError).WithError(err), nil, handler.AccessTokenErrorCode)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

// AuthInfoSetter is a middleware that sets the auth info(user id and username) for the request.
// It must be placed after jwt middleware.
func (m *MiddlewareImpl) AuthInfoSetter(next http.Handler) http.Handler {
	if !m.authCfg.Enabled {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, _ := jwtauth.FromContext(r.Context())
		username, _ := token.Get("username")
		userId, _ := token.Get("userId")
		ctx := context.WithValue(r.Context(), models.CtxKeyUserId, uint(userId.(float64)))
		ctx = context.WithValue(ctx, models.CtxKeyUsername, username.(string))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
