package handler

import (
	"context"
	"net/http"

	"github.com/ronannnn/infra/constant"
	"github.com/ronannnn/infra/i18n"
)

// GetLang get language from header
func GetLang(r *http.Request) i18n.Language {
	return GetLangByCtx(r.Context())
}

// GetLangByCtx get language from header
func GetLangByCtx(ctx context.Context) i18n.Language {
	acceptLanguage, ok := ctx.Value(string(constant.CtxKeyAcceptLanguage)).(i18n.Language)
	if ok {
		return acceptLanguage
	}
	return i18n.DefaultLanguage
}
