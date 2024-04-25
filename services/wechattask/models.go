package wechattask

import "github.com/ronannnn/infra/models"

type WechatTask struct {
	models.Base
	Uuid          *string            `json:"uuid" gorm:"type:varchar(255);uniqueIndex"` // 任务唯一标识
	Name          *string            `json:"name"`                                      // 任务名称
	Disabled      *bool              `json:"disabled"`                                  // 是否禁用该任务
	WechatUserIds *WechatTaskUserIds `json:"wechatUserIds" gorm:"many2many:wechat_tasks_user_ids"`
}

type WechatTaskUserIds []WechatTaskUserId

func (w *WechatTaskUserIds) UserIds() []string {
	var userIds []string
	for _, user := range *w {
		if user.Disabled != nil && *user.Disabled {
			continue
		}
		if user.WechatUserId != nil {
			userIds = append(userIds, *user.WechatUserId)
		}
	}
	return userIds
}

type WechatTaskUserId struct {
	models.Base
	Nickname     *string       `json:"nickname"`     // 微信用户昵称
	WechatUserId *string       `json:"wechatUserId"` // 微信用户唯一标识
	Disabled     *bool         `json:"disabled"`     // 是否禁用该用户
	WechatTasks  *[]WechatTask `json:"wechatTasks" gorm:"many2many:wechat_tasks_user_ids"`
}
