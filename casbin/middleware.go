package casbin

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/ronannnn/infra/constant"
	"github.com/ronannnn/infra/handler"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
)

type Middleware interface {
	CasbinInterceptor(http.Handler) http.Handler
}

func ProvideMiddleware(
	casbinEnforcer *casbin.SyncedCachedEnforcer,
	// handler
	httpHandler handler.HttpHandler,
) Middleware {
	return &MiddlewareImpl{
		casbinEnforcer: casbinEnforcer,
		httpHandler:    httpHandler,
	}
}

type MiddlewareImpl struct {
	casbinEnforcer *casbin.SyncedCachedEnforcer
	httpHandler    handler.HttpHandler
}

func (m *MiddlewareImpl) CasbinInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(constant.CtxKeyUserId).(uint) // sub
		path := r.URL.Path                                        // obj
		act := r.Method                                           // act
		_, err := m.casbinEnforcer.Enforce(userId, path, act)
		if err != nil {
			m.httpHandler.Fail(w, r, msg.NewError(reason.ForbiddenError), nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
