package biz

import "github.com/gin-gonic/gin"

type TerminalRouter struct{}

// InitTerminalRouter 注册 WebSSH WebSocket 路由
// WebSocket 走 public 组（浏览器无法设置请求头），在 handler 内通过 query token 校验
func (t *TerminalRouter) InitTerminalRouter(PublicRouter *gin.RouterGroup) {
	PublicRouter.Group("server").GET("terminal", terminalApi.ServerTerminal)
	PublicRouter.Group("docker").GET("containerLogsStream", terminalApi.DockerLogsStream)
	PublicRouter.Group("docker").GET("containerTerminal", terminalApi.DockerContainerTerminal)
}
