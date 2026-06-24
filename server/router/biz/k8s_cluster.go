package biz

import (
	"ai-devops/server/middleware"
	"github.com/gin-gonic/gin"
)

type K8sRouter struct{}

// InitK8sRouter 注册 K8s 管理路由
func (k *K8sRouter) InitK8sRouter(Router *gin.RouterGroup) {
	k8sRouter := Router.Group("k8s").Use(middleware.OperationRecord())
	k8sRouterWithoutRecord := Router.Group("k8s")
	{
		k8sRouter.POST("createK8sCluster", k8sApi.CreateK8sCluster)
		k8sRouter.PUT("updateK8sCluster", k8sApi.UpdateK8sCluster)
		k8sRouter.DELETE("deleteK8sCluster", k8sApi.DeleteK8sCluster)
		k8sRouter.DELETE("deleteK8sClusterByIds", k8sApi.DeleteK8sClusterByIds)
		k8sRouter.POST("deletePod", k8sApi.DeletePod)
	}
	{
		k8sRouterWithoutRecord.GET("findK8sCluster", k8sApi.FindK8sCluster)
		k8sRouterWithoutRecord.GET("getK8sClusterList", k8sApi.GetK8sClusterList)
		k8sRouterWithoutRecord.GET("testK8sCluster", k8sApi.TestK8sCluster)
		k8sRouterWithoutRecord.GET("getNamespaces", k8sApi.GetNamespaces)
		k8sRouterWithoutRecord.GET("getPods", k8sApi.GetPods)
		k8sRouterWithoutRecord.GET("getNodes", k8sApi.GetNodes)
		k8sRouterWithoutRecord.GET("getDeployments", k8sApi.GetDeployments)
		k8sRouterWithoutRecord.GET("getServices", k8sApi.GetServices)
		k8sRouterWithoutRecord.GET("getPodLogs", k8sApi.GetPodLogs)
	}
}
