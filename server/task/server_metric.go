package task

import (
	"ai-devops/server/service/biz"
)

// CollectServerMetric 采集所有服务器性能指标（供定时任务调用）
func CollectServerMetric() {
	biz.ServerMetricServiceApp.CollectAll()
}
