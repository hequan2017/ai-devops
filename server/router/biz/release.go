package biz

import (
	"ai-devops/server/middleware"
	"github.com/gin-gonic/gin"
)

type ReleaseRouter struct{}

// InitReleaseRouter 注册发版工作流路由
func (r *ReleaseRouter) InitReleaseRouter(Router *gin.RouterGroup) {
	releaseRouter := Router.Group("release").Use(middleware.OperationRecord())
	releaseRouterWithoutRecord := Router.Group("release")
	{
		releaseRouter.POST("createRelease", releaseApi.CreateRelease)
		releaseRouter.PUT("updateRelease", releaseApi.UpdateRelease)
		releaseRouter.DELETE("deleteRelease", releaseApi.DeleteRelease)
		releaseRouter.DELETE("deleteReleaseByIds", releaseApi.DeleteReleaseByIds)
		releaseRouter.POST("submitRelease", releaseApi.SubmitRelease)
		releaseRouter.POST("approveRelease", releaseApi.ApproveRelease)
		releaseRouter.POST("rejectRelease", releaseApi.RejectRelease)
		releaseRouter.POST("executeRelease", releaseApi.ExecuteRelease)
	}
	{
		releaseRouterWithoutRecord.GET("findRelease", releaseApi.FindRelease)
		releaseRouterWithoutRecord.GET("getReleaseList", releaseApi.GetReleaseList)
	}
}
