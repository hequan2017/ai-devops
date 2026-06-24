package request

import (
	"ai-devops/server/model/common/request"
	"ai-devops/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
