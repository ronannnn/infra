package models

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/ronannnn/infra/constant"
	"github.com/ronannnn/infra/models/request/query"
	"github.com/ronannnn/infra/utils"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

var (
	ErrUpdatedId = errors.New("id in updated or deleted entity must be greater than 0")
	ErrModified  = func(model string) error {
		return fmt.Errorf("[%s] data modified by others, please refresh the page", model)
	}
)

type Base struct {
	BaseModel          // 包含Id等基础字段
	OprBy              // 包含创建和更新等有User Scopes的数据
	DefaultQuerySetter // 包含默认的FieldColMapper，用于query
}

// 用于Services.Crud[T]接口
type Crudable interface {
	Identifiable
	query.QuerySetter
}

// 主要用于Crud的Update方法中判断Id是否为0
type Identifiable interface {
	GetId() uint
}

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

func (b BaseModel) GetId() uint {
	return b.Id
}

type OprBy struct {
	CreatedBy uint  `json:"createdBy"`
	Creator   *User `json:"creator" gorm:"foreignKey:CreatedBy"`
	UpdatedBy uint  `json:"updatedBy"`
	Updater   *User `json:"updater" gorm:"foreignKey:UpdatedBy"`
}

func (o *OprBy) GetOprFromReq(r *http.Request) {
	oprId := r.Context().Value(constant.CtxKeyUserId)
	if oprId != nil {
		if convertedOprId, ok := oprId.(uint); ok {
			o.CreatedBy = convertedOprId
			o.UpdatedBy = convertedOprId
		}
	}
}

func (o *OprBy) GetUpdaterFromReq(r *http.Request) {
	oprId := r.Context().Value(constant.CtxKeyUserId)
	if oprId != nil {
		if convertedOprId, ok := oprId.(uint); ok {
			o.UpdatedBy = convertedOprId
		}
	}
}

type DefaultQuerySetter struct{}

func (d DefaultQuerySetter) FieldColMapper() map[string]string {
	return CamelToSnakeFromStruct(d)
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

// utils

func CamelToSnakeWithBaseFromStrings(fields []string) map[string]string {
	mapper := make(map[string]string)
	for _, field := range fields {
		mapper[field] = utils.CamelToSnake(field)
	}
	return mapper
}

func CamelToSnakeFromStruct(obj any) map[string]string {
	fields := []string{}
	getJsonTagsFromStruct(obj, &fields)
	return CamelToSnakeWithBaseFromStrings(fields)
}

func getJsonTagsFromStruct(obj any, fields *[]string) {
	structType := reflect.TypeOf(obj)
	structValue := reflect.ValueOf(obj)
	for i := range structType.NumField() {
		jsonTag, jsonTagOk := structType.Field(i).Tag.Lookup("json")
		if !jsonTagOk {
			getJsonTagsFromStruct(structValue.Field(i).Interface(), fields)
			continue
		}
		if jsonTag == "-" {
			continue
		}
		gormTag, gormTagOk := structType.Field(i).Tag.Lookup("gorm")
		if gormTagOk && gormTag == "-" {
			continue
		}
		*fields = append(*fields, jsonTag)
	}
}
