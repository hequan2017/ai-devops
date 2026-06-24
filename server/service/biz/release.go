package biz

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
	bizReq "ai-devops/server/model/biz/request"
	"ai-devops/server/model/common/request"
	"fmt"
	"os/exec"
	"time"
)

type ReleaseService struct{}

var ReleaseServiceApp = new(ReleaseService)

// CreateRelease 新建发版（草稿）
func (s *ReleaseService) CreateRelease(r biz.Release) (err error) {
	if r.Status == "" {
		r.Status = "draft"
	}
	return global.GVA_DB.Create(&r).Error
}

// DeleteRelease 删除
func (s *ReleaseService) DeleteRelease(r biz.Release) (err error) {
	return global.GVA_DB.Delete(&r).Error
}

// DeleteReleaseByIds 批量删除
func (s *ReleaseService) DeleteReleaseByIds(ids request.IdsReq) (err error) {
	return global.GVA_DB.Delete(&[]biz.Release{}).Where("id in ?", ids.Ids).Error
}

// UpdateRelease 更新
func (s *ReleaseService) UpdateRelease(r biz.Release) (err error) {
	return global.GVA_DB.Save(&r).Error
}

// GetRelease 详情
func (s *ReleaseService) GetRelease(id uint) (r biz.Release, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&r).Error
	return
}

// GetReleaseList 分页列表
func (s *ReleaseService) GetReleaseList(info bizReq.ReleaseSearch) (list []biz.Release, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&biz.Release{})
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.Environment != "" {
		db = db.Where("environment = ?", info.Environment)
	}
	if info.Project != "" {
		db = db.Where("project LIKE ?", "%"+info.Project+"%")
	}
	if info.Keyword != "" {
		db = db.Where("project LIKE ? OR version LIKE ?", "%"+info.Keyword+"%", "%"+info.Keyword+"%")
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

// Submit 提交审批 draft -> pending
func (s *ReleaseService) Submit(id uint) (err error) {
	return global.GVA_DB.Model(&biz.Release{}).Where("id = ?", id).Update("status", "pending").Error
}

// Approve 审批通过 pending -> approved
func (s *ReleaseService) Approve(id, approverID uint) (err error) {
	now := time.Now()
	return global.GVA_DB.Model(&biz.Release{}).Where("id = ? AND status = ?", id, "pending").
		Updates(map[string]interface{}{"status": "approved", "approver_id": approverID, "approve_time": now}).Error
}

// Reject 审批拒绝 pending -> rejected
func (s *ReleaseService) Reject(id, approverID uint, reason string) (err error) {
	now := time.Now()
	return global.GVA_DB.Model(&biz.Release{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": "rejected", "approver_id": approverID, "approve_time": now, "result": reason}).Error
}

// Execute 执行发布脚本 approved -> executing -> done/failed
func (s *ReleaseService) Execute(id uint) (result string, err error) {
	var r biz.Release
	if err = global.GVA_DB.Where("id = ?", id).First(&r).Error; err != nil {
		return
	}
	if r.Status != "approved" {
		err = fmt.Errorf("仅审批通过的发版可执行, 当前状态: %s", r.Status)
		return
	}
	global.GVA_DB.Model(&r).Update("status", "executing")
	now := time.Now()
	cmd := exec.Command("sh", "-c", r.Script)
	out, execErr := cmd.CombinedOutput()
	result = string(out)
	status := "done"
	if execErr != nil {
		status = "failed"
		result += "\n" + execErr.Error()
	}
	global.GVA_DB.Model(&r).Updates(map[string]interface{}{"status": status, "result": result, "execute_time": now})
	if execErr != nil {
		err = fmt.Errorf("执行失败: %s", result)
	}
	return
}
