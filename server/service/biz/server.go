package biz

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/biz"
	bizReq "github.com/flipped-aurora/gin-vue-admin/server/model/biz/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type ServerService struct{}

var ServerServiceApp = new(ServerService)

// CreateServer 新增服务器
func (s *ServerService) CreateServer(server biz.Server) (err error) {
	err = global.GVA_DB.Create(&server).Error
	return err
}

// DeleteServer 删除服务器
func (s *ServerService) DeleteServer(server biz.Server) (err error) {
	err = global.GVA_DB.Delete(&server).Error
	return err
}

// DeleteServerByIds 批量删除服务器
func (s *ServerService) DeleteServerByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]biz.Server{}).Where("id in ?", ids.Ids).Error
	return err
}

// UpdateServer 更新服务器
func (s *ServerService) UpdateServer(server biz.Server) (err error) {
	err = global.GVA_DB.Save(&server).Error
	return err
}

// GetServer 根据 ID 获取服务器详情
func (s *ServerService) GetServer(id uint) (server biz.Server, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&server).Error
	return
}

// GetServerList 分页获取服务器列表，支持厂商 / 状态过滤与关键字模糊匹配
func (s *ServerService) GetServerList(info bizReq.ServerSearch) (list []biz.Server, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&biz.Server{})

	if info.Manufacturer != "" {
		db = db.Where("manufacturer = ?", info.Manufacturer)
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.Keyword != "" {
		like := "%" + info.Keyword + "%"
		db = db.Where("name LIKE ? OR serial_number LIKE ? OR host_ip LIKE ? OR manage_ip LIKE ?", like, like, like, like)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("id desc").Find(&list).Error
	return list, total, err
}
