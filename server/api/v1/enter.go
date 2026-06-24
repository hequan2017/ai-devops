package v1

import (
	"ai-devops/server/api/v1/biz"
	"ai-devops/server/api/v1/example"
	"ai-devops/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	BizApiGroup     biz.ApiGroup
}
