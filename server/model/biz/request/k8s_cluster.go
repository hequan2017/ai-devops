package request

import (
	"ai-devops/server/model/common/request"
)

// K8sClusterSearch K8s 集群分页搜索条件
type K8sClusterSearch struct {
	request.PageInfo
	Status string `json:"status" form:"status"`
}
