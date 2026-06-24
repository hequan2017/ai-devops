package biz

import (
	api "ai-devops/server/api/v1"
)

// RouterGroup 业务路由聚合，新增业务路由在此嵌入
type RouterGroup struct {
	ServerRouter
	DockerRouter
	K8sRouter
	ReleaseRouter
	TerminalRouter
	AlertRouter
}

var serverApi = api.ApiGroupApp.BizApiGroup.ServerApi

var dockerApi = api.ApiGroupApp.BizApiGroup.DockerApi

var k8sApi = api.ApiGroupApp.BizApiGroup.K8sApi

var releaseApi = api.ApiGroupApp.BizApiGroup.ReleaseApi

var terminalApi = api.ApiGroupApp.BizApiGroup.TerminalApi

var alertApi = api.ApiGroupApp.BizApiGroup.AlertApi
