package query

import (
	"fmt"
	"strings"

	"github.com/ronannnn/infra/utils"
	"gorm.io/gorm/schema"
)

// selectQuery
// SelectQueryItem 用于指定查询的字段和是否去重
type SelectQueryItem struct {
	Field    string `json:"field"`    // 字段名
	Distinct bool   `json:"distinct"` // 是否去重
}

// whereQuery
type WhereQuery struct {
	Items  []WhereQueryItem      `json:"items"`  // 简单查询条件
	Groups []WhereQueryItemGroup `json:"groups"` // 高级查询条件组
}

// WhereQueryItemGroup 用于表示一组条件，可以是and/or连接的多个条件项或子条件组
type WhereQueryItemGroup struct {
	AndOr  string                `json:"andOr"`  // 连接条件 and/or
	Items  []WhereQueryItem      `json:"items"`  // 条件项
	Groups []WhereQueryItemGroup `json:"groups"` // 子条件组
}
type WhereQueryItem struct {
	AndOr string `json:"andOr"` // 连接条件 and/or
	Field string `json:"field"` // 字段名
	Opr   string `json:"opr"`   // 操作类型
	Value any    `json:"value"` // 值
}

// orderQuery
// OrderQueryItem 用于指定排序的字段和排序方式
type OrderQueryItem struct {
	Field string `json:"field"` // 字段名
	Order string `json:"order"` // 排序方式 asc desc 等
}

type Query struct {
	Pagination  Pagination        `json:"pagination"`
	SelectQuery []SelectQueryItem `json:"selectQuery"`

	WhereQuery WhereQuery       `json:"whereQuery"`
	OrQuery    []WhereQueryItem `json:"orQuery"`

	OrderQuery []OrderQueryItem `json:"orderQuery"`

	SkipCount bool `json:"skipCount"`
}

type Range struct {
	Start any `json:"start"`
	End   any `json:"end"`
}

const (
	// query
	TypeCustom       = "custom" // 用户自定义gorm scope，不做任何处理
	TypeEq           = "eq"
	TypeNe           = "ne"
	TypeGt           = "gt"
	TypeGte          = "gte"
	TypeLt           = "lt"
	TypeLte          = "lte"
	TypeLike         = "like"
	TypeNotLike      = "not_like"
	TypeStartLike    = "start_like"
	TypeEndLike      = "end_like"
	TypeIn           = "in"
	TypeNotIn        = "not_in"
	TypeIs           = "is"
	TypeIsNot        = "is_not"
	TypeIsEmpty      = "is_empty"     // is null / is '' / is ' '
	TypeIsNotEmpty   = "is_not_empty" // is not null / is not '' / is not ' '
	TypeStringLenEq  = "str_len_eq"
	TypeStringLenNe  = "str_len_ne"
	TypeStringLenGt  = "str_len_gt"
	TypeStringLenGte = "str_len_gte"
	TypeStringLenLt  = "str_len_lt"
	TypeStringLenLte = "str_len_lte"

	// order
	TypeAsc           = "asc"
	TypeDesc          = "desc"
	TypeStringLenAsc  = "str_len_asc"
	TypeStringLenDesc = "str_len_desc"
)

func ResolveQuery(query Query, model schema.Tabler, condition DbCondition) (err error) {
	tblName := model.TableName()
	fieldColMapper := CamelToSnakeFromStruct(model)
	if err = ResolveSelectQuery(query.SelectQuery, tblName, fieldColMapper, condition); err != nil {
		return
	}
	// 把简单查询也放到groups中
	whereQueryGroups := query.WhereQuery.Groups
	if len(query.WhereQuery.Items) > 0 {
		whereQueryGroups = append(whereQueryGroups, WhereQueryItemGroup{
			AndOr: "and",
			Items: query.WhereQuery.Items,
		})
	}
	if err = ResolveWhereQuery(whereQueryGroups, tblName, fieldColMapper, condition); err != nil {
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
				// condition.SetDistinct(fmt.Sprintf("\"%s\".\"%s\"", tblName, col)) 似乎带上table的写法会失效
				// 好吧，没有tblName也不行
				// gorm这块没做好
				condition.SetDistinct(col)
			} else {
				condition.SetSelect(fmt.Sprintf("\"%s\".\"%s\"", tblName, col))
			}
		} else {
			return fmt.Errorf("field %s not found", item.Field)
		}
	}
	return
}

func ResolveWhereQuery(groups []WhereQueryItemGroup, tblName string, fieldColMapper map[string]string, condition DbCondition) (err error) {
	for _, group := range groups {
		var dbGroup DbConditionWhereGroup
		if dbGroup, err = ResolveWhereQueryGroup(group, tblName, fieldColMapper); err != nil {
			return
		}
		condition.SetWhere(dbGroup)
	}
	return
}

