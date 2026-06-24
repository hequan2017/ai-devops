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

type SshKeyApi struct{}

func (a *SshKeyApi) CreateSshKey(c *gin.Context) {
	var key biz.SshKey
	if err := c.ShouldBindJSON(&key); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sshKeyService.Create(key); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

func (a *SshKeyApi) DeleteSshKey(c *gin.Context) {
	var key biz.SshKey
	if err := c.ShouldBindJSON(&key); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(key.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sshKeyService.Delete(key); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (a *SshKeyApi) DeleteSshKeyByIds(c *gin.Context) {
	var ids request.IdsReq
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sshKeyService.DeleteByIds(ids); err != nil {
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

func (a *SshKeyApi) UpdateSshKey(c *gin.Context) {
	var key biz.SshKey
	if err := c.ShouldBindJSON(&key); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(key.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := sshKeyService.Update(key); err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *SshKeyApi) FindSshKey(c *gin.Context) {
	var key biz.SshKey
	if err := c.ShouldBindQuery(&key); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(key.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := sshKeyService.Get(key.ID)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

func (a *SshKeyApi) GetSshKeyList(c *gin.Context) {
	var pageInfo bizReq.SshKeySearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sshKeyService.GetList(pageInfo)
	if err != nil {
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetAllSshKey 全部密钥（下拉用）
func (a *SshKeyApi) GetAllSshKey(c *gin.Context) {
	list, err := sshKeyService.GetAll()
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}
