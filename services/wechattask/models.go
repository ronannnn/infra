package wechattask

import "github.com/ronannnn/infra/models"

type WechatTask struct {
	models.Base
	Uuid              *string            `json:"uuid" gorm:"type:varchar(36);uniqueIndex"` // 任务唯一标识
	Name              *string            `json:"name"`                                     // 任务名称
	Disabled          *bool              `json:"disabled"`                                 // 是否禁用该任务
	WechatTaskUserIds *WechatTaskUserIds `json:"wechatTaskUserIds" gorm:"many2many:wechat_tasks_user_ids_idx"`
}

type WechatTaskUserIds []WechatTaskUserId

func (w *WechatTaskUserIds) UserIds() []string {
	var userIds []string
	for _, user := range *w {
		if user.Disabled != nil && *user.Disabled {
			continue
		}
		if user.UserId != nil {
			userIds = append(userIds, *user.UserId)
		}
	}
	return userIds
}

func (w *WechatTaskUserIds) GroupIds() []string {
	var userIds []string
	for _, user := range *w {
		if user.Disabled != nil && *user.Disabled {
			continue
		}
		if user.UserId != nil && user.UserType != nil && *user.UserType == UserTypeGroup {
			userIds = append(userIds, *user.UserId)
		}
	}
	return userIds
}

func (w *WechatTaskUserIds) FriendIds() []string {
	var userIds []string
	for _, user := range *w {
		if user.Disabled != nil && *user.Disabled {
			continue
		}
		if user.UserId != nil && user.UserType != nil && *user.UserType == UserTypeFriend {
			userIds = append(userIds, *user.UserId)
		}
	}
	return userIds
}

func (w *WechatTaskUserIds) MpIds() []string {
	var userIds []string
	for _, user := range *w {
		if user.Disabled != nil && *user.Disabled {
			continue
		}
		if user.UserId != nil && user.UserType != nil && *user.UserType == UserTypeMp {
			userIds = append(userIds, *user.UserId)
		}
	}
	return userIds
}

type UserType int

const (
	UserTypeFriend UserType = iota + 1
	UserTypeGroup
	UserTypeMp
)

type WechatTaskUserId struct {
	models.Base
	Nickname    *string       `json:"nickname"` // 微信用户昵称
	UserId      *string       `json:"userId"`   // 微信用户唯一标识
	UserType    *UserType     `json:"userType"` // 微信用户类型: 好友，群组
	Disabled    *bool         `json:"disabled"` // 是否禁用该用户
	WechatTasks *[]WechatTask `json:"wechatTasks" gorm:"many2many:wechat_tasks_user_ids_idx"`
}
