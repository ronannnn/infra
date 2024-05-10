package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ronannnn/infra/utils"
)

const TagStr = "query"

type Range struct {
	Start any `json:"start"`
	End   any `json:"end"`
}

const (
	CategorySelect = "select"
	CategoryWhere  = "where"
	CategoryOrder  = "order"
)

const (
	TypeSkip  = "-"
	TypeEq    = "eq"
	TypeNe    = "ne"
	TypeGt    = "gt"
	TypeGte   = "gte"
	TypeLt    = "lt"
	TypeLte   = "lte"
	TypeLike  = "like"
	TypeIn    = "in"
	TypeRange = "range"
)

// queryTag
// example struct
//
//		type Example struct {
//			Pagination Pagination `json:"pagination"` // skip
//			Username   *string    `json:"username" query:"type:like;column:username"`
//	 }
type tag struct {
	Category string
	Type     string
	Table    string
	Column   string
}

func parseQueryTag(queryTag string) *tag {
	r := &tag{}
	queryTags := strings.Split(queryTag, ";")
	var ts []string
	for _, t := range queryTags {
		ts = strings.Split(t, ":")
		if len(ts) == 0 {
			continue
		}
		switch ts[0] {
		case "category":
			if len(ts) > 1 {
				r.Category = ts[1]
			}
		case "type":
			if len(ts) > 1 {
				r.Type = ts[1]
			}
		case "table":
			if len(ts) > 1 {
				r.Table = ts[1]
			}
		case "column":
			if len(ts) > 1 {
				r.Column = ts[1]
			}
		}
	}
	return r
}

func ResolveQuery(search any, condition DbCondition) {
	queryType := reflect.TypeOf(search)
	queryValue := reflect.ValueOf(search)
	for i := 0; i < queryType.NumField(); i++ {
		tagStr, ok := queryType.Field(i).Tag.Lookup(TagStr)
		// 如果没有tag则跳过
		if !ok {
			continue
		}
		// 如果没有tag或者是空值或者是空数组跳过
		// 没有tag
		if !ok ||
			// 是空值
			queryValue.Field(i).IsZero() ||
			// 空数组
			((queryValue.Field(i).Kind() == reflect.Array || queryValue.Field(i).Kind() == reflect.Slice) && queryValue.Field(i).Len() == 0) {
			continue
		}
		tag := parseQueryTag(tagStr)
		if tag.Category == CategorySelect {
			ResolveSelects(queryValue.Field(i).Interface(), condition)
		} else if tag.Category == CategoryWhere {
			ResolveWheres(queryValue.Field(i).Interface(), condition)
		} else if tag.Category == CategoryOrder {
			ResolveOrders(queryValue.Field(i).Interface(), condition)
		}
	}
}

func ResolveSelects(selects any, condition DbCondition) {
	kind := reflect.TypeOf(selects).Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < reflect.ValueOf(selects).Len(); i++ {
			ResolveSelect(reflect.ValueOf(selects).Index(i).Interface(), condition)
		}
	}
}

func ResolveSelect(selectModel any, condition DbCondition) {
	queryType := reflect.TypeOf(selectModel)
	queryValue := reflect.ValueOf(selectModel)
	for i := 0; i < queryType.NumField(); i++ {
		tagStr, ok := queryType.Field(i).Tag.Lookup(TagStr)
		if !ok {
			// 递归调用
			ResolveSelect(queryValue.Field(i).Interface(), condition)
			continue
		}
		// 如果是skip或者是空值则跳过
		if tagStr == TypeSkip || queryValue.Field(i).IsZero() {
			continue
		}
		tag := parseQueryTag(tagStr)
		var col string
		if tag.Table == "" {
			col = fmt.Sprintf("`%s`", tag.Column)
		} else {
			col = fmt.Sprintf("`%s`.`%s`", tag.Table, tag.Column)
		}
		condition.SetSelect(col)
	}
}

func ResolveWheres(wheres any, condition DbCondition) {
	kind := reflect.TypeOf(wheres).Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < reflect.ValueOf(wheres).Len(); i++ {
			ResolveWhere(reflect.ValueOf(wheres).Index(i).Interface(), condition)
		}
	}
}

