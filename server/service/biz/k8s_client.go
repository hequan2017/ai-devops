package biz

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ===== Kubernetes REST API 轻量客户端（基于 net/http，零外部依赖）=====

type k8sClient struct {
	apiServer string
	token     string
	httpCli   *http.Client
}

func newK8sClient(apiServer, token string, insecure bool) *k8sClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
	}
	return &k8sClient{
		apiServer: strings.TrimRight(strings.TrimSpace(apiServer), "/"),
		token:     strings.TrimSpace(token),
		httpCli:   &http.Client{Timeout: 15 * time.Second, Transport: tr},
	}
}

func (c *k8sClient) do(method, path string) (*http.Response, error) {
	req, err := http.NewRequest(method, c.apiServer+path, nil)
	if err != nil {
		return nil, err
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	return c.httpCli.Do(req)
}

// k8sVersion 集群版本摘要
type k8sVersion struct {
	GitVersion string `json:"gitVersion"`
}

// k8sList 通用列表包装
type k8sList struct {
	Items []map[string]interface{} `json:"items"`
}

// Ping 测试连通性并返回集群 gitVersion
func (c *k8sClient) Ping() (string, error) {
	resp, err := c.do(http.MethodGet, "/version")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("k8s 返回状态码 %d", resp.StatusCode)
	}
	var v k8sVersion
	_ = json.NewDecoder(resp.Body).Decode(&v)
	return v.GitVersion, nil
}

// Namespaces 命名空间列表（原始 items）
func (c *k8sClient) Namespaces() ([]map[string]interface{}, error) {
	return c.list("/api/v1/namespaces")
}

// Pods 指定命名空间下的 Pod 列表
func (c *k8sClient) Pods(ns string) ([]map[string]interface{}, error) {
	return c.list("/api/v1/namespaces/" + ns + "/pods")
}

// Nodes 集群节点列表
func (c *k8sClient) Nodes() ([]map[string]interface{}, error) {
	return c.list("/api/v1/nodes")
}

// Deployments 指定命名空间的 Deployment 列表
func (c *k8sClient) Deployments(ns string) ([]map[string]interface{}, error) {
	return c.list("/apis/apps/v1/namespaces/" + ns + "/deployments")
}

// Services 指定命名空间的 Service 列表
func (c *k8sClient) Services(ns string) ([]map[string]interface{}, error) {
	return c.list("/api/v1/namespaces/" + ns + "/services")
}

// PodLogs 获取 Pod 日志
func (c *k8sClient) PodLogs(ns, pod string) (string, error) {
	resp, err := c.do(http.MethodGet, "/api/v1/namespaces/"+ns+"/pods/"+pod+"/log?tailLines=500")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("k8s 返回状态码 %d: %s", resp.StatusCode, string(body))
	}
	body, _ := io.ReadAll(resp.Body)
	return string(body), nil
}

// DeletePod 删除 Pod
func (c *k8sClient) DeletePod(ns, pod string) error {
	resp, err := c.do(http.MethodDelete, "/api/v1/namespaces/"+ns+"/pods/"+pod)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("k8s 返回状态码 %d", resp.StatusCode)
	}
	return nil
}

func (c *k8sClient) list(path string) ([]map[string]interface{}, error) {
	resp, err := c.do(http.MethodGet, path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("k8s 返回状态码 %d: %s", resp.StatusCode, string(body))
	}
	var l k8sList
	err = json.NewDecoder(resp.Body).Decode(&l)
	return l.Items, err
}
