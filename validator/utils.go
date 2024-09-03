package validator

import (
	"reflect"
	"strconv"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"github.com/ronannnn/infra/i18n"
	"github.com/ronannnn/infra/models"
	"github.com/shopspring/decimal"
)

func createValidateWithCustomValidations(lang i18n.Language, trans ut.Translator) *validator.Validate {
	validate := validator.New()
	// 注册自定义类型
	validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if valuer, ok := field.Interface().(models.DecimalSafe); ok {
			return valuer.String()
		}
		return nil
	}, models.DecimalSafe{})
	validate.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {
		if valuer, ok := field.Interface().(decimal.Decimal); ok {
			return valuer.String()
		}
		return nil
	}, decimal.Decimal{})
	_ = validate.RegisterValidation("not_blank", notBlank)
	validate.RegisterTranslation("not_blank", trans, func(ut ut.Translator) error {
		switch lang {
		case i18n.LanguageChinese:
			return ut.Add("not_blank", "{0}不能为空", true)
		case i18n.LanguageEnglish:
			return ut.Add("not_blank", "{0} must not be empty", true)
		default:
			return ut.Add("not_blank", "{0} must not be empty", true)
		}
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("not_blank", fe.Field())
		return t
	})
	_ = validate.RegisterValidation("sanitizer", sanitizer)
	_ = validate.RegisterValidation("d_gt", decimalGreaterThan)
	_ = validate.RegisterValidation("d_lt", decimalLessThan)
	_ = validate.RegisterValidation("d_decimal_len_lte", decimalDecimalPartsLenLessThanOrEqual)
	return validate
}

func notBlank(fl validator.FieldLevel) (res bool) {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		trimSpace := strings.TrimSpace(field.String())
		if len(trimSpace) > 0 {
			field.SetString(trimSpace)
			return true
		}
		return false
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

func decimalGreaterThan(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	value, err := decimal.NewFromString(data)
	if err != nil {
		return false
	}
	baseValue, err := decimal.NewFromString(fl.Param())
	if err != nil {
		return false
	}
	return value.GreaterThan(baseValue)
}

func decimalLessThan(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	value, err := decimal.NewFromString(data)
	if err != nil {
		return false
	}
	baseValue, err := decimal.NewFromString(fl.Param())
	if err != nil {
		return false
	}
	return value.LessThan(baseValue)
}

// decimalDecimalPartsLenLessThanOrEqual 验证小数点位数是否小于等于指定值
func decimalDecimalPartsLenLessThanOrEqual(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	value, err := decimal.NewFromString(data)
	if err != nil {
		return false
	}
	// 获取小数点位数
	var actualLen int
	split := strings.Split(value.String(), ".")
	if len(split) == 1 {
		actualLen = 0
	} else {
		actualLen = len(split[1])
	}
	expectedLen, err := strconv.Atoi(fl.Param())
	if err != nil {
		return false
	}
	return actualLen <= expectedLen
}
