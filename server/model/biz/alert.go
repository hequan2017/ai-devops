package biz

import (
	"ai-devops/server/global"
)

// AlertRule 告警规则（与性能采集联动）
type AlertRule struct {
	global.GVA_MODEL
	Name      string  `json:"name" form:"name" gorm:"comment:规则名称"`
	Metric    string  `json:"metric" form:"metric" gorm:"comment:指标(cpu/mem/disk)"` // cpu/mem/disk 百分比
	Operator  string  `json:"operator" form:"operator" gorm:"comment:操作符(>,>=,<,<=,==)"`
	Threshold float64 `json:"threshold" form:"threshold" gorm:"comment:阈值"`
	Level     string  `json:"level" form:"level" gorm:"comment:级别(warn/critical)"`
	Enabled   bool    `json:"enabled" form:"enabled" gorm:"comment:是否启用"`
	Remark    string  `json:"remark" form:"remark" gorm:"comment:备注"`
}

func (AlertRule) TableName() string {
	return "biz_alert_rule"
}

// AlertRecord 告警记录
type AlertRecord struct {
	global.GVA_MODEL
	RuleID     uint    `json:"ruleId" gorm:"index;comment:规则ID"`
	RuleName   string  `json:"ruleName" gorm:"comment:规则名"`
	ServerID   uint    `json:"serverId" gorm:"index;comment:服务器ID"`
	ServerName string  `json:"serverName" gorm:"comment:服务器名"`
	Metric     string  `json:"metric" gorm:"comment:指标"`
	Value      float64 `json:"value" gorm:"comment:实际值"`
	Level      string  `json:"level" gorm:"comment:级别"`
	Message    string  `json:"message" gorm:"comment:告警信息"`
	Resolved   bool    `json:"resolved" gorm:"comment:是否已处理"`
}

func (AlertRecord) TableName() string {
	return "biz_alert_record"
}
