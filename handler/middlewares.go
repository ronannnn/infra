package handler

import (
	"context"
	"net/http"

	"github.com/ronannnn/infra/constant"
	"github.com/ronannnn/infra/i18n"
	"github.com/ronannnn/infra/models"
)

type Middleware interface {
	// header info
	Lang(http.Handler) http.Handler
	Ua(http.Handler) http.Handler
	DeviceId(http.Handler) http.Handler
}

func ProvideMiddleware() Middleware {
	return &MiddlewareImpl{}
}

type MiddlewareImpl struct {
}

func (m *MiddlewareImpl) Lang(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get(string(constant.CtxKeyAcceptLanguage))
		if lang == "" {
			lang = string(i18n.DefaultLanguage)
		}
		ctx := context.WithValue(r.Context(), models.CtxKeyLang, lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareImpl) Ua(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get(string(constant.CtxKeyUa))
		ctx := context.WithValue(r.Context(), models.CtxKeyUa, ua)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareImpl) DeviceId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deviceId := r.Header.Get(string(constant.CtxKeyDeviceId))
		ctx := context.WithValue(r.Context(), models.CtxKeyDeviceId, deviceId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
