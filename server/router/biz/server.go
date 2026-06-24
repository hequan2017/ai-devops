package biz

import (
	"ai-devops/server/middleware"
	"github.com/gin-gonic/gin"
)

type ServerRouter struct{}

// InitServerRouter 注册服务器管理路由
func (s *ServerRouter) InitServerRouter(Router *gin.RouterGroup) {
	serverRouter := Router.Group("server").Use(middleware.OperationRecord())
	serverRouterWithoutRecord := Router.Group("server")
	{
		serverRouter.POST("createServer", serverApi.CreateServer)           // 新增服务器
		serverRouter.PUT("updateServer", serverApi.UpdateServer)            // 更新服务器
		serverRouter.DELETE("deleteServer", serverApi.DeleteServer)         // 删除服务器
		serverRouter.DELETE("deleteServerByIds", serverApi.DeleteServerByIds) // 批量删除服务器
		serverRouter.POST("ipmiPower", serverApi.IPMIPower)                 // IPMI 电源控制
		serverRouter.POST("execCmd", serverApi.ExecCmd)                     // 批量执行命令
	}
	{
		serverRouterWithoutRecord.GET("findServer", serverApi.FindServer)     // 获取服务器详情
		serverRouterWithoutRecord.GET("getServerList", serverApi.GetServerList) // 获取服务器列表
		serverRouterWithoutRecord.GET("getServerMetrics", serverApi.GetServerMetrics)           // 性能历史
		serverRouterWithoutRecord.GET("getServerMetricLatest", serverApi.GetServerMetricLatest) // 最新指标
	}
}
