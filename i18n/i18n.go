package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	goI18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type Language string

const (
	LanguageChinese Language = "zh_CN"
	LanguageEnglish Language = "en_US"

	DefaultLanguage = LanguageChinese
)

type I18n interface {
	Tr(lang Language, key string) string
	TrWithData(lang Language, key string, templateData any) string
}

// New new i18n from Bundle/resource directory
func New(cfg *Cfg) (i18n I18n, err error) {
	impl := &Impl{
		localizes: make(map[Language]*goI18n.Localizer),
	}
	bundle := goI18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	// 验证bundleDir是否是一个目录
	var stat os.FileInfo
	if stat, err = os.Stat(cfg.BundleDir); err != nil {
		return
	}
	if !stat.IsDir() {
		err = fmt.Errorf("%s is not a directory", cfg.BundleDir)
		return
	}

	var entries []os.DirEntry
	if entries, err = os.ReadDir(cfg.BundleDir); err != nil {
		return
	}

	// read the Bundle resources file from entries
	for _, file := range entries {
		// ignore directory
		if file.IsDir() {
			continue
		}
		// ignore non-YAML file
		if filepath.Ext(file.Name()) != ".yaml" {
			continue
		}

		// read the resource from single entry file
		var fileBytes []byte
		if fileBytes, err = os.ReadFile(filepath.Join(cfg.BundleDir, file.Name())); err != nil {
			return nil, err
		}

		// the default localizes format is yaml
		if _, err = bundle.ParseMessageFileBytes(fileBytes, file.Name()); err != nil {
			err = fmt.Errorf("parse language message file [%s] failed: %s", file.Name(), err)
			return
		}
		languageName := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		impl.localizes[Language(languageName)] = goI18n.NewLocalizer(bundle, languageName)
	}
	return impl, nil
}

// Impl save the language type with localizer mapping
type Impl struct {
	localizes map[Language]*goI18n.Localizer
}

// Tr to Translate from specified language and string
func (tr *Impl) Tr(la Language, key string) string {
	return tr.TrWithData(la, key, nil)
}

func (tr *Impl) TrWithData(la Language, key string, templateData any) string {
	l, ok := tr.localizes[la]
	if !ok {
		l = tr.localizes[DefaultLanguage]
	}
	if l == nil {
		return fmt.Sprintf("Localizer for language %s not found, try to translate key %s", la, key)
	}

	translation, err := l.Localize(&goI18n.LocalizeConfig{MessageID: key, TemplateData: templateData})
	if _, ok := err.(*goI18n.MessageNotFoundErr); ok {
		return key // 返回key即表示未找到对应的翻译
	} else if err != nil {
		return fmt.Sprintf("Translation for key %s failed: %s", key, err)
	}

	return translation
}
