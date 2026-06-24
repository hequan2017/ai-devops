package biz

import (
	"ai-devops/server/global"
)

// DockerHost Docker 引擎接入点（参考 Portainer endpoint）
type DockerHost struct {
	global.GVA_MODEL
	Name      string `json:"name" form:"name" gorm:"comment:接入点名称"`     // 接入点名称
	Protocol  string `json:"protocol" form:"protocol" gorm:"comment:协议"`   // 协议: unix / tcp
	Host      string `json:"host" form:"host" gorm:"comment:主机地址或socket路径"` // tcp 时为 IP，unix 时为 socket 路径
	Port      string `json:"port" form:"port" gorm:"comment:端口"`          // tcp 时使用的端口
	TLSVerify bool   `json:"tlsVerify" form:"tlsVerify" gorm:"comment:是否启用TLS"` // 是否启用 TLS（预留）
	Status    string `json:"status" form:"status" gorm:"comment:连接状态"`    // online / offline
	Version   string `json:"version" form:"version" gorm:"comment:Docker版本"` // Docker 版本
	Remark    string `json:"remark" form:"remark" gorm:"comment:备注"`      // 备注
}

func (DockerHost) TableName() string {
	return "biz_docker_host"
}

// Dsn 拼接 Docker daemon 连接地址
// unix: unix:///var/run/docker.sock ; tcp: tcp://192.168.1.10:2375
func (d DockerHost) Dsn() string {
	switch d.Protocol {
	case "tcp":
		port := d.Port
		if port == "" {
			port = "2375"
		}
		return "tcp://" + d.Host + ":" + port
	default:
		return "unix://" + d.Host
	}
}
