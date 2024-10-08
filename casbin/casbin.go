package casbin

import (
	"time"

	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	gormAdapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

func NewCasbinEnforcer(
	db *gorm.DB,
) (enforcer *casbin.SyncedCachedEnforcer, err error) {
	var adapter *gormAdapter.Adapter
	if adapter, err = gormAdapter.NewAdapterByDB(db); err != nil {
		return
	}
	rbac_rule := `
	[request_definition]
	r = sub, obj, act
	
	[policy_definition]
	p = sub, obj, act
	
	[role_definition]
	g = _, _
	
	[policy_effect]
	e = some(where (p.eft == allow))
	
	[matchers]
	m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
	`
	var cModel casbinModel.Model
	if cModel, err = casbinModel.NewModelFromString(rbac_rule); err != nil {
		return
	}
	if enforcer, err = casbin.NewSyncedCachedEnforcer(cModel, adapter); err != nil {
		return
	}
	enforcer.SetExpireTime(60 * 60 * time.Second)
	if err = enforcer.LoadPolicy(); err != nil {
		return
	}
	// TODO set watcher
	return
}
