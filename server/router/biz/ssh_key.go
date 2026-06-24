package biz

import (
	"ai-devops/server/middleware"
	"github.com/gin-gonic/gin"
)

type SshKeyRouter struct{}

func (a *SshKeyRouter) InitSshKeyRouter(Router *gin.RouterGroup) {
	sshKey := Router.Group("sshKey").Use(middleware.OperationRecord())
	sshKeyWR := Router.Group("sshKey")
	{
		sshKey.POST("createSshKey", sshKeyApi.CreateSshKey)
		sshKey.PUT("updateSshKey", sshKeyApi.UpdateSshKey)
		sshKey.DELETE("deleteSshKey", sshKeyApi.DeleteSshKey)
		sshKey.DELETE("deleteSshKeyByIds", sshKeyApi.DeleteSshKeyByIds)
	}
	{
		sshKeyWR.GET("findSshKey", sshKeyApi.FindSshKey)
		sshKeyWR.GET("getSshKeyList", sshKeyApi.GetSshKeyList)
		sshKeyWR.GET("getAllSshKey", sshKeyApi.GetAllSshKey)
	}
}
