package biz

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
	bizReq "ai-devops/server/model/biz/request"
	"ai-devops/server/model/common/request"
)

type DockerHostService struct{}

var DockerHostServiceApp = new(DockerHostService)

// CreateDockerHost 新增 Docker 接入点
func (s *DockerHostService) CreateDockerHost(d biz.DockerHost) (err error) {
	return global.GVA_DB.Create(&d).Error
}

// DeleteDockerHost 删除 Docker 接入点
func (s *DockerHostService) DeleteDockerHost(d biz.DockerHost) (err error) {
	return global.GVA_DB.Delete(&d).Error
}

// DeleteDockerHostByIds 批量删除
func (s *DockerHostService) DeleteDockerHostByIds(ids request.IdsReq) (err error) {
	return global.GVA_DB.Delete(&[]biz.DockerHost{}).Where("id in ?", ids.Ids).Error
}

// UpdateDockerHost 更新 Docker 接入点
func (s *DockerHostService) UpdateDockerHost(d biz.DockerHost) (err error) {
	return global.GVA_DB.Save(&d).Error
}

// GetDockerHost 根据 ID 获取接入点
func (s *DockerHostService) GetDockerHost(id uint) (d biz.DockerHost, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&d).Error
	return
}

// GetDockerHostList 分页获取接入点列表
func (s *DockerHostService) GetDockerHostList(info bizReq.DockerHostSearch) (list []biz.DockerHost, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&biz.DockerHost{})
	if info.Protocol != "" {
		db = db.Where("protocol = ?", info.Protocol)
	}
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

// TestDockerHost 测试接入点连通性，返回版本并更新状态
func (s *DockerHostService) TestDockerHost(id uint) (version string, err error) {
	d, err := s.GetDockerHost(id)
	if err != nil {
		return
	}
	version, err = newDockerClient(d.Dsn()).Ping()
	if err != nil {
		global.GVA_DB.Model(&d).Updates(map[string]interface{}{"status": "offline"})
		return
	}
	global.GVA_DB.Model(&d).Updates(map[string]interface{}{"status": "online", "version": version})
	return
}

// ListContainers 容器列表
func (s *DockerHostService) ListContainers(id uint) ([]DockerContainer, error) {
	d, err := s.GetDockerHost(id)
	if err != nil {
		return nil, err
	}
	return newDockerClient(d.Dsn()).ListContainers()
}

// ContainerLogs 容器日志
func (s *DockerHostService) ContainerLogs(id uint, containerID string) (string, error) {
	d, err := s.GetDockerHost(id)
	if err != nil {
		return "", err
	}
	return newDockerClient(d.Dsn()).ContainerLogs(containerID)
}

// ContainerAction 容器操作 start/stop/restart/remove
func (s *DockerHostService) ContainerAction(id uint, containerID, action string) error {
	d, err := s.GetDockerHost(id)
	if err != nil {
		return err
	}
	return newDockerClient(d.Dsn()).ContainerAction(containerID, action)
}

// ListImages 镜像列表
func (s *DockerHostService) ListImages(id uint) ([]DockerImage, error) {
	d, err := s.GetDockerHost(id)
	if err != nil {
		return nil, err
	}
	return newDockerClient(d.Dsn()).ListImages()
}

// ListNetworks 网络列表
func (s *DockerHostService) ListNetworks(id uint) ([]DockerNetwork, error) {
	d, err := s.GetDockerHost(id)
	if err != nil {
		return nil, err
	}
	return newDockerClient(d.Dsn()).ListNetworks()
}

// ListVolumes 数据卷列表
func (s *DockerHostService) ListVolumes(id uint) ([]DockerVolume, error) {
	d, err := s.GetDockerHost(id)
	if err != nil {
		return nil, err
	}
	return newDockerClient(d.Dsn()).ListVolumes()
}
