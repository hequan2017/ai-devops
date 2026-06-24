package biz

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type ServerMetricService struct{}

var ServerMetricServiceApp = new(ServerMetricService)

// CollectAll 采集所有配置了 SSH 凭证的服务器（并发）
func (s *ServerMetricService) CollectAll() {
	var servers []biz.Server
	if err := global.GVA_DB.Find(&servers).Error; err != nil {
		global.GVA_LOG.Error("采集: 查询服务器失败", zap.Error(err))
		return
	}
	for i := range servers {
		sv := servers[i]
		if sv.SshUser == "" || sv.HostIP == "" {
			continue
		}
		go func(server biz.Server) {
			if err := s.collectOne(server); err != nil {
				global.GVA_LOG.Warn("采集失败 "+server.Name, zap.Error(err))
			}
		}(sv)
	}
}

func (s *ServerMetricService) collectOne(server biz.Server) error {
	m, err := collectViaSSH(server)
	if err != nil {
		return err
	}
	m.ServerID = server.ID
	return global.GVA_DB.Create(m).Error
}

// ListByServer 某服务器最近 N 条采集记录
func (s *ServerMetricService) ListByServer(serverID uint, limit int) (list []biz.ServerMetric, err error) {
	if limit <= 0 || limit > 2000 {
		limit = 60
	}
	err = global.GVA_DB.Where("server_id = ?", serverID).Order("id desc").Limit(limit).Find(&list).Error
	return
}

// LatestByServer 某服务器最新一条采集记录
func (s *ServerMetricService) LatestByServer(serverID uint) (m biz.ServerMetric, err error) {
	err = global.GVA_DB.Where("server_id = ?", serverID).Order("id desc").First(&m).Error
	return
}

// 通过 SSH 执行命令采集 CPU/内存/磁盘/网络（兼容主流 Linux）
const metricCmd = `vmstat 1 2 | tail -1 | awk '{print "CPU="100-$15}'; ` +
	`free -b | awk '/Mem/{print "MEM="$2" "$3}'; ` +
	`df -B1 / | awk 'NR==2{print "DISK="$2" "$3}'; ` +
	`awk 'NR>2{r+=$2;t+=$10}END{print "NET="r" "t}' /proc/net/dev`

func collectViaSSH(s biz.Server) (*biz.ServerMetric, error) {
	client, err := sshDial(s)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	out, err := session.CombinedOutput(metricCmd)
	if err != nil {
		return nil, err
	}
	return parseMetric(string(out))
}

// parseMetric 解析 "CPU=x MEM=t u DISK=t u NET=rx tx" 形式输出
func parseMetric(out string) (*biz.ServerMetric, error) {
	m := &biz.ServerMetric{}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		idx := strings.Index(line, "=")
		if idx <= 0 {
			continue
		}
		key := line[:idx]
		fields := strings.Fields(line[idx+1:])
		switch key {
		case "CPU":
			if len(fields) > 0 {
				m.CpuUsage, _ = strconv.ParseFloat(fields[0], 64)
			}
		case "MEM":
			if len(fields) >= 2 {
				m.MemTotal, _ = strconv.ParseUint(fields[0], 10, 64)
				m.MemUsed, _ = strconv.ParseUint(fields[1], 10, 64)
			}
		case "DISK":
			if len(fields) >= 2 {
				m.DiskTotal, _ = strconv.ParseUint(fields[0], 10, 64)
				m.DiskUsed, _ = strconv.ParseUint(fields[1], 10, 64)
			}
		case "NET":
			if len(fields) >= 2 {
				m.NetRx, _ = strconv.ParseUint(fields[0], 10, 64)
				m.NetTx, _ = strconv.ParseUint(fields[1], 10, 64)
			}
		}
	}
	return m, nil
}
