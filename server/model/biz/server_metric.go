package biz

import (
	"ai-devops/server/global"
)

// ServerMetric 服务器性能采集记录
type ServerMetric struct {
	global.GVA_MODEL
	ServerID  uint    `json:"serverId" gorm:"index;comment:服务器ID"`
	CpuUsage  float64 `json:"cpuUsage" gorm:"comment:CPU使用率%"`
	MemTotal  uint64  `json:"memTotal" gorm:"comment:内存总量(B)"`
	MemUsed   uint64  `json:"memUsed" gorm:"comment:内存已用(B)"`
	DiskTotal uint64  `json:"diskTotal" gorm:"comment:磁盘总量(B)"`
	DiskUsed  uint64  `json:"diskUsed" gorm:"comment:磁盘已用(B)"`
	NetRx     uint64  `json:"netRx" gorm:"comment:网络接收累计(B)"`
	NetTx     uint64  `json:"netTx" gorm:"comment:网络发送累计(B)"`
}

func (ServerMetric) TableName() string {
	return "biz_server_metric"
}
