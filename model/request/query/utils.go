package query

import (
	"reflect"

	"github.com/ronannnn/infra/utils"
)

func CamelToSnakeFromStruct(obj any) map[string]string {
	fields := []string{}
	getJsonTagsFromStruct(obj, &fields)
	return camelToSnakeWithBaseFromStrings(fields)
}

func camelToSnakeWithBaseFromStrings(fields []string) map[string]string {
	mapper := make(map[string]string)
	for _, field := range fields {
		mapper[field] = utils.CamelToSnake(field)
	}
	return mapper
}

func getJsonTagsFromStruct(obj any, fields *[]string) {
	structValue := reflect.ValueOf(obj)
	if structValue.Kind() == reflect.Ptr {
		if structValue.IsNil() {
			return
		}
		structValue = structValue.Elem()
	}
	structType := structValue.Type()
	for i := range structType.NumField() {
		jsonTag, jsonTagOk := structType.Field(i).Tag.Lookup("json")
		if !jsonTagOk {
			getJsonTagsFromStruct(structValue.Field(i).Interface(), fields)
			continue
		}
		if jsonTag == "-" {
			continue
		}
		gormTag, gormTagOk := structType.Field(i).Tag.Lookup("gorm")
		if gormTagOk && gormTag == "-" {
			continue
		}
		*fields = append(*fields, jsonTag)
	}
}
