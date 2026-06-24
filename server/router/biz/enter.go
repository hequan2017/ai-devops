package biz

import (
	api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
)

// RouterGroup 业务路由聚合，新增业务路由在此嵌入
type RouterGroup struct {
	ServerRouter
}

var serverApi = api.ApiGroupApp.BizApiGroup.ServerApi
