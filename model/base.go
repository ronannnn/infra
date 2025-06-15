package model

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ronannnn/infra/constant"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/optimisticlock"
)

var (
	ErrUpdatedId = errors.New("id in updated or deleted entity must be greater than 0")
	ErrModified  = func(model string) error {
		return fmt.Errorf("[%s] data modified by others, please refresh the page", model)
	}
)

type Base struct {
	BaseModel // 包含Id等基础字段
	OprBy     // 包含创建和更新等有User Scopes的数据
}

func (b Base) GetId() uint {
	return b.Id
}

func (b Base) TableName() string {
	return ""
}

func (b Base) WithOprFromReq(r *http.Request) Crudable {
	oprId := r.Context().Value(constant.CtxKeyUserId)
	if oprId != nil {
		if convertedOprId, ok := oprId.(uint); ok {
			b.CreatedBy = convertedOprId
			b.UpdatedBy = convertedOprId
		}
	}
	return b
}
func (b Base) WithUpdaterFromReq(r *http.Request) Crudable {
	oprId := r.Context().Value(constant.CtxKeyUserId)
	if oprId != nil {
		if convertedOprId, ok := oprId.(uint); ok {
			b.UpdatedBy = convertedOprId
		}
	}
	return b
}

// 用于Services.Crud[T]接口
type Crudable interface {
	// 主要用于Crud的Update方法中判断Id是否为0
	GetId() uint
	// 获取操作人信息
	WithOprFromReq(r *http.Request) Crudable
	// 获取更新人信息
	WithUpdaterFromReq(r *http.Request) Crudable
	schema.Tabler
}

type BaseModel struct {
	Id        uint                   `json:"id" gorm:"primaryKey;autoIncrement:true"`
	CreatedAt time.Time              `json:"createdAt"`
	UpdatedAt time.Time              `json:"updatedAt"`
	DeletedAt gorm.DeletedAt         `gorm:"index" json:"-"`
	Version   optimisticlock.Version `json:"version"`
}

type OprBy struct {
	CreatedBy uint  `json:"createdBy"`
	Creator   *User `json:"creator" gorm:"foreignKey:CreatedBy"`
	UpdatedBy uint  `json:"updatedBy"`
	Updater   *User `json:"updater" gorm:"foreignKey:UpdatedBy"`
}

func GetOprFromReq(r *http.Request) OprBy {
	oprId := r.Context().Value(constant.CtxKeyUserId)
	if oprId != nil {
		if convertedOprId, ok := oprId.(uint); ok {
			return OprBy{
				CreatedBy: convertedOprId,
				UpdatedBy: convertedOprId,
			}
		}
	}
	return OprBy{}
}

func GetUpdaterFromReq(r *http.Request) OprBy {
	oprId := r.Context().Value(constant.CtxKeyUserId)
	if oprId != nil {
		if convertedOprId, ok := oprId.(uint); ok {
			return OprBy{
				UpdatedBy: convertedOprId,
			}
		}
	}
	return OprBy{}
}
