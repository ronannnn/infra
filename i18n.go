package infra

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/ronannnn/infra/cfg"
	"golang.org/x/text/language"
)

type LangType string

var DefaultLangType = LangTypeZhCn

const (
	LangTypeZhCn LangType = "zh-cn"
	LangTypeEnUs LangType = "en-us"
)

type I18n interface {
	TFromCtx(ctx context.Context, key string, args ...any) string
	T(lang LangType, key string, args ...any) string
}

func NewI18n(cfg cfg.I18n) I18n {
	// init zh-CN bundle
	zhCnBundle := i18n.NewBundle(language.SimplifiedChinese)
	zhCnBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	zhCnBundle.LoadMessageFile(filepath.Join(cfg.LocalesDir, fmt.Sprintf("%s.toml", cfg.ZhCnTomlFilenamePrefix)))
	// init en-US bundle
	enUsBundle := i18n.NewBundle(language.AmericanEnglish)
	enUsBundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	enUsBundle.LoadMessageFile(filepath.Join(cfg.LocalesDir, fmt.Sprintf("%s.toml", cfg.EnUsTomlFilenamePrefix)))
	return &I18nImpl{
		zhCnLocalizer: i18n.NewLocalizer(zhCnBundle, string(LangTypeZhCn)),
		enUsLocalizer: i18n.NewLocalizer(enUsBundle, string(LangTypeEnUs)),
		ctxKey:        cfg.CtxKey,
	}
}

type I18nImpl struct {
	zhCnLocalizer *i18n.Localizer
	enUsLocalizer *i18n.Localizer
	ctxKey        string
}

func (i *I18nImpl) TFromCtx(ctx context.Context, key string, args ...any) string {
	return i.T(LangType(strings.ToLower(ctx.Value(i.ctxKey).(string))), key, args...)
}

func (i *I18nImpl) T(lang LangType, key string, args ...any) string {
	// get accept-language from context
	switch lang {
	case LangTypeEnUs:
		return i.enUsT(key, args...)
	case LangTypeZhCn:
		fallthrough
	default:
		return i.zhCnT(key, args...)
	}
}

func (i *I18nImpl) zhCnT(key string, args ...any) string {
	return i.zhCnLocalizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: args,
	})
}

func (i *I18nImpl) enUsT(key string, args ...any) string {
	return i.enUsLocalizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: args,
	})
}
