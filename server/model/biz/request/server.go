package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// ServerSearch 服务器分页搜索条件
type ServerSearch struct {
	request.PageInfo
	Manufacturer string `json:"manufacturer" form:"manufacturer"` // 厂商过滤: dell / huawei
	Status       string `json:"status" form:"status"`             // 状态过滤
}
