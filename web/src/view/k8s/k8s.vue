<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="集群">
          <el-input v-model="searchInfo.keyword" placeholder="集群名称" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增集群</el-button>
      </div>
      <el-table :data="tableData" row-key="ID" highlight-current-row @current-change="handleCurrentChange">
        <el-table-column label="集群名称" prop="name" min-width="140" show-overflow-tooltip />
        <el-table-column label="API Server" prop="apiServer" min-width="220" show-overflow-tooltip />
        <el-table-column label="状态" prop="status" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'online' ? 'success' : 'info'" effect="dark">
              {{ scope.row.status === 'online' ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="版本" prop="version" width="130" />
        <el-table-column label="操作" fixed="right" min-width="260">
          <template #default="scope">
            <el-button type="primary" link @click="testCluster(scope.row)">连接测试</el-button>
            <el-button type="primary" link icon="edit" @click="editCluster(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper"
          :current-page="page" :page-size="pageSize" :page-sizes="[10,30,50,100]" :total="total"
          @current-change="(v)=>{page=v;getTableData()}" @size-change="(v)=>{pageSize=v;getTableData()}" />
      </div>
    </div>

    <el-card v-if="currentCluster.ID" class="mt-4">
      <template #header>
        <div class="flex items-center justify-between">
          <span>{{ currentCluster.name }} - 资源浏览</span>
          <div class="flex items-center gap-2">
            <el-select v-model="namespace" placeholder="命名空间" style="width:180px" @change="loadResource">
              <el-option v-for="ns in namespaces" :key="ns.name" :label="ns.name" :value="ns.name" />
            </el-select>
            <el-button icon="refresh" link @click="loadResource">刷新</el-button>
          </div>
        </div>
      </template>
      <el-tabs v-model="activeTab" @tab-change="loadResource">
        <el-tab-pane label="Pod" name="pods">
          <el-table :data="pods" size="small">
            <el-table-column label="名称" min-width="200">
              <template #default="scope">{{ scope.row.metadata?.name }}</template>
            </el-table-column>
            <el-table-column label="状态" width="110">
              <template #default="scope">
                <el-tag size="small" :type="podTagType(scope.row.status?.phase)">{{ scope.row.status?.phase }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="节点" width="150">
              <template #default="scope">{{ scope.row.spec?.nodeName }}</template>
            </el-table-column>
            <el-table-column label="重启" width="80">
              <template #default="scope">{{ (scope.row.status?.containerStatuses||[]).reduce((s,c)=>s+(c.restartCount||0),0) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="180">
              <template #default="scope">
                <el-button type="primary" link @click="viewPodLogs(scope.row.metadata.name)">日志</el-button>
                <el-button type="danger" link @click="removePod(scope.row.metadata.name)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="Deployment" name="deployments">
          <el-table :data="deployments" size="small">
            <el-table-column label="名称" min-width="200">
              <template #default="scope">{{ scope.row.metadata?.name }}</template>
            </el-table-column>
            <el-table-column label="副本" width="120">
              <template #default="scope">{{ scope.row.status?.readyReplicas||0 }} / {{ scope.row.status?.replicas||0 }}</template>
            </el-table-column>
            <el-table-column label="镜像" min-width="220">
              <template #default="scope">{{ (scope.row.spec?.template?.spec?.containers||[]).map(c=>c.image).join(', ') }}</template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="Service" name="services">
          <el-table :data="services" size="small">
            <el-table-column label="名称" min-width="160">
              <template #default="scope">{{ scope.row.metadata?.name }}</template>
            </el-table-column>
            <el-table-column label="类型" width="120">
              <template #default="scope">{{ scope.row.spec?.type }}</template>
            </el-table-column>
            <el-table-column label="ClusterIP" width="150">
              <template #default="scope">{{ scope.row.spec?.clusterIP }}</template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="Node" name="nodes">
          <el-table :data="nodes" size="small">
            <el-table-column label="名称" min-width="160">
              <template #default="scope">{{ scope.row.metadata?.name }}</template>
            </el-table-column>
            <el-table-column label="状态" width="110">
              <template #default="scope">{{ (scope.row.status?.conditions||[]).find(c=>c.type==='Ready')?.status }}</template>
            </el-table-column>
            <el-table-column label="版本" width="140">
              <template #default="scope">{{ scope.row.status?.nodeInfo?.kubeletVersion }}</template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <el-drawer destroy-on-close size="560" v-model="dialogVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增 K8s 集群' : '编辑 K8s 集群' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>
      <el-form :model="formData" label-position="top">
        <el-form-item label="集群名称"><el-input v-model="formData.name" /></el-form-item>
        <el-form-item label="API Server"><el-input v-model="formData.apiServer" placeholder="https://1.2.3.4:6443" /></el-form-item>
        <el-form-item label="Token"><el-input v-model="formData.token" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="跳过TLS校验(自签证书)">
          <el-switch v-model="formData.insecureTls" />
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="formData.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer destroy-on-close size="720" v-model="logsVisible" title="Pod 日志">
      <pre style="white-space:pre-wrap;font-size:12px;max-height:70vh;overflow:auto">{{ logsContent }}</pre>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    getK8sClusterList, createK8sCluster, updateK8sCluster, deleteK8sCluster, findK8sCluster,
    testK8sCluster, getNamespaces, getPods, getNodes, getDeployments, getServices, getPodLogs, deletePod
  } from '@/api/k8s'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref } from 'vue'

  defineOptions({ name: 'K8s' })

  const page = ref(1); const pageSize = ref(10); const total = ref(0)
  const tableData = ref([]); const searchInfo = ref({})
  const currentCluster = ref({}); const activeTab = ref('pods')
  const namespace = ref('default')
  const namespaces = ref([]); const pods = ref([]); const nodes = ref([]); const deployments = ref([]); const services = ref([])

  const getTableData = async () => {
    const res = await getK8sClusterList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (res.code === 0) { tableData.value = res.data.list; total.value = res.data.total }
  }
  getTableData()

  const onSubmit = () => { page.value = 1; getTableData() }
  const onReset = () => { searchInfo.value = {}; getTableData() }

  const handleCurrentChange = async (row) => {
    if (!row) return
    currentCluster.value = row
    const r = await getNamespaces({ ID: row.ID })
    if (r.code === 0) {
      namespaces.value = (r.data || []).map((i) => ({ name: i.metadata?.name }))
      loadResource()
    }
  }

  const loadResource = async () => {
    if (!currentCluster.value.ID) return
    const id = currentCluster.value.ID; const ns = namespace.value
    if (activeTab.value === 'pods') { const r = await getPods({ ID: id, namespace: ns }); if (r.code === 0) pods.value = r.data || [] }
    if (activeTab.value === 'deployments') { const r = await getDeployments({ ID: id, namespace: ns }); if (r.code === 0) deployments.value = r.data || [] }
    if (activeTab.value === 'services') { const r = await getServices({ ID: id, namespace: ns }); if (r.code === 0) services.value = r.data || [] }
    if (activeTab.value === 'nodes') { const r = await getNodes({ ID: id }); if (r.code === 0) nodes.value = r.data || [] }
  }

  const podTagType = (phase) => (phase === 'Running' ? 'success' : phase === 'Failed' ? 'danger' : 'warning')

  const testCluster = async (row) => {
    const r = await testK8sCluster({ ID: row.ID })
    if (r.code === 0) { ElMessage.success('连接成功: ' + (r.data.version || '')); getTableData() }
  }

  const logsVisible = ref(false); const logsContent = ref('')
  const viewPodLogs = async (podName) => {
    const r = await getPodLogs({ ID: currentCluster.value.ID, namespace: namespace.value, pod: podName })
    if (r.code === 0) { logsContent.value = r.data.logs || '(空)'; logsVisible.value = true }
  }
  const removePod = (podName) => {
    ElMessageBox.confirm('确定删除 Pod ' + podName + '?', '提示', { type: 'warning' }).then(async () => {
      const r = await deletePod({ clusterId: currentCluster.value.ID, namespace: namespace.value, pod: podName })
      if (r.code === 0) { ElMessage.success('已删除'); loadResource() }
    })
  }

  const type = ref(''); const dialogVisible = ref(false)
  const defaultForm = () => ({ name: '', apiServer: '', token: '', insecureTls: true, remark: '' })
  const formData = ref(defaultForm())
  const openDialog = () => { type.value = 'create'; formData.value = defaultForm(); dialogVisible.value = true }
  const editCluster = async (row) => { const r = await findK8sCluster({ ID: row.ID }); if (r.code === 0) { formData.value = r.data; type.value = 'update'; dialogVisible.value = true } }
  const closeDialog = () => { dialogVisible.value = false; formData.value = defaultForm() }
  const enterDialog = async () => {
    const res = type.value === 'update' ? await updateK8sCluster(formData.value) : await createK8sCluster(formData.value)
    if (res.code === 0) { ElMessage.success('保存成功'); closeDialog(); getTableData() }
  }
  const deleteRow = (row) => {
    ElMessageBox.confirm('确定删除该集群接入点?', '提示', { type: 'warning' }).then(async () => {
      const r = await deleteK8sCluster({ ID: row.ID })
      if (r.code === 0) { ElMessage.success('删除成功'); getTableData() }
    })
  }
</script>
