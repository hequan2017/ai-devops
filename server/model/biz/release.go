package biz

import (
	"ai-devops/server/global"
	"time"
)

// Release 发版工作流记录（类 Jenkins 流水线审批）
type Release struct {
	global.GVA_MODEL
	Project     string     `json:"project" form:"project" gorm:"comment:项目名"`           // 项目名
	Version     string     `json:"version" form:"version" gorm:"comment:版本号"`           // 版本号
	Branch      string     `json:"branch" form:"branch" gorm:"comment:代码分支"`           // 代码分支
	Environment string     `json:"environment" form:"environment" gorm:"comment:发布环境"`     // dev/staging/prod
	Description string     `json:"description" form:"description" gorm:"comment:发版说明"`     // 发版说明
	Status      string     `json:"status" form:"status" gorm:"comment:状态"`             // draft/pending/approved/rejected/executing/done/failed
	Script      string     `json:"script" form:"script" gorm:"comment:发布脚本"`           // 发布脚本/命令
	ApplicantID uint       `json:"applicantId" form:"applicantId" gorm:"comment:申请人"`     // 申请人
	ApproverID  uint       `json:"approverId" form:"approverId" gorm:"comment:审批人"`       // 审批人
	ApproveTime *time.Time `json:"approveTime" form:"approveTime" gorm:"comment:审批时间"`     // 审批时间
	ExecuteTime *time.Time `json:"executeTime" form:"executeTime" gorm:"comment:执行时间"`     // 执行时间
	Result      string     `json:"result" form:"result" gorm:"comment:执行结果"`           // 执行结果/日志
}

func (Release) TableName() string {
	return "biz_release"
}
