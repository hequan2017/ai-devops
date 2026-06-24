package biz

import (
	"ai-devops/server/model/biz"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
)

type TerminalService struct{}

var TerminalServiceApp = new(TerminalService)

// wsWriter 线程安全的 websocket 写入器
type wsWriter struct {
	ws *websocket.Conn
	mu sync.Mutex
}

func (w *wsWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err := w.ws.WriteMessage(websocket.TextMessage, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

// sshDial 建立 SSH 连接（支持密钥 / 密码）
func sshDial(s biz.Server) (*ssh.Client, error) {
	port := s.SshPort
	if port == 0 {
		port = 22
	}
	var auths []ssh.AuthMethod
	if s.SshAuthType == "key" && s.SshKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(s.SshKey))
		if err != nil {
			return nil, fmt.Errorf("私钥解析失败: %w", err)
		}
		auths = append(auths, ssh.PublicKeys(signer))
	} else {
		auths = append(auths, ssh.Password(s.SshPassword))
	}
	cfg := &ssh.ClientConfig{
		User:            s.SshUser,
		Auth:            auths,
		Timeout:         10 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := net.JoinHostPort(s.HostIP, fmt.Sprintf("%d", port))
	return ssh.Dial("tcp", addr, cfg)
}

// ServerTerminal 服务器 SSH 终端：桥接 websocket 与 ssh pty
func (t *TerminalService) ServerTerminal(s biz.Server, ws *websocket.Conn) error {
	client, err := sshDial(s)
	if err != nil {
		return err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}
	w := &wsWriter{ws: ws}
	session.Stdout = w
	session.Stderr = w

	modes := ssh.TerminalModes{ssh.ECHO: 1, ssh.TTY_OP_ISPEED: 14400, ssh.TTY_OP_OSPEED: 14400}
	if err := session.RequestPty("xterm-256color", 50, 200, modes); err != nil {
		return err
	}
	if err := session.Shell(); err != nil {
		return err
	}

	// 前端输入 -> ssh stdin（JSON 协议：data/resize）
	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				_ = stdin.Close()
				return
			}
			var m struct {
				Type string `json:"type"`
				Data string `json:"data"`
				Cols int    `json:"cols"`
				Rows int    `json:"rows"`
			}
			if json.Unmarshal(msg, &m) == nil {
				switch m.Type {
				case "resize":
					_ = session.WindowChange(m.Rows, m.Cols)
					continue
				case "data":
					_, _ = stdin.Write([]byte(m.Data))
					continue
				}
			}
			_, _ = stdin.Write(msg) // 兼容纯文本
		}
	}()

	return session.Wait()
}

// DockerLogsStream Docker 容器日志实时流（follow），解析 multiplexed 分帧
func (t *TerminalService) DockerLogsStream(d biz.DockerHost, containerID string, ws *websocket.Conn) error {
	resp, err := newDockerClient(d.Dsn()).do(http.MethodGet, "/containers/"+containerID+"/logs?follow=true&stdout=true&stderr=true&tail=300")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	w := &wsWriter{ws: ws}
	header := make([]byte, 8)
	for {
		if _, err := io.ReadFull(resp.Body, header); err != nil {
			break
		}
		n := binary.BigEndian.Uint32(header[4:8])
		if n == 0 {
			continue
		}
		buf := make([]byte, n)
		if _, err := io.ReadFull(resp.Body, buf); err != nil {
			break
		}
		_, _ = w.Write(buf)
	}
	return nil
}
