package query

import (
	"fmt"
	"strings"

	"github.com/ronannnn/infra/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Paginate for gorm pagination scopes
func Paginate(pagination Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pagination.PageSize <= 0 || pagination.PageNum <= 0 {
			return db
		}
		offset := pagination.PageSize * (pagination.PageNum - 1)
		return db.Offset(offset).Limit(pagination.PageSize)
	}
}

func MakeConditionFromQuery(
	query Query,
	model schema.Tabler,
	filter ConditionFilter,
) (fn func(db *gorm.DB) *gorm.DB, err error) {
	condition := &DbConditionImpl{}
	if err = ResolveQuery(query, model, condition, filter); err != nil {
		return
	}
	fn = func(db *gorm.DB) *gorm.DB {
		return MakeCondition(condition)(db)
	}
	return
}

func MakeCondition(condition *DbConditionImpl) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for _, group := range condition.Where {
			db = MakeDbConditionWhereQueryGroup(db, group)
		}
		for _, o := range condition.Order {
			db = db.Order(o)
		}
		if len(condition.Select) > 0 {
			db = db.Select(condition.Select)
		}
		if len(condition.Distinct) > 0 {
			db = db.Distinct(condition.Distinct)
		}
		return db
	}
}

func MakeDbConditionWhereQueryGroup(db *gorm.DB, group DbConditionWhereGroup) *gorm.DB {
	if len(group.Items) == 0 && len(group.Groups) == 0 {
		return db
	}

	// gorm有bug，如果db.Where(db.Where(), ...)，没有与第二个db.Where()有同级别的条件，会导致第二个db.Where()被忽略
	// 因此，如果没有Items的话，即没有同级别的条件，直接用传入db，不需要新建db.Session()
	if len(group.Items) == 0 {
		for _, subGroup := range group.Groups {
			db = MakeDbConditionWhereQueryGroup(db, subGroup)
		}
	} else {
		whereSubDb := db.Session(&gorm.Session{NewDB: true})
		whereSubDb = MakeDbConditionWhereQueryItems(whereSubDb, group.Items)
		for _, subGroup := range group.Groups {
			whereSubDb = MakeDbConditionWhereQueryGroup(whereSubDb, subGroup)
		}

		switch strings.ToLower(group.AndOr) {
		case "and", "":
			db = db.Where(whereSubDb)
		case "or":
			db = db.Or(whereSubDb)
		default:
			fmt.Printf("error: invalid and/or condition in where query group: %s", group.AndOr)
		}
	}

	return db
}

func MakeDbConditionWhereQueryItems(db *gorm.DB, items []DbConditionWhereItem) *gorm.DB {
	if len(items) == 0 {
		return db
	}
	for _, item := range items {
		switch strings.ToLower(item.AndOr) {
		case "and", "":
			if item.Value == nil {
				db = db.Where(item.Key)
			} else {
				db = db.Where(item.Key, item.Value)
			}
		case "or":
			if item.Value == nil {
				db = db.Or(item.Key)
			} else {
				db = db.Or(item.Key, item.Value)
			}
		default:
			fmt.Printf("error: invalid and/or condition in where query item: %s", item.AndOr)
		}
	}
	return db
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