func ResolveWhere(where any, condition DbCondition) {
	queryType := reflect.TypeOf(where)
	queryValue := reflect.ValueOf(where)
	for i := 0; i < queryType.NumField(); i++ {
		tagStr, ok := queryType.Field(i).Tag.Lookup(TagStr)
		if !ok {
			// 递归调用
			ResolveWhere(queryValue.Field(i).Interface(), condition)
			continue
		}
		// 如果是skip或者是空值或者是空数组跳过
		if tagStr == TypeSkip ||
			queryValue.Field(i).IsZero() ||
			((queryValue.Field(i).Kind() == reflect.Array || queryValue.Field(i).Kind() == reflect.Slice) && queryValue.Field(i).Len() == 0) {
			continue
		}
		tag := parseQueryTag(tagStr)
		var col string
		if tag.Table == "" {
			col = fmt.Sprintf("`%s`", tag.Column)
		} else {
			col = fmt.Sprintf("`%s`.`%s`", tag.Table, tag.Column)
		}
		//解析
		switch tag.Type {
		case TypeEq:
			condition.SetWhere(fmt.Sprintf("%s = ?", col), []any{queryValue.Field(i).Interface()})
		case TypeNe:
			condition.SetNot(fmt.Sprintf("%s = ?", col), []any{queryValue.Field(i).Interface()})
		case TypeLike:
			condition.SetWhere(fmt.Sprintf("%s like ?", col), []any{"%" + queryValue.Field(i).String() + "%"})
		case TypeGt:
			condition.SetWhere(fmt.Sprintf("%s > ?", col), []any{queryValue.Field(i).Interface()})
		case TypeGte:
			condition.SetWhere(fmt.Sprintf("%s >= ?", col), []any{queryValue.Field(i).Interface()})
		case TypeLt:
			condition.SetWhere(fmt.Sprintf("%s < ?", col), []any{queryValue.Field(i).Interface()})
		case TypeLte:
			condition.SetWhere(fmt.Sprintf("%s <= ?", col), []any{queryValue.Field(i).Interface()})
		case TypeIn:
			condition.SetWhere(fmt.Sprintf("%s in (?)", col), []any{queryValue.Field(i).Interface()})
		case TypeRange:
			start := queryValue.Field(i).Interface().(Range).Start
			end := queryValue.Field(i).Interface().(Range).End
			if !utils.IsZeroValue(start) {
				condition.SetWhere(fmt.Sprintf("%s >= ?", col), []any{start})
			}
			if !utils.IsZeroValue(end) {
				condition.SetWhere(fmt.Sprintf("%s <= ?", col), []any{end})
			}
		}
	}
}

func ResolveOrders(orders any, condition DbCondition) {
	kind := reflect.TypeOf(orders).Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < reflect.ValueOf(orders).Len(); i++ {
			ResolveOrder(reflect.ValueOf(orders).Index(i).Interface(), condition)
		}
	}
}

func ResolveOrder(order any, condition DbCondition) {
	queryType := reflect.TypeOf(order)
	queryValue := reflect.ValueOf(order)
	for i := 0; i < queryType.NumField(); i++ {
		tagStr, ok := queryType.Field(i).Tag.Lookup(TagStr)
		if !ok {
			// 递归调用
			ResolveOrder(queryValue.Field(i).Interface(), condition)
			continue
		}
		if tagStr == TypeSkip ||
			queryValue.Field(i).IsZero() ||
			((queryValue.Field(i).Kind() == reflect.Array || queryValue.Field(i).Kind() == reflect.Slice) && queryValue.Field(i).Len() == 0) {
			continue
		}
		tag := parseQueryTag(tagStr)
		var col string
		if tag.Table == "" {
			col = fmt.Sprintf("`%s`", tag.Column)
		} else {
			col = fmt.Sprintf("`%s`.`%s`", tag.Table, tag.Column)
		}
		//解析
		switch strings.ToLower(queryValue.Field(i).String()) {
		case "desc", "asc":
			condition.SetOrder(fmt.Sprintf("%s %s", col, queryValue.Field(i).String()))
		}
	}
}
