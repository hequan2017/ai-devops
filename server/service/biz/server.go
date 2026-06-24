package biz

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
	bizReq "ai-devops/server/model/biz/request"
	"ai-devops/server/model/common/request"
	"fmt"
	"net"
	"strconv"
	"time"
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

// IPMIPower 通过 IPMI 控制服务器电源（开关机/重启/状态查询）
func (s *ServerService) IPMIPower(id uint, action string) (output string, err error) {
	var server biz.Server
	if err = global.GVA_DB.Where("id = ?", id).First(&server).Error; err != nil {
		return
	}
	if server.IpmiIP == "" {
		err = fmt.Errorf("该服务器未配置 IPMI 地址")
		return
	}
	return ipmiPower(server.IpmiIP, server.IpmiUser, server.IpmiPassword, action)
}

// ProbePort TCP 端口探活，返回是否可达与延迟(ms)
func (s *ServerService) ProbePort(host string, port int) (alive bool, rtt int64, err error) {
	start := time.Now()
	conn, derr := net.DialTimeout("tcp", net.JoinHostPort(host, strconv.Itoa(port)), 3*time.Second)
	rtt = time.Since(start).Milliseconds()
	if derr != nil {
		return false, rtt, derr
	}
	conn.Close()
	return true, rtt, nil
}
