package biz

import (
	"ai-devops/server/model/biz"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// DockerExecTerminal Docker 容器交互式终端（exec + hijack）
// 流程：创建 exec -> 用 raw conn 以 Upgrade 方式 start -> 双向桥接 ws 与容器 stdin/stdout
func (t *TerminalService) DockerExecTerminal(d biz.DockerHost, containerID string, ws *websocket.Conn) error {
	cli := newDockerClient(d.Dsn())

	// 1. 创建 exec 实例
	createBody, _ := json.Marshal(map[string]interface{}{
		"AttachStdin": true, "AttachStdout": true, "AttachStderr": true,
		"Tty": true, "Cmd": []string{"sh"},
	})
	createReq, err := http.NewRequest(http.MethodPost, cli.host+"/containers/"+containerID+"/exec", strings.NewReader(string(createBody)))
	if err != nil {
		return err
	}
	createReq.Header.Set("Content-Type", "application/json")
	createResp, err := cli.httpCli.Do(createReq)
	if err != nil {
		return err
	}
	defer createResp.Body.Close()
	var execInfo struct {
		Id string `json:"Id"`
	}
	_ = json.NewDecoder(createResp.Body).Decode(&execInfo)
	if execInfo.Id == "" {
		return fmt.Errorf("创建exec失败")
	}

	// 2. hijack raw conn 启动 exec
	conn, err := dialDocker(d)
	if err != nil {
		return err
	}
	defer conn.Close()

	startBody := `{"Detach":false,"Tty":true}`
	startReq := fmt.Sprintf("POST /exec/%s/start HTTP/1.1\r\nHost: docker\r\nContent-Type: application/json\r\nConnection: Upgrade\r\nUpgrade: tcp\r\nContent-Length: %d\r\n\r\n%s",
		execInfo.Id, len(startBody), startBody)
	if _, err := conn.Write([]byte(startReq)); err != nil {
		return err
	}

	// 读响应头直到空行（之后为 raw 流）
	br := bufio.NewReader(conn)
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			return err
		}
		if line == "\r\n" || line == "\n" {
			break
		}
	}

	// 3. 双向桥接：ws -> 容器 stdin；容器 stdout -> ws
	done := make(chan struct{})
	go func() {
		_, _ = io.Copy(conn, &wsReader{ws: ws})
		if c, ok := conn.(interface{ CloseWrite() error }); ok {
			_ = c.CloseWrite()
		}
		close(done)
	}()
	_, _ = io.Copy(&wsWriter{ws: ws}, conn)
	<-done
	return nil
}

// dialDocker 建立到 Docker daemon 的 raw 连接（unix socket / tcp）
func dialDocker(d biz.DockerHost) (net.Conn, error) {
	dsn := d.Dsn()
	if strings.HasPrefix(dsn, "unix://") {
		return net.DialTimeout("unix", strings.TrimPrefix(dsn, "unix://"), 10*time.Second)
	}
	addr := strings.TrimPrefix(dsn, "tcp://")
	return net.DialTimeout("tcp", addr, 10*time.Second)
}

// wsReader 将 websocket 作为 io.Reader（逐消息消费）
type wsReader struct {
	ws  *websocket.Conn
	buf []byte
}

func (r *wsReader) Read(p []byte) (int, error) {
	if len(r.buf) == 0 {
		_, msg, err := r.ws.ReadMessage()
		if err != nil {
			return 0, err
		}
		r.buf = msg
	}
	n := copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}
