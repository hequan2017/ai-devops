package request

import (
	"ai-devops/server/model/common/request"
)

// AlertRuleSearch 告警规则搜索
type AlertRuleSearch struct {
	request.PageInfo
	Metric string `json:"metric" form:"metric"`
}

// AlertRecordSearch 告警记录搜索
type AlertRecordSearch struct {
	request.PageInfo
	ServerID uint   `json:"serverId" form:"serverId"`
	Level    string `json:"level" form:"level"`
	Resolved *bool  `json:"resolved" form:"resolved"`
}
