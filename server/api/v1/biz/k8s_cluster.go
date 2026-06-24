package biz

import (
	"ai-devops/server/global"
	"ai-devops/server/model/biz"
	bizReq "ai-devops/server/model/biz/request"
	"ai-devops/server/model/common/request"
	"ai-devops/server/model/common/response"
	"ai-devops/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type K8sApi struct{}

// CreateK8sCluster 新增集群
// @Router /k8s/createK8sCluster [post]
func (k *K8sApi) CreateK8sCluster(c *gin.Context) {
	var cluster biz.K8sCluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := k8sClusterService.CreateK8sCluster(cluster); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// DeleteK8sCluster 删除集群
// @Router /k8s/deleteK8sCluster [delete]
func (k *K8sApi) DeleteK8sCluster(c *gin.Context) {
	var cluster biz.K8sCluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(cluster.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := k8sClusterService.DeleteK8sCluster(cluster); err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteK8sClusterByIds 批量删除
// @Router /k8s/deleteK8sClusterByIds [delete]
func (k *K8sApi) DeleteK8sClusterByIds(c *gin.Context) {
	var ids request.IdsReq
	if err := c.ShouldBindJSON(&ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := k8sClusterService.DeleteK8sClusterByIds(ids); err != nil {
		response.FailWithMessage("批量删除失败", c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateK8sCluster 更新集群
// @Router /k8s/updateK8sCluster [put]
func (k *K8sApi) UpdateK8sCluster(c *gin.Context) {
	var cluster biz.K8sCluster
	if err := c.ShouldBindJSON(&cluster); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(cluster.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := k8sClusterService.UpdateK8sCluster(cluster); err != nil {
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindK8sCluster 查询集群详情
// @Router /k8s/findK8sCluster [get]
func (k *K8sApi) FindK8sCluster(c *gin.Context) {
	var cluster biz.K8sCluster
	if err := c.ShouldBindQuery(&cluster); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(cluster.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := k8sClusterService.GetK8sCluster(cluster.ID)
	if err != nil {
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetK8sClusterList 分页列表
// @Router /k8s/getK8sClusterList [get]
func (k *K8sApi) GetK8sClusterList(c *gin.Context) {
	var pageInfo bizReq.K8sClusterSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := k8sClusterService.GetK8sClusterList(pageInfo)
	if err != nil {
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List: list, Total: total, Page: pageInfo.Page, PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// TestK8sCluster 连接测试
// @Router /k8s/testK8sCluster [get]
func (k *K8sApi) TestK8sCluster(c *gin.Context) {
	var cluster biz.K8sCluster
	if err := c.ShouldBindQuery(&cluster); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(cluster.GVA_MODEL, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	version, err := k8sClusterService.TestCluster(cluster.ID)
	if err != nil {
		response.FailWithMessage("连接失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]interface{}{"version": version}, "连接成功", c)
}

// GetNamespaces 命名空间列表
// @Router /k8s/getNamespaces [get]
func (k *K8sApi) GetNamespaces(c *gin.Context) {
	var cluster biz.K8sCluster
	_ = c.ShouldBindQuery(&cluster)
	list, err := k8sClusterService.ListNamespaces(cluster.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetPods Pod 列表
// @Router /k8s/getPods [get]
func (k *K8sApi) GetPods(c *gin.Context) {
	var cluster biz.K8sCluster
	_ = c.ShouldBindQuery(&cluster)
	ns := c.Query("namespace")
	list, err := k8sClusterService.ListPods(cluster.ID, ns)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetNodes 节点列表
// @Router /k8s/getNodes [get]
func (k *K8sApi) GetNodes(c *gin.Context) {
	var cluster biz.K8sCluster
	_ = c.ShouldBindQuery(&cluster)
	list, err := k8sClusterService.ListNodes(cluster.ID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetDeployments Deployment 列表
// @Router /k8s/getDeployments [get]
func (k *K8sApi) GetDeployments(c *gin.Context) {
	var cluster biz.K8sCluster
	_ = c.ShouldBindQuery(&cluster)
	ns := c.Query("namespace")
	list, err := k8sClusterService.ListDeployments(cluster.ID, ns)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetServices Service 列表
// @Router /k8s/getServices [get]
func (k *K8sApi) GetServices(c *gin.Context) {
	var cluster biz.K8sCluster
	_ = c.ShouldBindQuery(&cluster)
	ns := c.Query("namespace")
	list, err := k8sClusterService.ListServices(cluster.ID, ns)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}

// GetPodLogs Pod 日志
// @Router /k8s/getPodLogs [get]
func (k *K8sApi) GetPodLogs(c *gin.Context) {
	var cluster biz.K8sCluster
	_ = c.ShouldBindQuery(&cluster)
	ns := c.Query("namespace")
	pod := c.Query("pod")
	logs, err := k8sClusterService.PodLogs(cluster.ID, ns, pod)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]interface{}{"logs": logs}, "获取成功", c)
}

// k8sDeletePodReq 删除 Pod 入参
type k8sDeletePodReq struct {
	ClusterID uint   `json:"clusterId"`
	Namespace string `json:"namespace"`
	Pod       string `json:"pod"`
}

// DeletePod 删除 Pod
// @Router /k8s/deletePod [post]
func (k *K8sApi) DeletePod(c *gin.Context) {
	var req k8sDeletePodReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := k8sClusterService.DeletePod(req.ClusterID, req.Namespace, req.Pod); err != nil {
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
