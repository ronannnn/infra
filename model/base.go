package model

import (
	"context"
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
	CreatedBy *uint `json:"createdBy"`
	Creator   *User `json:"creator" gorm:"foreignKey:CreatedBy"`
	UpdatedBy *uint `json:"updatedBy"`
	Updater   *User `json:"updater" gorm:"foreignKey:UpdatedBy"`
}

func GetOprFromReq(r *http.Request) OprBy {
	return GetOprFromCtx(r.Context())
}
func GetOprFromCtx(ctx context.Context) OprBy {
	oprId := ctx.Value(constant.CtxKeyUserId)
	if oprId != nil {
		if convertedOprId, ok := oprId.(uint); ok {
			return OprBy{
				CreatedBy: &convertedOprId,
				UpdatedBy: &convertedOprId,
			}
		}
	}
	return OprBy{}
}

func GetUpdaterFromReq(r *http.Request) OprBy {
	return GetUpdaterFromCtx(r.Context())
}
func GetUpdaterFromCtx(ctx context.Context) OprBy {
	oprId := ctx.Value(constant.CtxKeyUserId)
	if oprId != nil {
		if convertedOprId, ok := oprId.(uint); ok {
			return OprBy{
				UpdatedBy: &convertedOprId,
			}
		}
	}
	return OprBy{}
}
