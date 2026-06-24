package biz

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
	bizReq "ai-devops/server/model/biz/request"
	"ai-devops/server/model/common/request"
	"ai-devops/server/model/common/response"
	"ai-devops/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AlertApi struct{}

// ===== 告警规则 =====

// CreateAlertRule 新建告警规则
// @Router /alert/createAlertRule [post]
func (a *AlertApi) CreateAlertRule(c *gin.Context) {
	var rule biz.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := alertService.CreateRule(rule); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteAlertRule 删除规则
// @Router /alert/deleteAlertRule [delete]
func (a *AlertApi) DeleteAlertRule(c *gin.Context) {
	var rule biz.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(rule.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := alertService.DeleteRule(rule); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteAlertRuleByIds 批量删除
// @Router /alert/deleteAlertRuleByIds [delete]
func (a *AlertApi) DeleteAlertRuleByIds(c *gin.Context) {
	var ids request.IdsReq
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := alertService.DeleteRuleByIds(ids); err != nil {
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateAlertRule 更新规则
// @Router /alert/updateAlertRule [put]
func (a *AlertApi) UpdateAlertRule(c *gin.Context) {
	var rule biz.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(rule.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := alertService.UpdateRule(rule); err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindAlertRule 查询规则
// @Router /alert/findAlertRule [get]
func (a *AlertApi) FindAlertRule(c *gin.Context) {
	var rule biz.AlertRule
	if err := c.ShouldBindQuery(&rule); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(rule.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := alertService.GetRule(rule.ID)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetAlertRuleList 规则列表
// @Router /alert/getAlertRuleList [get]
func (a *AlertApi) GetAlertRuleList(c *gin.Context) {
	var pageInfo bizReq.AlertRuleSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := alertService.GetRuleList(pageInfo)
	if err != nil {
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// ===== 告警记录 =====

// GetAlertRecordList 告警记录列表
// @Router /alert/getAlertRecordList [get]
func (a *AlertApi) GetAlertRecordList(c *gin.Context) {
	var pageInfo bizReq.AlertRecordSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := alertService.GetRecordList(pageInfo)
	if err != nil {
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// resolveReq 处理告警入参
type resolveReq struct {
	ID uint `json:"ID"`
}

// ResolveAlertRecord 标记告警已处理
// @Router /alert/resolveAlertRecord [post]
func (a *AlertApi) ResolveAlertRecord(c *gin.Context) {
	var req resolveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := alertService.ResolveRecord(req.ID); err != nil {
		response.FailWithMessage("处理失败", c)
		return
	}
	response.OkWithMessage("已处理", c)
}
