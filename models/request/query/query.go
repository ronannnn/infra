package query

import (
	"fmt"
	"strings"

	"github.com/ronannnn/infra/utils"
	"gorm.io/gorm/schema"
)

type SelectQueryItem struct {
	Field    string `json:"field"`    // 字段名
	Distinct bool   `json:"distinct"` // 是否去重
}

type WhereQueryItem struct {
	Field string `json:"field"` // 字段名
	Opr   string `json:"opr"`   // 操作类型
	Value any    `json:"value"` // 值
}

type OrderQueryItem struct {
	Field string `json:"field"` // 字段名
	Order string `json:"order"` // 排序方式 asc desc
}

type Query struct {
	Pagination  Pagination        `json:"pagination"`
	SelectQuery []SelectQueryItem `json:"selectQuery"`
	WhereQuery  []WhereQueryItem  `json:"whereQuery"`
	OrderQuery  []OrderQueryItem  `json:"orderQuery"`
}

type Range struct {
	Start any `json:"start"`
	End   any `json:"end"`
}

const (
	TypeCustom    = "custom" // 用户自定义gorm scope，不做任何处理
	TypeEq        = "eq"
	TypeNe        = "ne"
	TypeGt        = "gt"
	TypeGte       = "gte"
	TypeLt        = "lt"
	TypeLte       = "lte"
	TypeLike      = "like"
	TypeStartLike = "start_like"
	TypeEndLike   = "end_like"
	TypeIn        = "in"
	TypeNotIn     = "not_in"
	TypeRange     = "range"
)

type QuerySetter interface {
	schema.Tabler
	FieldColMapper() map[string]string
}

func ResolveQuery(query Query, setter QuerySetter, condition DbCondition) (err error) {
	tblName := setter.TableName()
	fieldColMapper := setter.FieldColMapper()
	if err = ResolveSelectQuery(query.SelectQuery, tblName, fieldColMapper, condition); err != nil {
		return
	}
	if err = ResolveWhereQuery(query.WhereQuery, tblName, fieldColMapper, condition); err != nil {
		return
	}
	if err = ResolveOrderQuery(query.OrderQuery, tblName, fieldColMapper, condition); err != nil {
		return
	}
	return
}

func ResolveSelectQuery(items []SelectQueryItem, tblName string, fieldColMapper map[string]string, condition DbCondition) (err error) {
	for _, item := range items {
		if col, ok := fieldColMapper[item.Field]; ok {
			if item.Distinct {
				// condition.SetDistinct(fmt.Sprintf("`%s`.`%s`", tblName, col)) 似乎带上table的写法会失效
				condition.SetDistinct(col)
			} else {
				condition.SetSelect(fmt.Sprintf("`%s`.`%s`", tblName, col))
			}
		} else {
			return fmt.Errorf("field %s not found", item.Field)
		}
	}
	return
}

func ResolveWhereQuery(items []WhereQueryItem, tblName string, fieldColMapper map[string]string, condition DbCondition) (err error) {
	for _, item := range items {
		if utils.IsZeroValue(item.Value) || item.Opr == TypeCustom {
			continue
		}
		if col, ok := fieldColMapper[item.Field]; ok {
			fullColName := fmt.Sprintf("`%s`.`%s`", tblName, col)
			switch item.Opr {
			case TypeEq:
				condition.SetWhere(fmt.Sprintf("%s = ?", fullColName), []any{item.Value})
			case TypeNe:
				condition.SetNot(fmt.Sprintf("%s != ?", fullColName), []any{item.Value})
			case TypeGt:
				condition.SetWhere(fmt.Sprintf("%s > ?", fullColName), []any{item.Value})
			case TypeGte:
				condition.SetWhere(fmt.Sprintf("%s >= ?", fullColName), []any{item.Value})
			case TypeLt:
				condition.SetWhere(fmt.Sprintf("%s < ?", fullColName), []any{item.Value})
			case TypeLte:
				condition.SetWhere(fmt.Sprintf("%s <= ?", fullColName), []any{item.Value})
			case TypeLike:
				condition.SetWhere(fmt.Sprintf("%s like ?", fullColName), []any{"%" + item.Value.(string) + "%"})
			case TypeStartLike:
				condition.SetWhere(fmt.Sprintf("%s like ?", fullColName), []any{item.Value.(string) + "%"})
			case TypeEndLike:
				condition.SetWhere(fmt.Sprintf("%s like ?", fullColName), []any{"%" + item.Value.(string)})
			case TypeIn:
				condition.SetWhere(fmt.Sprintf("%s in (?)", fullColName), []any{item.Value})
			case TypeNotIn:
				condition.SetWhere(fmt.Sprintf("%s not in (?)", fullColName), []any{item.Value})
			case TypeRange:
				start := item.Value.(map[string]any)["start"]
				end := item.Value.(map[string]any)["end"]
				if !utils.IsZeroValue(start) {
					condition.SetWhere(fmt.Sprintf("`%s`.`%s` >= ?", tblName, col), []any{start})
				}
				if !utils.IsZeroValue(end) {
					condition.SetWhere(fmt.Sprintf("`%s`.`%s` <= ?", tblName, col), []any{end})
				}
			default:
				return fmt.Errorf("opr %s not found", item.Opr)
			}
		} else {
			return fmt.Errorf("field %s not found", item.Field)
		}
	}
	return
}

func ResolveOrderQuery(items []OrderQueryItem, tblName string, fieldColMapper map[string]string, condition DbCondition) (err error) {
	for _, item := range items {
		if col, ok := fieldColMapper[item.Field]; ok {
			switch strings.ToLower(item.Order) {
			case "desc", "asc":
				condition.SetOrder(fmt.Sprintf("`%s`.`%s` %s", tblName, col, item.Order))
			default:
				return fmt.Errorf("order %s not found", item.Order)
			}
		} else {
			return fmt.Errorf("field %s not found", item.Field)
		}
	}
	return
}
