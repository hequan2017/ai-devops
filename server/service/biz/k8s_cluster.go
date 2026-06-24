package biz

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
	bizReq "ai-devops/server/model/biz/request"
	"ai-devops/server/model/common/request"
)

type K8sClusterService struct{}

var K8sClusterServiceApp = new(K8sClusterService)

func (s *K8sClusterService) cli(c biz.K8sCluster) *k8sClient {
	return newK8sClient(c.APIServer, c.Token, c.InsecureTLS)
}

// CreateK8sCluster 新增集群接入点
func (s *K8sClusterService) CreateK8sCluster(c biz.K8sCluster) (err error) {
	return global.GVA_DB.Create(&c).Error
}

// DeleteK8sCluster 删除集群接入点
func (s *K8sClusterService) DeleteK8sCluster(c biz.K8sCluster) (err error) {
	return global.GVA_DB.Delete(&c).Error
}

// DeleteK8sClusterByIds 批量删除
func (s *K8sClusterService) DeleteK8sClusterByIds(ids request.IdsReq) (err error) {
	return global.GVA_DB.Delete(&[]biz.K8sCluster{}).Where("id in ?", ids.Ids).Error
}

// UpdateK8sCluster 更新集群接入点
func (s *K8sClusterService) UpdateK8sCluster(c biz.K8sCluster) (err error) {
	return global.GVA_DB.Save(&c).Error
}

// GetK8sCluster 根据 ID 获取集群
func (s *K8sClusterService) GetK8sCluster(id uint) (c biz.K8sCluster, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&c).Error
	return
}

// GetK8sClusterList 分页列表
func (s *K8sClusterService) GetK8sClusterList(info bizReq.K8sClusterSearch) (list []biz.K8sCluster, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&biz.K8sCluster{})
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.Keyword != "" {
		db = db.Where("name LIKE ?", "%"+info.Keyword+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("id desc").Find(&list).Error
	return
}

// TestCluster 测试连通性并更新状态
func (s *K8sClusterService) TestCluster(id uint) (version string, err error) {
	c, err := s.GetK8sCluster(id)
	if err != nil {
		return
	}
	version, err = s.cli(c).Ping()
	if err != nil {
		global.GVA_DB.Model(&c).Updates(map[string]interface{}{"status": "offline"})
		return
	}
	global.GVA_DB.Model(&c).Updates(map[string]interface{}{"status": "online", "version": version})
	return
}

// ListNamespaces 命名空间列表
func (s *K8sClusterService) ListNamespaces(id uint) ([]map[string]interface{}, error) {
	c, err := s.GetK8sCluster(id)
	if err != nil {
		return nil, err
	}
	return s.cli(c).Namespaces()
}

// ListPods 指定命名空间的 Pod 列表
func (s *K8sClusterService) ListPods(id uint, ns string) ([]map[string]interface{}, error) {
	c, err := s.GetK8sCluster(id)
	if err != nil {
		return nil, err
	}
	return s.cli(c).Pods(ns)
}

// ListNodes 节点列表
func (s *K8sClusterService) ListNodes(id uint) ([]map[string]interface{}, error) {
	c, err := s.GetK8sCluster(id)
	if err != nil {
		return nil, err
	}
	return s.cli(c).Nodes()
}

// ListDeployments 指定命名空间的 Deployment 列表
func (s *K8sClusterService) ListDeployments(id uint, ns string) ([]map[string]interface{}, error) {
	c, err := s.GetK8sCluster(id)
	if err != nil {
		return nil, err
	}
	return s.cli(c).Deployments(ns)
}

// ListServices 指定命名空间的 Service 列表
func (s *K8sClusterService) ListServices(id uint, ns string) ([]map[string]interface{}, error) {
	c, err := s.GetK8sCluster(id)
	if err != nil {
		return nil, err
	}
	return s.cli(c).Services(ns)
}

// PodLogs Pod 日志
func (s *K8sClusterService) PodLogs(id uint, ns, pod string) (string, error) {
	c, err := s.GetK8sCluster(id)
	if err != nil {
		return "", err
	}
	return s.cli(c).PodLogs(ns, pod)
}

// DeletePod 删除 Pod
func (s *K8sClusterService) DeletePod(id uint, ns, pod string) error {
	c, err := s.GetK8sCluster(id)
	if err != nil {
		return err
	}
	return s.cli(c).DeletePod(ns, pod)
}
