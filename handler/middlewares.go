package handler

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/ronannnn/infra/constant"
	"github.com/ronannnn/infra/i18n"
	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/service"
	"go.uber.org/zap"
)

type Middleware interface {
	// header info
	Lang(http.Handler) http.Handler
	Ua(http.Handler) http.Handler
	DeviceId(http.Handler) http.Handler
	ReqRecorder(http.Handler) http.Handler
}

func ProvideMiddleware(
	log *zap.SugaredLogger,
	apiRecordService *service.ApiRecordService,
) Middleware {
	return &MiddlewareImpl{
		log:              log,
		apiRecordService: apiRecordService,
	}
}

type MiddlewareImpl struct {
	log              *zap.SugaredLogger
	apiRecordService *service.ApiRecordService
}

func (m *MiddlewareImpl) Lang(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get(string(constant.CtxKeyAcceptLanguage))
		if lang == "" {
			lang = string(i18n.DefaultLanguage)
		}
		ctx := context.WithValue(r.Context(), constant.CtxKeyAcceptLanguage, lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareImpl) Ua(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get(string(constant.CtxKeyUa))
		ctx := context.WithValue(r.Context(), constant.CtxKeyUa, ua)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *MiddlewareImpl) DeviceId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deviceId := r.Header.Get(string(constant.CtxKeyDeviceId))
		ctx := context.WithValue(r.Context(), constant.CtxKeyDeviceId, deviceId)
		next.ServeHTTP(w, r.WithContext(ctx))
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
			ctxUserId := r.Context().Value(constant.CtxKeyUserId)
			if ctxUserId != nil {
				userId = ctxUserId.(uint)
			}
			latency := time.Since(start)
			apiRecord := &model.ApiRecord{
				Base: model.Base{
					OprBy: model.OprBy{
						CreatedBy: &userId,
						UpdatedBy: &userId,
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
