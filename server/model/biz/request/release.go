package request

import (
	"ai-devops/server/model/common/request"
)

// ReleaseSearch 发版分页搜索条件
type ReleaseSearch struct {
	request.PageInfo
	Status      string `json:"status" form:"status"`
	Environment string `json:"environment" form:"environment"`
	Project     string `json:"project" form:"project"`
}
