package request

import (
	"ai-devops/server/model/common/request"
)

// DockerHostSearch Docker 接入点分页搜索条件
type DockerHostSearch struct {
	request.PageInfo
	Protocol string `json:"protocol" form:"protocol"`
	Status   string `json:"status" form:"status"`
}
