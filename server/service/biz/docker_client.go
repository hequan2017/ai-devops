package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// ===== Docker Engine API 轻量客户端（基于 net/http，零外部依赖）=====

type dockerClient struct {
	host    string // 基址，如 http://unix 或 tcp://1.2.3.4:2375
	httpCli *http.Client
}

func newDockerClient(dsn string) *dockerClient {
	timeout := 15 * time.Second
	if strings.HasPrefix(dsn, "unix://") {
		socket := strings.TrimPrefix(dsn, "unix://")
		return &dockerClient{
			host: "http://unix",
			httpCli: &http.Client{
				Timeout: timeout,
				Transport: &http.Transport{
					DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
						return net.DialTimeout("unix", socket, timeout)
					},
				},
			},
		}
	}
	return &dockerClient{
		host:    dsn,
		httpCli: &http.Client{Timeout: timeout},
	}
}

func (c *dockerClient) do(method, path string) (*http.Response, error) {
	req, err := http.NewRequest(method, c.host+path, nil)
	if err != nil {
		return nil, err
	}
	return c.httpCli.Do(req)
}

// DockerVersion Docker 版本摘要
type DockerVersion struct {
	Version string `json:"Version"`
}

// DockerContainer 容器摘要
type DockerContainer struct {
	ID     string   `json:"Id" gorm:"-"`
	Names  []string `json:"Names" gorm:"-"`
	Image  string   `json:"Image" gorm:"-"`
	State  string   `json:"State" gorm:"-"`
	Status string   `json:"Status" gorm:"-"`
}

// DockerImage 镜像摘要
type DockerImage struct {
	ID       string   `json:"Id" gorm:"-"`
	RepoTags []string `json:"RepoTags" gorm:"-"`
	Size     int64    `json:"Size" gorm:"-"`
}

// Ping 测试连通性并返回 Docker 版本
func (c *dockerClient) Ping() (string, error) {
	resp, err := c.do(http.MethodGet, "/version")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("docker 返回状态码 %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	var v DockerVersion
	_ = json.Unmarshal(body, &v)
	return v.Version, nil
}

// ListContainers 列出全部容器（含已停止）
func (c *dockerClient) ListContainers() ([]DockerContainer, error) {
	resp, err := c.do(http.MethodGet, "/containers/json?all=1")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("docker 返回状态码 %d", resp.StatusCode)
	}
	var list []DockerContainer
	err = json.NewDecoder(resp.Body).Decode(&list)
	return list, err
}

// ListImages 列出本地镜像
func (c *dockerClient) ListImages() ([]DockerImage, error) {
	resp, err := c.do(http.MethodGet, "/images/json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("docker 返回状态码 %d", resp.StatusCode)
	}
	var list []DockerImage
	err = json.NewDecoder(resp.Body).Decode(&list)
	return list, err
}

// ContainerAction 容器操作: start / stop / restart / remove
func (c *dockerClient) ContainerAction(id, action string) error {
	var method, path string
	switch action {
	case "start", "stop", "restart":
		method = http.MethodPost
		path = "/containers/" + id + "/" + action
	case "remove":
		method = http.MethodDelete
		path = "/containers/" + id + "?force=true"
	default:
		return fmt.Errorf("不支持的容器操作: %s", action)
	}
	resp, err := c.do(method, path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("docker 返回状态码 %d", resp.StatusCode)
	}
	return nil
}

// ContainerLogs 获取容器最近日志
func (c *dockerClient) ContainerLogs(id string) (string, error) {
	resp, err := c.do(http.MethodGet, "/containers/"+id+"/logs?stdout=true&stderr=true&tail=500")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("docker 返回状态码 %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

// DockerNetwork 网络摘要
type DockerNetwork struct {
	ID     string `json:"Id"`
	Name   string `json:"Name"`
	Driver string `json:"Driver"`
	Scope  string `json:"Scope"`
}

// ListNetworks 网络列表
func (c *dockerClient) ListNetworks() ([]DockerNetwork, error) {
	resp, err := c.do(http.MethodGet, "/networks")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("docker 返回状态码 %d", resp.StatusCode)
	}
	var list []DockerNetwork
	err = json.NewDecoder(resp.Body).Decode(&list)
	return list, err
}

// DockerVolume 数据卷摘要
type DockerVolume struct {
	Name       string `json:"Name"`
	Driver     string `json:"Driver"`
	Mountpoint string `json:"Mountpoint"`
}

type dockerVolumeList struct {
	Volumes []DockerVolume `json:"Volumes"`
}

// ListVolumes 数据卷列表
func (c *dockerClient) ListVolumes() ([]DockerVolume, error) {
	resp, err := c.do(http.MethodGet, "/volumes")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("docker 返回状态码 %d", resp.StatusCode)
	}
	var l dockerVolumeList
	err = json.NewDecoder(resp.Body).Decode(&l)
	return l.Volumes, err
}
