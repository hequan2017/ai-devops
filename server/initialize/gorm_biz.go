package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/biz"
)

// bizModel 业务模块自动迁移，新增业务表在此登记
func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(
		biz.Server{},
	)
	if err != nil {
		return err
	}
	return nil
}
