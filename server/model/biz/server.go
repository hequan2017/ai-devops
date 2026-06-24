package biz

import (
	"ai-devops/server/global"
)

// Server 服务器资产信息（支持戴尔 / 华为等物理服务器）
type Server struct {
	global.GVA_MODEL
	Name         string `json:"name" form:"name" gorm:"comment:服务器名称"`               // 服务器名称 / 主机名
	Manufacturer string `json:"manufacturer" form:"manufacturer" gorm:"comment:厂商"`     // 厂商: dell / huawei
	Model        string `json:"model" form:"model" gorm:"comment:型号"`                 // 型号，如 PowerEdge R740 / 2288H V5
	SerialNumber string `json:"serialNumber" form:"serialNumber" gorm:"comment:序列号SN"` // 序列号
	ManageIP     string `json:"manageIp" form:"manageIp" gorm:"comment:带外管理IP"`        // 带外管理 IP（iDRAC / iBMC）
	HostIP       string `json:"hostIp" form:"hostIp" gorm:"comment:业务IP"`              // 业务 IP
	DataCenter   string `json:"dataCenter" form:"dataCenter" gorm:"comment:机房"`         // 所在机房
	Cabinet      string `json:"cabinet" form:"cabinet" gorm:"comment:机柜"`              // 机柜编号
	UPosition    string `json:"uPosition" form:"uPosition" gorm:"comment:U位"`           // U 位
	CPU          string `json:"cpu" form:"cpu" gorm:"comment:CPU配置"`                   // CPU
	Memory       string `json:"memory" form:"memory" gorm:"comment:内存"`                 // 内存
	Disk         string `json:"disk" form:"disk" gorm:"comment:磁盘"`                   // 磁盘
	OS           string `json:"os" form:"os" gorm:"comment:操作系统"`                   // 操作系统
	Status       string `json:"status" form:"status" gorm:"comment:状态"`                 // 状态: online/offline/maintenance/fault
	Owner        string `json:"owner" form:"owner" gorm:"comment:负责人"`                 // 负责人
	Remark       string `json:"remark" form:"remark" gorm:"comment:备注"`                 // 备注
	IpmiIP       string `json:"ipmiIp" form:"ipmiIp" gorm:"comment:IPMI地址"`               // IPMI / BMC 地址
	IpmiUser     string `json:"ipmiUser" form:"ipmiUser" gorm:"comment:IPMI用户名"`           // IPMI 用户名
	IpmiPassword string `json:"ipmiPassword" form:"ipmiPassword" gorm:"comment:IPMI密码"`       // IPMI 密码
	SshPort      int    `json:"sshPort" form:"sshPort" gorm:"comment:SSH端口"`               // SSH 端口
	SshUser      string `json:"sshUser" form:"sshUser" gorm:"comment:SSH用户名"`             // SSH 用户名
	SshAuthType  string `json:"sshAuthType" form:"sshAuthType" gorm:"comment:SSH认证方式"`       // password / key
	SshPassword  string `json:"sshPassword" form="sshPassword" gorm:"comment:SSH密码"`         // SSH 密码
	SshKey       string `json:"sshKey" form:"sshKey" gorm:"comment:SSH私钥"`               // SSH 私钥(PEM)
}

// TableName 业务表统一使用 biz_ 前缀
func (Server) TableName() string {
	return "biz_server"
}
