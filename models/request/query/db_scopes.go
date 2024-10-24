package query

import (
	"fmt"

	"github.com/ronannnn/infra/utils"
	"gorm.io/gorm"
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

func MakeConditionFromQuery(query Query, setter QuerySetter) (fn func(db *gorm.DB) *gorm.DB, err error) {
	condition := &DbConditionImpl{}
	if err = ResolveQuery(query, setter, condition); err != nil {
		return
	}
	fn = func(db *gorm.DB) *gorm.DB {
		return MakeCondition(condition)(db)
	}
	return
}

func MakeCondition(condition *DbConditionImpl) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// where
		whereSubDb := db.Session(&gorm.Session{NewDB: true})
		for k, vs := range condition.Where {
			for _, v := range vs {
				whereSubDb = whereSubDb.Where(k, v...)
			}
		}
		db = db.Where(whereSubDb)

		// or
		orSubDb := db.Session(&gorm.Session{NewDB: true})
		for k, vs := range condition.Or {
			for _, v := range vs {
				orSubDb = orSubDb.Or(k, v...)
			}
		}
		db = db.Where(orSubDb)

		// not
		notSubDb := db.Session(&gorm.Session{NewDB: true})
		for k, vs := range condition.Not {
			for _, v := range vs {
				notSubDb = notSubDb.Not(k, v...)
			}
		}
		db = db.Where(notSubDb)

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
