package query

import "gorm.io/gorm"

// Paginate for gorm pagination scopes
func Paginate(pageNum, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if pageSize <= 0 || pageNum <= 0 {
			return db
		}
		offset := pageSize * (pageNum - 1)
		return db.Offset(offset).Limit(pageSize)
	}
}

func MakeConditionFromQuery(searchPayload any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		condition := &DbConditionImpl{}
		ResolveQuery(searchPayload, condition)
		return MakeCondition(condition)(db)
	}
}

func MakeCondition(condition *DbConditionImpl) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		for k, v := range condition.Where {
			db = db.Where(k, v...)
		}
		for k, v := range condition.Or {
			db = db.Or(k, v...)
		}
		for k, v := range condition.Not {
			db = db.Not(k, v...)
		}
		for _, o := range condition.Order {
			db = db.Order(o)
		}
		if len(condition.Select) > 0 {
			db = db.Select(condition.Select)
		}
		return db
	}
}
