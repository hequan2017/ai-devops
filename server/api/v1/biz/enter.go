package biz

import (
	"ai-devops/server/service"
)

// ApiGroup 业务接口聚合，新增业务接口在此嵌入
type ApiGroup struct {
	ServerApi
	DockerApi
	K8sApi
	ReleaseApi
	TerminalApi
}

var serverService = service.ServiceGroupApp.BizServiceGroup.ServerService

var dockerHostService = service.ServiceGroupApp.BizServiceGroup.DockerHostService

var k8sClusterService = service.ServiceGroupApp.BizServiceGroup.K8sClusterService

var releaseService = service.ServiceGroupApp.BizServiceGroup.ReleaseService

var terminalService = service.ServiceGroupApp.BizServiceGroup.TerminalService
