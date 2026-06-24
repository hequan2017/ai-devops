package service

import (
	"ai-devops/server/service/biz"
	"ai-devops/server/service/example"
	"ai-devops/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	BizServiceGroup     biz.ServiceGroup
}
