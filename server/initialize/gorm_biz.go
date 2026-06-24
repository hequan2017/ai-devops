package initialize

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
)

// bizModel 业务模块自动迁移，新增业务表在此登记
func bizModel() error {
	db := global.GVA_DB
	err := db.AutoMigrate(
		biz.Server{},
		biz.DockerHost{},
		biz.K8sCluster{},
		biz.Release{},
		biz.ServerMetric{},
		biz.AlertRule{},
		biz.AlertRecord{},
	)
	if err != nil {
		return err
	}
	return nil
}
