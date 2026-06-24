package biz

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
	bizReq "ai-devops/server/model/biz/request"
	"ai-devops/server/model/common/request"
)

type SshKeyService struct{}

var SshKeyServiceApp = new(SshKeyService)

func (s *SshKeyService) Create(k biz.SshKey) (err error) { return global.GVA_DB.Create(&k).Error }
func (s *SshKeyService) Delete(k biz.SshKey) (err error) { return global.GVA_DB.Delete(&k).Error }
func (s *SshKeyService) DeleteByIds(ids request.IdsReq) (err error) {
	return global.GVA_DB.Delete(&[]biz.SshKey{}).Where("id in ?", ids.Ids).Error
}
func (s *SshKeyService) Update(k biz.SshKey) (err error) { return global.GVA_DB.Save(&k).Error }
func (s *SshKeyService) Get(id uint) (k biz.SshKey, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&k).Error
	return
}
func (s *SshKeyService) GetList(info bizReq.SshKeySearch) (list []biz.SshKey, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&biz.SshKey{})
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

// GetAll 全部密钥（供服务器下拉选择）
func (s *SshKeyService) GetAll() (list []biz.SshKey, err error) {
	err = global.GVA_DB.Order("id desc").Find(&list).Error
	return
}
