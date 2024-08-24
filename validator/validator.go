package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/locales"
	english "github.com/go-playground/locales/en"
	chinese "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/go-playground/validator/v10/translations/zh"
	"github.com/ronannnn/infra/i18n"
	"github.com/ronannnn/infra/msg"
	"github.com/ronannnn/infra/reason"
	"github.com/ronannnn/infra/utils"
)

type TranslatorLocal struct {
	La           i18n.Language
	Lo           locales.Translator
	RegisterFunc func(v *validator.Validate, trans ut.Translator) (err error)
}

var (
	allLanguageTranslators = []*TranslatorLocal{
		{La: i18n.LanguageChinese, Lo: chinese.New(), RegisterFunc: zh.RegisterDefaultTranslations},
		{La: i18n.LanguageEnglish, Lo: english.New(), RegisterFunc: en.RegisterDefaultTranslations},
	}
)

// ValidatorEle validator
type ValidatorEle struct {
	Validate *validator.Validate
	Tran     ut.Translator
}

// FormErrorField indicates the current form error content. which field is error and error message.
type FormErrorField struct {
	ErrorField         string `json:"errorField"`
	ErrorWithNamespace string `json:"errorFieldWithNamespace"`
	ErrorMsg           string `json:"errorMsg"`
}

type Validator interface {
	Check(lang i18n.Language, value interface{}) (errFields []*FormErrorField, err error)
}

func New(wiredI18n i18n.I18n) Validator {
	impl := &Impl{
		i18n:           wiredI18n,
		validatorItems: make(map[i18n.Language]*ValidatorEle, 0),
	}
	for _, t := range allLanguageTranslators {
		tran, ok := ut.New(t.Lo, t.Lo).GetTranslator(t.Lo.Locale())
		if !ok {
			panic(fmt.Sprintf("not found translator %s", t.Lo.Locale()))
		}
		val := createValidateWithCustomValidations(t.La, wiredI18n)
		if t.RegisterFunc != nil {
			if err := t.RegisterFunc(val, tran); err != nil {
				panic(err)
			}
		}
		impl.validatorItems[t.La] = &ValidatorEle{Validate: val, Tran: tran}
	}
	return impl
}

type Impl struct {
	i18n           i18n.I18n
	validatorItems map[i18n.Language]*ValidatorEle
}

func (m *Impl) Check(lang i18n.Language, value interface{}) (errFields []*FormErrorField, err error) {
	v, ok := m.validatorItems[lang]
	if !ok {
		err = msg.NewError(reason.ValidatorLangNotFound)
		return
	}
	err = v.Validate.Struct(value)
	if err != nil {
		var valErrors validator.ValidationErrors
		if !errors.As(err, &valErrors) {
			err = fmt.Errorf("validate check exception, %v", err)
			return
		}

		for _, fieldError := range valErrors {
			errFields = append(errFields, &FormErrorField{
				ErrorField:         utils.LowercaseFirstLetter(fieldError.StructField()),
				ErrorWithNamespace: utils.LowercaseFirstLetterAndJoin(fieldError.StructNamespace(), "."),
				ErrorMsg:           fieldError.Translate(v.Tran),
			})
		}
		err = msg.NewError("fields validation failed").WithMsg("fields validation failed")
	}
	return
}
