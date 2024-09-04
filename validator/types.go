package validator

import (
	"fmt"
	"reflect"

	"github.com/ronannnn/infra/utils"
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

	for i := 0; i < val.NumField(); i++ {
		rawField := val.Field(i)
		// 跳过非validate字段
		tag := typ.Field(i).Tag.Get("validate")
		if tag == "" {
			continue
		}
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
			getNonZeroFieldsRecursively(field.Interface(), subPrefix, nonZeroFields)
		default:
			if rawField.Kind() == reflect.Ptr || !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
				nonZeroFields[subPrefix] = nil
			}
		}
	}
}
