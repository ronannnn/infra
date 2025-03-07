package user

import "gorm.io/gorm"

func userPreload() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.
			Preload("Roles").
			Preload("Roles.Menus").
			Preload("Menus")
		return db
	}
}
