package wechattask

import "gorm.io/gorm"

func wechatTaskPreload() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Preload("WechatTaskUserIds")
		return db
	}
}
