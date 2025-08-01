package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ronannnn/infra/constant"
	"github.com/ronannnn/infra/i18n"
)

// GetLang get language from header
func GetLang(r *http.Request) i18n.Language {
	return GetLangByCtx(r.Context())
}

// GetLangByCtx get language from header
func GetLangByCtx(ctx context.Context) i18n.Language {
	acceptLanguage, ok := ctx.Value(constant.CtxKeyAcceptLanguage).(i18n.Language)
	if ok {
		return acceptLanguage
	}
	return i18n.DefaultLanguage
}

// GetShowType get show type from header
func GetShowType(r *http.Request, defaultShowType ShowType) ShowType {
	return GetShowTypeByCtx(r.Context(), defaultShowType)
}

// GetShowTypeByCtx get show type from header
func GetShowTypeByCtx(ctx context.Context, defaultShowType ShowType) ShowType {
	intShowType, err := strconv.Atoi(ctx.Value(constant.CtxKeyShowType).(string))
	if err != nil {
		return defaultShowType
	}

	switch ShowType(intShowType) {
	case ShowTypeSilent,
		ShowTypeSuccessMessage,
		ShowTypeInfoMessage,
		ShowTypeWarningMessage,
		ShowTypeErrorMessage,
		ShowTypeSuccessNotification,
		ShowTypeInfoNotification,
		ShowTypeWarningNotification,
		ShowTypeErrorNotification:
		return ShowType(intShowType)
	default:
		return defaultShowType
	}
}
