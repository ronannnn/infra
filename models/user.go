package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

var (
	ErrWrongUsernameOrPassword = fmt.Errorf("incorrect username or password")
)

type User struct {
	Id        uint                   `json:"id" gorm:"primaryKey;autoIncrement:true"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	DeletedAt gorm.DeletedAt         `gorm:"index" json:"-"`
	Version   optimisticlock.Version `json:"version"`
	// user info
	Nickname *string `json:"nickname"`
	// login info
	Username *string `json:"username"`
	Email    *string `json:"email"`
	TelNo    *string `json:"telNo"`
	Password *string `json:"-"`
}
