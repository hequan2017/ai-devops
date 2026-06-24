package biz

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
)

// ExecResult 单台服务器命令执行结果
type ExecResult struct {
	ServerID   uint   `json:"serverId"`
	ServerName string `json:"serverName"`
	Output     string `json:"output"`
	Error      string `json:"error"`
}

// ExecCmd 在多台服务器上批量执行 shell 命令（通过 SSH）
func (s *ServerService) ExecCmd(ids []uint, cmd string) []ExecResult {
	var servers []biz.Server
	global.GVA_DB.Where("id in ?", ids).Find(&servers)
	results := make([]ExecResult, 0, len(servers))
	for _, sv := range servers {
		r := ExecResult{ServerID: sv.ID, ServerName: sv.Name}
		if sv.SshUser == "" || sv.HostIP == "" {
			r.Error = "未配置SSH凭证或业务IP"
			results = append(results, r)
			continue
		}
		client, err := sshDial(sv)
		if err != nil {
			r.Error = err.Error()
			results = append(results, r)
			continue
		}
		session, err := client.NewSession()
		if err != nil {
			r.Error = err.Error()
			client.Close()
			results = append(results, r)
			continue
		}
		out, err := session.CombinedOutput(cmd)
		r.Output = string(out)
		if err != nil {
			r.Error = err.Error()
		}
		session.Close()
		client.Close()
		results = append(results, r)
	}
	return results
}
