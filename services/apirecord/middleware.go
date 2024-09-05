package apirecord

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/ronannnn/infra/constant"
	"github.com/ronannnn/infra/models"
	"go.uber.org/zap"
)

type Middleware interface {
	ReqRecorder(http.Handler) http.Handler
}

func ProvideMiddleware(
	log *zap.SugaredLogger,
	apiRecordService Service,
) Middleware {
	return &MiddlewareImpl{
		log:              log,
		apirecordService: apiRecordService,
	}
}

type MiddlewareImpl struct {
	log              *zap.SugaredLogger
	apirecordService Service
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
			apiRecord := &ApiRecord{
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
