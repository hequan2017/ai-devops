package router

import (
	"ai-devops/server/router/biz"
	"ai-devops/server/router/example"
	"ai-devops/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Biz     biz.RouterGroup
}
