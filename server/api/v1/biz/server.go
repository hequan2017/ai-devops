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
	"strconv"
)

type ServerApi struct{}

// CreateServer
// @Tags      Server
// @Summary   新增服务器
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      biz.Server                     true  "服务器信息"
// @Success   200   {object}  response.Response{msg=string}  "新增服务器"
// @Router    /server/createServer [post]
func (s *ServerApi) CreateServer(c *gin.Context) {
	var server biz.Server
	if err := c.ShouldBindJSON(&server); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serverService.CreateServer(server); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteServer
// @Tags      Server
// @Summary   删除服务器
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      biz.Server                     true  "服务器ID"
// @Success   200   {object}  response.Response{msg=string}  "删除服务器"
// @Router    /server/deleteServer [delete]
func (s *ServerApi) DeleteServer(c *gin.Context) {
	var server biz.Server
	if err := c.ShouldBindJSON(&server); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(server.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serverService.DeleteServer(server); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteServerByIds
// @Tags      Server
// @Summary   批量删除服务器
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.IdsReq                 true  "服务器ID集合"
// @Success   200   {object}  response.Response{msg=string}  "批量删除服务器"
// @Router    /server/deleteServerByIds [delete]
func (s *ServerApi) DeleteServerByIds(c *gin.Context) {
	var ids request.IdsReq
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serverService.DeleteServerByIds(ids); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateServer
// @Tags      Server
// @Summary   更新服务器
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      biz.Server                     true  "服务器信息"
// @Success   200   {object}  response.Response{msg=string}  "更新服务器"
// @Router    /server/updateServer [put]
func (s *ServerApi) UpdateServer(c *gin.Context) {
	var server biz.Server
	if err := c.ShouldBindJSON(&server); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(server.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serverService.UpdateServer(server); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindServer
// @Tags      Server
// @Summary   根据ID获取服务器详情
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     biz.Server                     true  "服务器ID"
// @Success   200   {object}  response.Response{data=biz.Server,msg=string}  "根据ID获取服务器详情"
// @Router    /server/findServer [get]
func (s *ServerApi) FindServer(c *gin.Context) {
	var server biz.Server
	if err := c.ShouldBindQuery(&server); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(server.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := serverService.GetServer(server.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetServerList
// @Tags      Server
// @Summary   分页获取服务器列表
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  query     bizReq.ServerSearch            true  "页码, 每页大小, 厂商, 状态, 关键字"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "分页获取服务器列表"
// @Router    /server/getServerList [get]
func (s *ServerApi) GetServerList(c *gin.Context) {
	var pageInfo bizReq.ServerSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serverService.GetServerList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// ipmiPowerReq IPMI 电源控制入参
type ipmiPowerReq struct {
	ID     uint   `json:"ID"`
	Action string `json:"action"` // status / on / off / reset / soft
}

// IPMIPower IPMI 电源控制
// @Tags      Server
// @Router    /server/ipmiPower [post]
func (s *ServerApi) IPMIPower(c *gin.Context) {
	var req ipmiPowerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	out, err := serverService.IPMIPower(req.ID, req.Action)
	if err != nil {
		response.FailWithMessage("IPMI 操作失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]interface{}{"output": out}, "操作成功", c)
}

// GetServerMetrics 服务器性能历史
// @Tags      Server
// @Router    /server/getServerMetrics [get]
func (s *ServerApi) GetServerMetrics(c *gin.Context) {
	var host biz.Server
	_ = c.ShouldBindQuery(&host)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "60"))
	list, err := metricService.ListByServer(host.ID, limit)
	if err != nil {
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetServerMetricLatest 服务器最新性能指标
// @Tags      Server
// @Router    /server/getServerMetricLatest [get]
func (s *ServerApi) GetServerMetricLatest(c *gin.Context) {
	var host biz.Server
	_ = c.ShouldBindQuery(&host)
	data, err := metricService.LatestByServer(host.ID)
	if err != nil {
		response.FailWithMessage("暂无采集数据", c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// execCmdReq 批量执行命令入参
type execCmdReq struct {
	IDs []uint `json:"ids"`
	Cmd string `json:"cmd"`
}

// ExecCmd 批量在多台服务器执行命令（SSH）
// @Tags      Server
// @Router    /server/execCmd [post]
func (s *ServerApi) ExecCmd(c *gin.Context) {
	var req execCmdReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(req.IDs) == 0 || req.Cmd == "" {
		response.FailWithMessage("请选择服务器并输入命令", c)
		return
	}
	results := serverService.ExecCmd(req.IDs, req.Cmd)
	response.OkWithDetailed(results, "执行完成", c)
}

// ProbePort 端口探活（TCP 连通性 + 延迟）
// @Tags      Server
// @Router    /server/probePort [get]
func (s *ServerApi) ProbePort(c *gin.Context) {
	host := c.Query("host")
	port, _ := strconv.Atoi(c.Query("port"))
	if host == "" || port <= 0 {
		response.FailWithMessage("请提供 host 与 port", c)
		return
	}
	ok, ms, err := serverService.ProbePort(host, port)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	response.OkWithDetailed(map[string]interface{}{"alive": ok, "rtt": ms, "error": errMsg}, "检测完成", c)
}

// GetOverview 运维概览统计
// @Tags      Server
// @Router    /server/overview [get]
func (s *ServerApi) GetOverview(c *gin.Context) {
	response.OkWithDetailed(serverService.Overview(), "获取成功", c)
}
