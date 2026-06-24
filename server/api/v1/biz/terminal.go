package biz

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
	"ai-devops/server/model/common/response"
	"ai-devops/server/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type TerminalApi struct{}

// wsAuth 校验 query 中的 token（浏览器 WebSocket 无法自定义请求头，故走 query）
func wsAuth(c *gin.Context) bool {
	token := c.Query("token")
	if token == "" {
		response.FailWithMessage("未授权", c)
		return false
	}
	if _, err := utils.NewJWT().ParseToken(token); err != nil {
		response.FailWithMessage("token无效", c)
		return false
	}
	return true
}

// ServerTerminal 服务器 SSH 终端 (WebSocket)
// @Tags      Terminal
// @Router    /server/terminal [get]
func (t *TerminalApi) ServerTerminal(c *gin.Context) {
	if !wsAuth(c) {
		return
	}
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	var server biz.Server
	if err := global.GVA_DB.Where("id = ?", id).First(&server).Error; err != nil {
		response.FailWithMessage("服务器不存在", c)
		return
	}
	if server.SshUser == "" || server.HostIP == "" {
		response.FailWithMessage("未配置SSH凭证或业务IP", c)
		return
	}
	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("ws升级失败", zap.Error(err))
		return
	}
	defer ws.Close()
	if err := terminalService.ServerTerminal(server, ws); err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("\r\n[连接错误] "+err.Error()+"\r\n"))
	}
}

// DockerLogsStream Docker 容器日志实时流 (WebSocket)
// @Tags      Terminal
// @Router    /docker/containerLogsStream [get]
func (t *TerminalApi) DockerLogsStream(c *gin.Context) {
	if !wsAuth(c) {
		return
	}
	hostID, _ := strconv.ParseUint(c.Query("hostId"), 10, 64)
	containerID := c.Query("containerId")
	var host biz.DockerHost
	if err := global.GVA_DB.Where("id = ?", hostID).First(&host).Error; err != nil {
		response.FailWithMessage("Docker接入点不存在", c)
		return
	}
	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("ws升级失败", zap.Error(err))
		return
	}
	defer ws.Close()
	_ = terminalService.DockerLogsStream(host, containerID, ws)
}

// DockerContainerTerminal Docker 容器交互式终端 (WebSocket)
// @Tags      Terminal
// @Router    /docker/containerTerminal [get]
func (t *TerminalApi) DockerContainerTerminal(c *gin.Context) {
	if !wsAuth(c) {
		return
	}
	hostID, _ := strconv.ParseUint(c.Query("hostId"), 10, 64)
	containerID := c.Query("containerId")
	var host biz.DockerHost
	if err := global.GVA_DB.Where("id = ?", hostID).First(&host).Error; err != nil {
		response.FailWithMessage("Docker接入点不存在", c)
		return
	}
	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.GVA_LOG.Error("ws升级失败", zap.Error(err))
		return
	}
	defer ws.Close()
	if err := terminalService.DockerExecTerminal(host, containerID, ws); err != nil {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("\r\n[连接错误] "+err.Error()+"\r\n"))
	}
}
