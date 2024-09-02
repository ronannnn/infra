package models

import (
	"fmt"
	"reflect"

	"github.com/ronannnn/infra/utils"
	"github.com/shopspring/decimal"
)

type DecimalSafe struct {
	decimal.Decimal
}

func (d *DecimalSafe) UnmarshalJSON(decimalBytes []byte) error {
	if string(decimalBytes) == `""` {
		return nil
	}

	return d.Decimal.UnmarshalJSON(decimalBytes)
}

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
		case reflect.Struct:
			getNonZeroFieldsRecursively(field.Interface(), subPrefix, nonZeroFields)
		case reflect.Slice, reflect.Array:
			for j := 0; j < field.Len(); j++ {
				rawElem := field.Index(j)
				var elem reflect.Value
				// 如果字段是指针类型，需要进一步解引用
				if rawElem.Kind() == reflect.Ptr {
					if rawElem.IsNil() {
						continue
					}
					elem = rawElem.Elem()
				} else {
					elem = rawElem
				}
				sliceSubPrefix := utils.JoinNonEmptyStrings(".", prefix, fmt.Sprintf("%s[%d]", fieldName, j))
				if elem.Kind() == reflect.Struct {
					getNonZeroFieldsRecursively(elem.Interface(), sliceSubPrefix, nonZeroFields)
				} else {
					if rawElem.Kind() == reflect.Ptr || !reflect.DeepEqual(rawElem.Interface(), reflect.Zero(rawElem.Type()).Interface()) {
						nonZeroFields[sliceSubPrefix] = nil
					}
				}
			}
		default:
			if rawField.Kind() == reflect.Ptr || !reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
				nonZeroFields[subPrefix] = nil
			}
		}
	}
}
