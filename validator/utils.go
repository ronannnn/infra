package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"github.com/ronannnn/infra/i18n"
)

func createValidateWithCustomValidations(la i18n.Language, wiredI18n i18n.I18n) *validator.Validate {
	validate := validator.New()
	_ = validate.RegisterValidation("not-blank", notBlank)
	_ = validate.RegisterValidation("sanitizer", sanitizer)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) (res string) {
		// TODO: 提issue，加入namespace信息
		return wiredI18n.Tr(la, fld.Name)
	})
	return validate
}

func notBlank(fl validator.FieldLevel) (res bool) {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		trimSpace := strings.TrimSpace(field.String())
		res := len(trimSpace) > 0
		if !res {
			field.SetString(trimSpace)
		}
		return true
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func sanitizer(fl validator.FieldLevel) (res bool) {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		filter := bluemonday.UGCPolicy()
		content := strings.Replace(filter.Sanitize(field.String()), "&amp;", "&", -1)
		field.SetString(content)
		return true
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}
