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

type DockerApi struct{}

// CreateDockerHost 新增 Docker 接入点
// @Tags Docker
// @Router /docker/createDockerHost [post]
func (d *DockerApi) CreateDockerHost(c *gin.Context) {
	var host biz.DockerHost
	if err := c.ShouldBindJSON(&host); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := dockerHostService.CreateDockerHost(host); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteDockerHost 删除接入点
// @Router /docker/deleteDockerHost [delete]
func (d *DockerApi) DeleteDockerHost(c *gin.Context) {
	var host biz.DockerHost
	if err := c.ShouldBindJSON(&host); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(host.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := dockerHostService.DeleteDockerHost(host); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteDockerHostByIds 批量删除
// @Router /docker/deleteDockerHostByIds [delete]
func (d *DockerApi) DeleteDockerHostByIds(c *gin.Context) {
	var ids request.IdsReq
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := dockerHostService.DeleteDockerHostByIds(ids); err != nil {
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateDockerHost 更新接入点
// @Router /docker/updateDockerHost [put]
func (d *DockerApi) UpdateDockerHost(c *gin.Context) {
	var host biz.DockerHost
	if err := c.ShouldBindJSON(&host); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(host.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := dockerHostService.UpdateDockerHost(host); err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindDockerHost 查询接入点详情
// @Router /docker/findDockerHost [get]
func (d *DockerApi) FindDockerHost(c *gin.Context) {
	var host biz.DockerHost
	if err := c.ShouldBindQuery(&host); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(host.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := dockerHostService.GetDockerHost(host.ID)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetDockerHostList 分页列表
// @Router /docker/getDockerHostList [get]
func (d *DockerApi) GetDockerHostList(c *gin.Context) {
	var pageInfo bizReq.DockerHostSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := dockerHostService.GetDockerHostList(pageInfo)
	if err != nil {
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// TestDockerHost 连接测试
// @Router /docker/testDockerHost [get]
func (d *DockerApi) TestDockerHost(c *gin.Context) {
	var host biz.DockerHost
	if err := c.ShouldBindQuery(&host); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(host.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	version, err := dockerHostService.TestDockerHost(host.ID)
	if err != nil {
		response.FailWithMessage("连接失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]interface{}{"version": version}, "连接成功", c)
}

// GetContainers 容器列表
// @Router /docker/getContainers [get]
func (d *DockerApi) GetContainers(c *gin.Context) {
	var host biz.DockerHost
	_ = c.ShouldBindQuery(&host)
	list, err := dockerHostService.ListContainers(host.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// dockerContainerActionReq 容器操作入参
type dockerContainerActionReq struct {
	HostID      uint   `json:"hostId"`
	ContainerID string `json:"containerId"`
	Action      string `json:"action"` // start / stop / restart / remove
}

// ContainerAction 容器操作
// @Router /docker/containerAction [post]
func (d *DockerApi) ContainerAction(c *gin.Context) {
	var req dockerContainerActionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := dockerHostService.ContainerAction(req.HostID, req.ContainerID, req.Action); err != nil {
		response.FailWithMessage("操作失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("操作成功", c)
}

// GetContainerLogs 容器日志
// @Router /docker/getContainerLogs [get]
func (d *DockerApi) GetContainerLogs(c *gin.Context) {
	var host biz.DockerHost
	_ = c.ShouldBindQuery(&host)
	containerID := c.Query("containerId")
	logs, err := dockerHostService.ContainerLogs(host.ID, containerID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]interface{}{"logs": logs}, "获取成功", c)
}

// GetImages 镜像列表
// @Router /docker/getImages [get]
func (d *DockerApi) GetImages(c *gin.Context) {
	var host biz.DockerHost
	_ = c.ShouldBindQuery(&host)
	list, err := dockerHostService.ListImages(host.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetNetworks 网络列表
// @Router /docker/getNetworks [get]
func (d *DockerApi) GetNetworks(c *gin.Context) {
	var host biz.DockerHost
	_ = c.ShouldBindQuery(&host)
	list, err := dockerHostService.ListNetworks(host.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetVolumes 数据卷列表
// @Router /docker/getVolumes [get]
func (d *DockerApi) GetVolumes(c *gin.Context) {
	var host biz.DockerHost
	_ = c.ShouldBindQuery(&host)
	list, err := dockerHostService.ListVolumes(host.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}
