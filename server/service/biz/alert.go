package biz

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
	bizReq "ai-devops/server/model/biz/request"
	"ai-devops/server/model/common/request"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type AlertService struct{}

var AlertServiceApp = new(AlertService)

// ===== 规则 CRUD =====

func (s *AlertService) CreateRule(r biz.AlertRule) (err error) { return global.GVA_DB.Create(&r).Error }
func (s *AlertService) DeleteRule(r biz.AlertRule) (err error) { return global.GVA_DB.Delete(&r).Error }
func (s *AlertService) DeleteRuleByIds(ids request.IdsReq) (err error) {
	return global.GVA_DB.Delete(&[]biz.AlertRule{}).Where("id in ?", ids.Ids).Error
}
func (s *AlertService) UpdateRule(r biz.AlertRule) (err error) { return global.GVA_DB.Save(&r).Error }
func (s *AlertService) GetRule(id uint) (r biz.AlertRule, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&r).Error
	return
}
func (s *AlertService) GetRuleList(info bizReq.AlertRuleSearch) (list []biz.AlertRule, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&biz.AlertRule{})
	if info.Metric != "" {
		db = db.Where("metric = ?", info.Metric)
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

// ===== 告警检测（采集后调用）=====

// Check 检查单台服务器指标是否触发告警，超阈值则记录（避免同规则重复告警）
func (s *AlertService) Check(server biz.Server, m *biz.ServerMetric) {
	var rules []biz.AlertRule
	if err := global.GVA_DB.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		return
	}
	for _, r := range rules {
		val := metricPercent(m, r.Metric)
		if !matchRule(val, r.Operator, r.Threshold) {
			continue
		}
		var exist biz.AlertRecord
		_ = global.GVA_DB.Where("server_id = ? AND rule_id = ? AND resolved = ?", server.ID, r.ID, false).First(&exist).Error
		if exist.ID > 0 {
			continue // 已存在未处理告警，不重复
		}
		rec := biz.AlertRecord{
			RuleID:     r.ID,
			RuleName:   r.Name,
			ServerID:   server.ID,
			ServerName: server.Name,
			Metric:     r.Metric,
			Value:      val,
			Level:      r.Level,
			Message:    fmt.Sprintf("%s %.2f%% %s %.2f%% (%s)", r.Metric, val, r.Operator, r.Threshold, server.Name),
		}
		global.GVA_DB.Create(&rec)
		if r.NotifyWebhook != "" {
			go s.postWebhook(r.NotifyWebhook, rec, r)
		}
	}
}

func metricPercent(m *biz.ServerMetric, metric string) float64 {
	switch metric {
	case "cpu":
		return m.CpuUsage
	case "mem":
		if m.MemTotal > 0 {
			return float64(m.MemUsed) / float64(m.MemTotal) * 100
		}
	case "disk":
		if m.DiskTotal > 0 {
			return float64(m.DiskUsed) / float64(m.DiskTotal) * 100
		}
	}
	return 0
}

func matchRule(val float64, op string, threshold float64) bool {
	switch op {
	case ">":
		return val > threshold
	case ">=":
		return val >= threshold
	case "<":
		return val < threshold
	case "<=":
		return val <= threshold
	case "==":
		return val == threshold
	}
	return false
}

// ===== 记录查询 =====

func (s *AlertService) GetRecordList(info bizReq.AlertRecordSearch) (list []biz.AlertRecord, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&biz.AlertRecord{})
	if info.ServerID > 0 {
		db = db.Where("server_id = ?", info.ServerID)
	}
	if info.Level != "" {
		db = db.Where("level = ?", info.Level)
	}
	if info.Resolved != nil {
		db = db.Where("resolved = ?", *info.Resolved)
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

// ResolveRecord 标记告警已处理
func (s *AlertService) ResolveRecord(id uint) (err error) {
	return global.GVA_DB.Model(&biz.AlertRecord{}).Where("id = ?", id).Update("resolved", true).Error
}

// postWebhook 发送 webhook 通知（钉钉/企业微信/飞书通用 text 格式）
func (s *AlertService) postWebhook(url string, rec biz.AlertRecord, rule biz.AlertRule) {
	content := fmt.Sprintf("[AI运维告警][%s] %s\n服务器: %s\n规则: %s\n时间: %s",
		strings.ToUpper(rec.Level), rec.Message, rec.ServerName, rule.Name, time.Now().Format("2006-01-02 15:04:05"))
	payload := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": content},
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		global.GVA_LOG.Warn("webhook通知失败 "+rule.Name, zap.Error(err))
		return
	}
	defer resp.Body.Close()
}
