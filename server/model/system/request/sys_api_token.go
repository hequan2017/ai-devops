package request

import (
	"ai-devops/server/model/common/request"
	"ai-devops/server/model/system"
)

type SysApiTokenSearch struct {
	system.SysApiToken
	request.PageInfo
    Status *bool `json:"status" form:"status"`
}
