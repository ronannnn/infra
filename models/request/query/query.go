package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ronannnn/infra/utils"
	"gorm.io/gorm"
)

type Range struct {
	Start any `json:"start"`
	End   any `json:"end"`
}

const TagStr = "query"

const (
	TypeSkip   = "-"
	TypeEq     = "eq"
	TypeNe     = "ne"
	TypeGt     = "gt"
	TypeGte    = "gte"
	TypeLt     = "lt"
	TypeLte    = "lte"
	TypeLike   = "like"
	TypeIn     = "in"
	TypeOrder  = "order"
	TypeSelect = "select"
	TypeRange  = "range"
)

// queryTag
// example struct
//
//		type Example struct {
//			Pagination Pagination `json:"pagination" query:"-"` // skip
//			Username   *string    `json:"username" query:"type:like;column:username"`
//	 }
type tag struct {
	Type   string
	Column string
	Table  string
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
		case "type":
			if len(ts) > 1 {
				r.Type = ts[1]
			}
		case "column":
			if len(ts) > 1 {
				r.Column = ts[1]
			}
		case "table":
			if len(ts) > 1 {
				r.Table = ts[1]
			}
		}
	}
	return r
}

func ResolveQuery(search any, condition DbCondition) {
	searchType := reflect.TypeOf(search)
	searchValue := reflect.ValueOf(search)
	for i := 0; i < searchType.NumField(); i++ {
		tagStr, ok := searchType.Field(i).Tag.Lookup(TagStr)
		if !ok {
			// 递归调用
			ResolveQuery(searchValue.Field(i).Interface(), condition)
			continue
		}
		if tagStr == TypeSkip ||
			searchValue.Field(i).IsZero() ||
			((searchValue.Field(i).Kind() == reflect.Array || searchValue.Field(i).Kind() == reflect.Slice) && searchValue.Field(i).Len() == 0) {
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
			condition.SetWhere(fmt.Sprintf("%s = ?", col), []any{searchValue.Field(i).Interface()})
		case TypeNe:
			condition.SetNot(fmt.Sprintf("%s = ?", col), []any{searchValue.Field(i).Interface()})
		case TypeLike:
			condition.SetWhere(fmt.Sprintf("%s like ?", col), []any{"%" + searchValue.Field(i).String() + "%"})
		case TypeGt:
			condition.SetWhere(fmt.Sprintf("%s > ?", col), []any{searchValue.Field(i).Interface()})
		case TypeGte:
			condition.SetWhere(fmt.Sprintf("%s >= ?", col), []any{searchValue.Field(i).Interface()})
		case TypeLt:
			condition.SetWhere(fmt.Sprintf("%s < ?", col), []any{searchValue.Field(i).Interface()})
		case TypeLte:
			condition.SetWhere(fmt.Sprintf("%s <= ?", col), []any{searchValue.Field(i).Interface()})
		case TypeIn:
			condition.SetWhere(fmt.Sprintf("%s in (?)", col), []any{searchValue.Field(i).Interface()})
		case TypeOrder:
			switch strings.ToLower(searchValue.Field(i).String()) {
			case "desc", "asc":
				condition.SetOrder(fmt.Sprintf("%s %s", col, searchValue.Field(i).String()))
			}
		case TypeSelect:
			condition.SetSelect(col)
		case TypeRange:
			start := searchValue.Field(i).Interface().(Range).Start
			end := searchValue.Field(i).Interface().(Range).End
			if !utils.IsZeroValue(start) {
				condition.SetWhere(fmt.Sprintf("%s >= ?", col), []any{start})
			}
			if !utils.IsZeroValue(end) {
				condition.SetWhere(fmt.Sprintf("%s <= ?", col), []any{end})
			}
		}
	}
}

func ResolveQueryRange(queryRange Range, fieldName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		start := queryRange.Start
		end := queryRange.End
		if !utils.IsZeroValue(start) {
			db.Where(fmt.Sprintf("%s >= ?", fieldName), start)
		}
		if !utils.IsZeroValue(end) {
			db.Where(fmt.Sprintf("%s <= ?", fieldName), end)
		}
		return db
	}
}
