package models

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

var (
	ErrUpdatedId = errors.New("id in updated or deleted entity must be greater than 0")
	ErrModified  = func(model string) error {
		return fmt.Errorf("[%s] data modified by others, please refresh the page", model)
	}
)

type BaseModel struct {
	Id        uint                   `json:"id" gorm:"primaryKey;autoIncrement:true"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	DeletedAt gorm.DeletedAt         `gorm:"index" json:"-"`
	Version   optimisticlock.Version `json:"version"`
}

func (b *BaseModel) SetId(id uint) {
	b.Id = id
}

type Base struct {
	BaseModel
	OprBy
}

type OprBy struct {
	CreatedBy uint  `json:"createdBy"`
	Creator   *User `json:"creator" gorm:"foreignKey:CreatedBy"`
	UpdatedBy uint  `json:"updatedBy"`
	Updater   *User `json:"updater" gorm:"foreignKey:UpdatedBy"`
}

func (o *OprBy) GetOprFromReq(r *http.Request) {
	oprId := r.Context().Value(CtxKeyUserId)
	if oprId != nil {
		if convertedOprId, ok := oprId.(uint); ok {
			o.CreatedBy = convertedOprId
			o.UpdatedBy = convertedOprId
		}
	}
}

func (o *OprBy) GetUpdaterFromReq(r *http.Request) {
	oprId := r.Context().Value(CtxKeyUserId)
	if oprId != nil {
		if convertedOprId, ok := oprId.(uint); ok {
			o.UpdatedBy = convertedOprId
		}
	}
}
