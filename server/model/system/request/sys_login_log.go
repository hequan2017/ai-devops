package request

import (
	"ai-devops/server/model/common/request"
	"ai-devops/server/model/system"
)

type SysLoginLogSearch struct {
	system.SysLoginLog
	request.PageInfo
}
