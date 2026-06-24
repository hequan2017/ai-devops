package biz

import (
	"ai-devops/server/global"
)

// K8sCluster Kubernetes 集群接入点（参考 Kuboard / Dashboard）
type K8sCluster struct {
	global.GVA_MODEL
	Name        string `json:"name" form:"name" gorm:"comment:集群名称"`           // 集群名称
	APIServer   string `json:"apiServer" form:"apiServer" gorm:"comment:API Server地址"` // API Server 地址，如 https://1.2.3.4:6443
	Token       string `json:"token" form:"token" gorm:"comment:ServiceAccount Token"`  // Bearer Token
	InsecureTLS bool   `json:"insecureTls" form:"insecureTls" gorm:"comment:是否跳过TLS校验"` // 自签证书时跳过校验
	Status      string `json:"status" form:"status" gorm:"comment:连接状态"`          // online / offline
	Version     string `json:"version" form:"version" gorm:"comment:集群版本"`         // 集群 gitVersion
	Remark      string `json:"remark" form:"remark" gorm:"comment:备注"`           // 备注
}

func (K8sCluster) TableName() string {
	return "biz_k8s_cluster"
}
