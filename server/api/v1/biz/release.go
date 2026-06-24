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

type ReleaseApi struct{}

// CreateRelease 新建发版
// @Router /release/createRelease [post]
func (r *ReleaseApi) CreateRelease(c *gin.Context) {
	var release biz.Release
	if err := c.ShouldBindJSON(&release); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	release.ApplicantID = utils.GetUserID(c)
	if err := releaseService.CreateRelease(release); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteRelease 删除
// @Router /release/deleteRelease [delete]
func (r *ReleaseApi) DeleteRelease(c *gin.Context) {
	var release biz.Release
	if err := c.ShouldBindJSON(&release); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(release.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := releaseService.DeleteRelease(release); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteReleaseByIds 批量删除
// @Router /release/deleteReleaseByIds [delete]
func (r *ReleaseApi) DeleteReleaseByIds(c *gin.Context) {
	var ids request.IdsReq
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := releaseService.DeleteReleaseByIds(ids); err != nil {
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateRelease 更新
// @Router /release/updateRelease [put]
func (r *ReleaseApi) UpdateRelease(c *gin.Context) {
	var release biz.Release
	if err := c.ShouldBindJSON(&release); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(release.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := releaseService.UpdateRelease(release); err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindRelease 详情
// @Router /release/findRelease [get]
func (r *ReleaseApi) FindRelease(c *gin.Context) {
	var release biz.Release
	if err := c.ShouldBindQuery(&release); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(release.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := releaseService.GetRelease(release.ID)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetReleaseList 分页列表
// @Router /release/getReleaseList [get]
func (r *ReleaseApi) GetReleaseList(c *gin.Context) {
	var pageInfo bizReq.ReleaseSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := releaseService.GetReleaseList(pageInfo)
	if err != nil {
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// releaseActionReq 发版审批/执行入参
type releaseActionReq struct {
	ID     uint   `json:"ID"`
	Reason string `json:"reason"`
}

// SubmitRelease 提交审批
// @Router /release/submitRelease [post]
func (r *ReleaseApi) SubmitRelease(c *gin.Context) {
	var req releaseActionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := releaseService.Submit(req.ID); err != nil {
		response.FailWithMessage("提交失败", c)
		return
	}
	response.OkWithMessage("已提交审批", c)
}

// ApproveRelease 审批通过
// @Router /release/approveRelease [post]
func (r *ReleaseApi) ApproveRelease(c *gin.Context) {
	var req releaseActionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := releaseService.Approve(req.ID, utils.GetUserID(c)); err != nil {
		response.FailWithMessage("审批失败", c)
		return
	}
	response.OkWithMessage("审批通过", c)
}

// RejectRelease 审批拒绝
// @Router /release/rejectRelease [post]
func (r *ReleaseApi) RejectRelease(c *gin.Context) {
	var req releaseActionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := releaseService.Reject(req.ID, utils.GetUserID(c), req.Reason); err != nil {
		response.FailWithMessage("操作失败", c)
		return
	}
	response.OkWithMessage("已拒绝", c)
}

// ExecuteRelease 执行发布
// @Router /release/executeRelease [post]
func (r *ReleaseApi) ExecuteRelease(c *gin.Context) {
	var req releaseActionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	result, err := releaseService.Execute(req.ID)
	if err != nil {
		response.FailWithMessage("执行失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]interface{}{"result": result}, "执行完成", c)
}
