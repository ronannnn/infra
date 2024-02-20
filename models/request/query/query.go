package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ronannnn/infra"
	"gorm.io/gorm"
)

type QueryRange struct {
	Start any `json:"start"`
	End   any `json:"end"`
}

const (
	SearchTag = "search"
	Skip      = "-"
	Eq        = "eq"
	Ne        = "ne"
	Gt        = "gt"
	Gte       = "gte"
	Lt        = "lt"
	Lte       = "lte"
	Like      = "like"
	In        = "in"
	Order     = "order"
	Select    = "select"
	Range     = "range"
)

// searchTag
// example struct
//
//		type Example struct {
//			Pagination Pagination `json:"pagination" search:"-"` // skip
//			Username   *string    `json:"username" search:"type:like;column:username"`
//	 }
type searchTag struct {
	Type   string
	Column string
	Table  string
}

func parseQueryTag(tag string) *searchTag {
	r := &searchTag{}
	tags := strings.Split(tag, ";")
	var ts []string
	for _, t := range tags {
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
		tagStr, ok := searchType.Field(i).Tag.Lookup(SearchTag)
		if !ok {
			// 递归调用
			ResolveQuery(searchValue.Field(i).Interface(), condition)
			continue
		}
		if tagStr == Skip ||
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
		case Eq:
			condition.SetWhere(fmt.Sprintf("%s = ?", col), []any{searchValue.Field(i).Interface()})
		case Ne:
			condition.SetNot(fmt.Sprintf("%s = ?", col), []any{searchValue.Field(i).Interface()})
		case Like:
			condition.SetWhere(fmt.Sprintf("%s like ?", col), []any{"%" + searchValue.Field(i).String() + "%"})
		case Gt:
			condition.SetWhere(fmt.Sprintf("%s > ?", col), []any{searchValue.Field(i).Interface()})
		case Gte:
			condition.SetWhere(fmt.Sprintf("%s >= ?", col), []any{searchValue.Field(i).Interface()})
		case Lt:
			condition.SetWhere(fmt.Sprintf("%s < ?", col), []any{searchValue.Field(i).Interface()})
		case Lte:
			condition.SetWhere(fmt.Sprintf("%s <= ?", col), []any{searchValue.Field(i).Interface()})
		case In:
			condition.SetWhere(fmt.Sprintf("%s in (?)", col), []any{searchValue.Field(i).Interface()})
		case Order:
			switch strings.ToLower(searchValue.Field(i).String()) {
			case "desc", "asc":
				condition.SetOrder(fmt.Sprintf("%s %s", col, searchValue.Field(i).String()))
			}
		case Select:
			condition.SetSelect(col)
		case Range:
			start := searchValue.Field(i).Interface().(QueryRange).Start
			end := searchValue.Field(i).Interface().(QueryRange).End
			if !infra.IsZeroValue(start) {
				condition.SetWhere(fmt.Sprintf("%s >= ?", col), []any{start})
			}
			if !infra.IsZeroValue(end) {
				condition.SetWhere(fmt.Sprintf("%s <= ?", col), []any{end})
			}
		}
	}
}

func ResolveQueryRange(queryRange QueryRange, fieldName string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		start := queryRange.Start
		end := queryRange.End
		if !infra.IsZeroValue(start) {
			db.Where(fmt.Sprintf("%s >= ?", fieldName), start)
		}
		if !infra.IsZeroValue(end) {
			db.Where(fmt.Sprintf("%s <= ?", fieldName), end)
		}
		return db
	}
}
