package infra

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/ronannnn/infra/models"
	"github.com/ronannnn/infra/models/response"
	"github.com/ronannnn/infra/services/apirecord"
	"github.com/ronannnn/infra/services/auth/accesstoken"
	"go.uber.org/zap"
)

type Middleware interface {
	// set request language
	Lang(http.Handler) http.Handler
	// verify privilege
	CasbinInterceptor(http.Handler) http.Handler
	// auth handlers
	AuthHandlers() []func(http.Handler) http.Handler
	Verifier(http.Handler) http.Handler
	Authenticator(http.Handler) http.Handler
	AuthInfoSetter(next http.Handler) http.Handler
	// record request info
	ReqRecorder(http.Handler) http.Handler
}

func ProvideMiddleware(
	log *zap.SugaredLogger,
	casbinEnforcer *casbin.SyncedCachedEnforcer,
	accesstokenService accesstoken.Service,
	apirecordService apirecord.Service,
) Middleware {
	return &MiddlewareImpl{
		log:                log,
		casbinEnforcer:     casbinEnforcer,
		accesstokenService: accesstokenService,
		apirecordService:   apirecordService,
	}
}

type MiddlewareImpl struct {
	log                *zap.SugaredLogger
	casbinEnforcer     *casbin.SyncedCachedEnforcer
	accesstokenService accesstoken.Service
	apirecordService   apirecord.Service
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

func (m *MiddlewareImpl) AuthHandlers() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		m.Verifier,
		m.Authenticator,
		m.AuthInfoSetter,
	}
}

func (m *MiddlewareImpl) Verifier(next http.Handler) http.Handler {
	return jwtauth.Verifier(m.accesstokenService.GetJwtAuth())(next)
}

// Authenticator override chi.Authenticator
func (m *MiddlewareImpl) Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			response.ErrAccessToken(w, r, err)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			response.ErrAccessToken(w, r, fmt.Errorf("invalid token"))
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

// AuthInfoSetter is a middleware that sets the auth info(user id and username) for the request.
// It must be placed after jwt middleware.
func (m *MiddlewareImpl) AuthInfoSetter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, _ := jwtauth.FromContext(r.Context())
		username, _ := token.Get("username")
		userId, _ := token.Get("userId")
		ctx := context.WithValue(r.Context(), models.CtxKeyUserId, uint(userId.(float64)))
		ctx = context.WithValue(ctx, models.CtxKeyUsername, username.(string))
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
