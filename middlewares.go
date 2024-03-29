package infra

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/response"
	"github.com/ronannnn/infra/services/apirecord"
	"go.uber.org/zap"
)

type Middleware interface {
	// set request language
	Lang(http.Handler) http.Handler
	// verify privilege
	CasbinInterceptor(http.Handler) http.Handler
	// record request info
	ReqRecorder(http.Handler) http.Handler
}

func ProvideMiddleware(
	log *zap.SugaredLogger,
	casbinEnforcer *casbin.SyncedCachedEnforcer,
	apirecordService apirecord.Service,
) Middleware {
	return &MiddlewareImpl{
		log:              log,
		casbinEnforcer:   casbinEnforcer,
		apirecordService: apirecordService,
	}
}

type MiddlewareImpl struct {
	log              *zap.SugaredLogger
	casbinEnforcer   *casbin.SyncedCachedEnforcer
	apirecordService apirecord.Service
}

func (m *MiddlewareImpl) Lang(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get("Accept-Language")
		if lang == "" {
			lang = string(DefaultLangType)
		}
		ctx := context.WithValue(r.Context(), models.CtxKeyLang, lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareImpl) CasbinInterceptor(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(models.CtxKeyUserId).(uint) // sub
		path := r.URL.Path                                      // obj
		act := r.Method                                         // act
		_, err := m.casbinEnforcer.Enforce(userId, path, act)
		if err != nil {
			response.ErrPrivilege(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *MiddlewareImpl) ReqRecorder(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// get body
		contentType := r.Header.Get("Content-Type")
		var body string
		// assume that requests with contentType 'multipart/form-data' contains files
		// skip record files
		if !strings.Contains(contentType, "multipart/form-data") && r.Body != nil {
			rawData := r.Body
			defer rawData.Close()
			readBytes, _ := io.ReadAll(rawData)
			body = string(readBytes)
			r.Body = io.NopCloser(bytes.NewBuffer(readBytes))
		}
		// wrap response writer to get status code
		// Reference: chi.middleware.Logger
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		defer func() {
			path := r.URL.Path
			raw := r.URL.RawQuery
			if raw != "" {
				path = path + "?" + raw
			}
			userId := uint(0)
			ctxUserId := r.Context().Value(models.CtxKeyUserId)
			if ctxUserId != nil {
				userId = ctxUserId.(uint)
			}
			latency := time.Since(start)
			apiRecord := &models.ApiRecord{
				Base: models.Base{
					OprBy: models.OprBy{
						CreatedBy: userId,
						UpdatedBy: userId,
					},
				},
				Path:       path,
				Method:     r.Method,
				Ip:         r.RemoteAddr,
				Latency:    latency,
				StatusCode: ww.Status(),
				Body:       body,
			}
			m.log.Info(apiRecord)
			// if err := m.apirecordService.Create(apiRecord); err != nil {
			// 	m.log.Warnf("record save error, %v", err)
			// }
		}()
		next.ServeHTTP(ww, r)
	})
}
