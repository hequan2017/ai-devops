package biz

import (
	"ai-devops/server/middleware"
	"github.com/gin-gonic/gin"
)

type DockerRouter struct{}

// InitDockerRouter 注册 Docker 管理路由
func (d *DockerRouter) InitDockerRouter(Router *gin.RouterGroup) {
	dockerRouter := Router.Group("docker").Use(middleware.OperationRecord())
	dockerRouterWithoutRecord := Router.Group("docker")
	{
		dockerRouter.POST("createDockerHost", dockerApi.CreateDockerHost)
		dockerRouter.PUT("updateDockerHost", dockerApi.UpdateDockerHost)
		dockerRouter.DELETE("deleteDockerHost", dockerApi.DeleteDockerHost)
		dockerRouter.DELETE("deleteDockerHostByIds", dockerApi.DeleteDockerHostByIds)
		dockerRouter.POST("containerAction", dockerApi.ContainerAction)
	}
	{
		dockerRouterWithoutRecord.GET("findDockerHost", dockerApi.FindDockerHost)
		dockerRouterWithoutRecord.GET("getDockerHostList", dockerApi.GetDockerHostList)
		dockerRouterWithoutRecord.GET("testDockerHost", dockerApi.TestDockerHost)
		dockerRouterWithoutRecord.GET("getContainers", dockerApi.GetContainers)
		dockerRouterWithoutRecord.GET("getContainerLogs", dockerApi.GetContainerLogs)
		dockerRouterWithoutRecord.GET("getImages", dockerApi.GetImages)
		dockerRouterWithoutRecord.GET("getNetworks", dockerApi.GetNetworks)
		dockerRouterWithoutRecord.GET("getVolumes", dockerApi.GetVolumes)
	}
}
