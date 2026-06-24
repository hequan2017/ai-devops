package biz

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service"
)

// ApiGroup 业务接口聚合，新增业务接口在此嵌入
type ApiGroup struct {
	ServerApi
}

var serverService = service.ServiceGroupApp.BizServiceGroup.ServerService
