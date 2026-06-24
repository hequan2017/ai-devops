package biz

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
	"sync"
	"time"
)

// ExecResult 单台服务器命令执行结果
type ExecResult struct {
	ServerID   uint   `json:"serverId"`
	ServerName string `json:"serverName"`
	Output     string `json:"output"`
	Error      string `json:"error"`
}

// ExecCmd 在多台服务器上并发执行 shell 命令（SSH）
// 并发上限 10，单命令超时 30s，避免慢命令拖垮整体
func (s *ServerService) ExecCmd(ids []uint, cmd string) []ExecResult {
	var servers []biz.Server
	global.GVA_DB.Where("id in ?", ids).Find(&servers)
	results := make([]ExecResult, len(servers))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)
	for i, sv := range servers {
		wg.Add(1)
		sem <- struct{}{}
		go func(idx int, server biz.Server) {
			defer wg.Done()
			defer func() { <-sem }()
			results[idx] = s.execOne(server, cmd)
		}(i, sv)
	}
	wg.Wait()
	return results
}

// execOne 单台服务器执行命令（含 30s 超时）
func (s *ServerService) execOne(server biz.Server, cmd string) ExecResult {
	r := ExecResult{ServerID: server.ID, ServerName: server.Name}
	if server.SshUser == "" || server.HostIP == "" {
		r.Error = "未配置SSH凭证或业务IP"
		return r
	}
	client, err := sshDial(server)
	if err != nil {
		r.Error = err.Error()
		return r
	}
	defer client.Close()
	session, err := client.NewSession()
	if err != nil {
		r.Error = err.Error()
		return r
	}
	defer session.Close()

	// 命令超时控制（ssh.Session 无 SetTimeout，用 goroutine + timer）
	done := make(chan struct{})
	var out []byte
	var execErr error
	go func() {
		out, execErr = session.CombinedOutput(cmd)
		close(done)
	}()
	select {
	case <-done:
		r.Output = string(out)
		if execErr != nil {
			r.Error = execErr.Error()
		}
	case <-time.After(30 * time.Second):
		_ = session.Close() // 中断阻塞的命令
		<-done
		r.Output = string(out)
		r.Error = "命令执行超时(30s)"
	}
	return r
}
