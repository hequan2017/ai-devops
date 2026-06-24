package biz

import (
	"ai-devops/server/global"
)

// SshKey SSH 密钥（集中管理，服务器可引用）
type SshKey struct {
	global.GVA_MODEL
	Name       string `json:"name" form:"name" gorm:"comment:名称"`
	PrivateKey string `json:"privateKey" form:"privateKey" gorm:"comment:私钥(PEM)"`
	Remark     string `json:"remark" form:"remark" gorm:"comment:备注"`
}

func (SshKey) TableName() string {
	return "biz_ssh_key"
}
