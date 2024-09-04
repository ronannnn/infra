package validator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/microcosm-cc/bluemonday"
	"github.com/ronannnn/infra/constant"
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
	validate.RegisterTranslation("d_gt", trans, func(ut ut.Translator) error {
		switch lang {
		case i18n.LanguageChinese:
			return ut.Add("d_gt", "{0}必须大于{1}", true)
		case i18n.LanguageEnglish:
			return ut.Add("d_gt", "{0} must be greater than {1}", true)
		default:
			return ut.Add("d_gt", "{0} must be greater than {1}", true)
		}
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("d_gt", fe.Field(), fe.Param())
		return t
	})
	_ = validate.RegisterValidation("d_lt", decimalLessThan)
	validate.RegisterTranslation("d_lt", trans, func(ut ut.Translator) error {
		switch lang {
		case i18n.LanguageChinese:
			return ut.Add("d_lt", "{0}必须小于{1}", true)
		case i18n.LanguageEnglish:
			return ut.Add("d_lt", "{0} must be less than {1}", true)
		default:
			return ut.Add("d_lt", "{0} must be less than {1}", true)
		}
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("d_lt", fe.Field(), fe.Param())
		return t
	})
	_ = validate.RegisterValidation("d_decimal_len_lte", decimalDecimalPartsLenLessThanOrEqual)
	validate.RegisterTranslation("d_decimal_len_lte", trans, func(ut ut.Translator) error {
		switch lang {
		case i18n.LanguageChinese:
			return ut.Add("d_decimal_len_lte", "{0}小数点位数必须小于或等于{1}", true)
		case i18n.LanguageEnglish:
			return ut.Add("d_decimal_len_lte", "{0} decimals length must be less than or equal to {1}", true)
		default:
			return ut.Add("d_decimal_len_lte", "{0} decimals length must be less than or equal to {1}", true)
		}
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("d_decimal_len_lte", fe.Field(), fe.Param())
		return t
	})
	_ = validate.RegisterValidation("cn_car", chineseCarNo)
	validate.RegisterTranslation("cn_car", trans, func(ut ut.Translator) error {
		switch lang {
		case i18n.LanguageChinese:
			return ut.Add("cn_car", "{0}不是有效的中国车牌号", true)
		case i18n.LanguageEnglish:
			return ut.Add("cn_car", "{0} is not a qualified chinese car number", true)
		default:
			return ut.Add("cn_car", "{0} is not a qualified chinese car number", true)
		}
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("cn_car", fe.Value().(string), fe.Param())
		return t
	})
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

func chineseCarNo(fl validator.FieldLevel) bool {
	data, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return regexp.MustCompile(constant.CnCarRegexp).MatchString(data)
}