func ResolveWhereQueryGroup(group WhereQueryItemGroup, tblName string, fieldColMapper map[string]string) (dbGroup DbConditionWhereGroup, err error) {
	dbGroup.AndOr = group.AndOr
	if dbGroup.AndOr == "" {
		dbGroup.AndOr = "and" // 默认使用and连接
	}
	if dbGroup.Items, err = ResolveWhereQueryItems(group.Items, tblName, fieldColMapper); err != nil {
		return
	}
	for _, subGroup := range group.Groups {
		var dbSubGroup DbConditionWhereGroup
		if dbSubGroup, err = ResolveWhereQueryGroup(subGroup, tblName, fieldColMapper); err != nil {
			return
		}
		dbGroup.Groups = append(dbGroup.Groups, dbSubGroup)
	}
	return
}

func ResolveWhereQueryItems(items []WhereQueryItem, tblName string, fieldColMapper map[string]string) (dbItems []DbConditionWhereItem, err error) {
	for _, item := range items {
		if utils.IsZeroValue(item.Value) || item.Opr == TypeCustom {
			continue
		}
		if col, ok := fieldColMapper[item.Field]; ok {
			fullColName := fmt.Sprintf("\"%s\".\"%s\"", tblName, col)
			dbItem := DbConditionWhereItem{AndOr: item.AndOr, Value: item.Value}
			switch item.Opr {
			case TypeEq:
				dbItem.Key = fmt.Sprintf("%s = ?", fullColName)
			case TypeNe:
				dbItem.Key = fmt.Sprintf("%s != ?", fullColName)
			case TypeGt:
				dbItem.Key = fmt.Sprintf("%s > ?", fullColName)
			case TypeGte:
				dbItem.Key = fmt.Sprintf("%s >= ?", fullColName)
			case TypeLt:
				dbItem.Key = fmt.Sprintf("%s < ?", fullColName)
			case TypeLte:
				dbItem.Key = fmt.Sprintf("%s <= ?", fullColName)
			case TypeLike:
				dbItem.Key = fmt.Sprintf("%s like ?", fullColName)
				dbItem.Value = "%" + item.Value.(string) + "%"
			case TypeNotLike:
				dbItem.Key = fmt.Sprintf("%s not like ?", fullColName)
				dbItem.Value = "%" + item.Value.(string) + "%"
			case TypeStartLike:
				dbItem.Key = fmt.Sprintf("%s like ?", fullColName)
				dbItem.Value = item.Value.(string) + "%"
			case TypeEndLike:
				dbItem.Key = fmt.Sprintf("%s like ?", fullColName)
				dbItem.Value = "%" + item.Value.(string)
			case TypeIn:
				dbItem.Key = fmt.Sprintf("%s in (?)", fullColName)
			case TypeNotIn:
				dbItem.Key = fmt.Sprintf("%s not in (?)", fullColName)
			case TypeIs:
				dbItem.Key = fmt.Sprintf("%s is ?", fullColName)
			case TypeIsNot:
				dbItem.Key = fmt.Sprintf("%s is not ?", fullColName)
			case TypeIsEmpty:
				dbItem.Key = fmt.Sprintf("TRIM(COALESCE(%s, '')) = ''", fullColName)
			case TypeIsNotEmpty:
				dbItem.Key = fmt.Sprintf("TRIM(COALESCE(%s, '')) != ''", fullColName)
			case TypeStringLenEq:
				dbItem.Key = fmt.Sprintf("length(%s) = ?", fullColName)
			case TypeStringLenNe:
				dbItem.Key = fmt.Sprintf("length(%s) != ?", fullColName)
			case TypeStringLenGt:
				dbItem.Key = fmt.Sprintf("length(%s) > ?", fullColName)
			case TypeStringLenGte:
				dbItem.Key = fmt.Sprintf("length(%s) >= ?", fullColName)
			case TypeStringLenLt:
				dbItem.Key = fmt.Sprintf("length(%s) < ?", fullColName)
			case TypeStringLenLte:
				dbItem.Key = fmt.Sprintf("length(%s) <= ?", fullColName)
			default:
				err = fmt.Errorf("opr %s not found", item.Opr)
				return
			}
			dbItems = append(dbItems, dbItem)
		} else {
			err = fmt.Errorf("field %s not found", item.Field)
			return
		}
	}
	return
}

func ResolveOrderQuery(items []OrderQueryItem, tblName string, fieldColMapper map[string]string, condition DbCondition) (err error) {
	for _, item := range items {
		if col, ok := fieldColMapper[item.Field]; ok {
			switch strings.ToLower(item.Order) {
			case TypeAsc:
				condition.SetOrder(fmt.Sprintf("\"%s\".\"%s\" asc", tblName, col))
			case TypeDesc:
				condition.SetOrder(fmt.Sprintf("\"%s\".\"%s\" desc", tblName, col))
			case TypeStringLenAsc:
				condition.SetOrder(fmt.Sprintf("length(\"%s\".\"%s\") asc", tblName, col))
			case TypeStringLenDesc:
				condition.SetOrder(fmt.Sprintf("length(\"%s\".\"%s\") desc", tblName, col))
			default:
				return fmt.Errorf("opr %s not found", item.Order)
			}
		} else {
			return fmt.Errorf("field %s not found", item.Field)
		}
	}
	return
}
