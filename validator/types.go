package validator

import (
	"fmt"
	"reflect"
	"time"

	"github.com/ronannnn/infra/model"
	"github.com/ronannnn/infra/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

func GetNonZeroFields(data interface{}) (fields []string) {
	nonZeroFields := make(map[string]interface{})
	getNonZeroFieldsRecursively(data, "", nonZeroFields)
	for field := range nonZeroFields {
		fields = append(fields, field)
	}
	return
}

func getNonZeroFieldsRecursively(data interface{}, prefix string, nonZeroFields map[string]interface{}) {
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	// 遍历结构体的字段
	for i := 0; i < val.NumField(); i++ {
		rawField := val.Field(i)
		fieldName := typ.Field(i).Name

		var field reflect.Value

		// 如果字段是指针类型，需要进一步解引用
		if rawField.Kind() == reflect.Ptr {
			if rawField.IsNil() {
				continue
			}
			field = rawField.Elem()
		} else {
			field = rawField
		}

		subPrefix := utils.JoinNonEmptyStrings(".", prefix, fieldName)
		switch field.Kind() {
		case reflect.Slice, reflect.Array:
			for j := 0; j < field.Len(); j++ {
				getNonZeroFieldsRecursively(field.Index(j).Interface(), utils.JoinNonEmptyStrings(".", prefix, fmt.Sprintf("%s[%d]", fieldName, j)), nonZeroFields)
			}
		case reflect.Struct:
			switch field.Type() {
			case reflect.TypeOf(time.Time{}):
				if rawField.Kind() == reflect.Ptr || !field.Interface().(time.Time).IsZero() {
					nonZeroFields[subPrefix] = nil
				}
			case reflect.TypeOf(model.DecimalSafe{}):
				if rawField.Kind() == reflect.Ptr || !field.Interface().(model.DecimalSafe).Decimal.IsZero() {
					nonZeroFields[subPrefix] = nil
				}
			case reflect.TypeOf(model.BaseModel{}), reflect.TypeOf(gorm.DeletedAt{}), reflect.TypeOf(optimisticlock.Version{}):
				// do nothing
			default:
				getNonZeroFieldsRecursively(field.Interface(), subPrefix, nonZeroFields)
			}
		default:
			if rawField.Kind() == reflect.Ptr || !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
				nonZeroFields[subPrefix] = nil
			}
		}
	}
}
