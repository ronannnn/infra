package model

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ronannnn/infra/constant"
	"gorm.io/gorm"
	"gorm.io/plugin/optimisticlock"
)

var (
	ErrWrongUsernameOrPassword = fmt.Errorf("incorrect username or password")
)

type User struct {
	Id            uint                   `json:"id" gorm:"primaryKey;autoIncrement:true"`
	CreatedAt     time.Time              `json:"createdAt"`
	UserCreatedBy uint                   `json:"userCreatedBy"`
	UserCreator   *User                  `json:"userCreator" gorm:"foreignKey:UserCreatedBy;references:Id"`
	UpdatedAt     time.Time              `json:"updatedAt"`
	UserUpdatedBy uint                   `json:"userUpdatedBy"`
	UserUpdater   *User                  `json:"userUpdater" gorm:"foreignKey:UserCreatedBy;references:Id"`
	DeletedAt     gorm.DeletedAt         `gorm:"index" json:"-"`
	Version       optimisticlock.Version `json:"version"`

	// user info
	Nickname     *string     `json:"nickname"`
	JobTitleId   *uint       `json:"jobTitleId"` // 职务ID
	JobTitle     *JobTitle   `json:"jobTitle" gorm:"foreignKey:JobTitleId;references:Id"`
	JobGradeId   *uint       `json:"jobGradeId"` // 职级ID
	JobGrade     *JobGrade   `json:"jobGrade" gorm:"foreignKey:JobGradeId;references:Id"`
	DepartmentId *uint       `json:"departmentId"` // 部门ID
	Department   *Department `json:"department" gorm:"foreignKey:DepartmentId;references:Id"`
	EntryDate    *time.Time  `json:"entryDate"`  // 入职日期
	ResignDate   *time.Time  `json:"resignDate"` // 离职日期
	Gender       *int        `json:"gender"`     // 性别，1表示男，2表示女

	// login info
	Username *string `json:"username"`
	Email    *string `json:"email"`
	TelNo    *string `json:"telNo"`
	Password *string `json:"-"`

	// 登陆方式：可自定义，比如2代WMS登录，账户密码登录，手机验证码登录等
	// 由逗号分隔
	LoginType *string `json:"loginType"`

	// wechat info
	WechatOpenId  *string `json:"wechatOpenId"`
	WechatUnionId *string `json:"wechatUnionId"`

	// permission
	Roles []*Role `json:"roles" gorm:"many2many:user_roles;comment:用户角色"`
	Menus []*Menu `json:"menus" gorm:"many2many:user_menus;"`
}

func (u User) TableName() string {
	return "users"
}

func (u User) GetId() uint {
	return u.Id
}

func (u *User) GetOprFromReq(r *http.Request) {
	oprId := r.Context().Value(constant.CtxKeyUserId)
	if oprId != nil {
		if convertedOprId, ok := oprId.(uint); ok {
			u.UserCreatedBy = convertedOprId
			u.UserUpdatedBy = convertedOprId
		}
	}
}

func (u *User) GetUpdaterFromReq(r *http.Request) {
	oprId := r.Context().Value(constant.CtxKeyUserId)
	if oprId != nil {
		if convertedOprId, ok := oprId.(uint); ok {
			u.UserUpdatedBy = convertedOprId
		}
	}
}

func (u User) HasLoginType(loginType string) error {
	if u.LoginType == nil {
		return fmt.Errorf("user login type is not defined")
	}
	allowedLoginTypes := strings.Split(*u.LoginType, ",")
	for _, allowedLoginType := range allowedLoginTypes {
		if allowedLoginType == loginType {
			return nil
		}
	}
	return fmt.Errorf("user login type %s is not allowed", loginType)
}
