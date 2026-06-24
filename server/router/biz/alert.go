package biz

import (
	"ai-devops/server/middleware"
	"github.com/gin-gonic/gin"
)

type AlertRouter struct{}

// InitAlertRouter 注册告警管理路由
func (a *AlertRouter) InitAlertRouter(Router *gin.RouterGroup) {
	alert := Router.Group("alert").Use(middleware.OperationRecord())
	alertWR := Router.Group("alert")
	{
		alert.POST("createAlertRule", alertApi.CreateAlertRule)
		alert.PUT("updateAlertRule", alertApi.UpdateAlertRule)
		alert.DELETE("deleteAlertRule", alertApi.DeleteAlertRule)
		alert.DELETE("deleteAlertRuleByIds", alertApi.DeleteAlertRuleByIds)
		alert.POST("resolveAlertRecord", alertApi.ResolveAlertRecord)
	}
	{
		alertWR.GET("findAlertRule", alertApi.FindAlertRule)
		alertWR.GET("getAlertRuleList", alertApi.GetAlertRuleList)
		alertWR.GET("getAlertRecordList", alertApi.GetAlertRecordList)
	}
}
