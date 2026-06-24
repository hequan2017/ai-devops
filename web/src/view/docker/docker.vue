<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="关键字">
          <el-input v-model="searchInfo.keyword" placeholder="接入点名称" clearable />
        </el-form-item>
        <el-form-item label="协议">
          <el-select v-model="searchInfo.protocol" placeholder="协议" clearable style="width:120px">
            <el-option label="unix" value="unix" />
            <el-option label="tcp" value="tcp" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增接入点</el-button>
      </div>
      <el-table :data="tableData" row-key="ID" highlight-current-row @current-change="handleCurrentChange">
        <el-table-column label="名称" prop="name" min-width="140" show-overflow-tooltip />
        <el-table-column label="协议" prop="protocol" width="80" />
        <el-table-column label="主机/Socket" prop="host" min-width="180" show-overflow-tooltip />
        <el-table-column label="端口" prop="port" width="80" />
        <el-table-column label="状态" prop="status" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'online' ? 'success' : 'info'" effect="dark">
              {{ scope.row.status === 'online' ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="版本" prop="version" width="110" />
        <el-table-column label="操作" fixed="right" min-width="280">
          <template #default="scope">
            <el-button type="primary" link @click="testHost(scope.row)">连接测试</el-button>
            <el-button type="primary" link icon="edit" @click="editHost(scope.row)">编辑</el-button>
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

    <!-- 资源面板 -->
    <el-card v-if="currentHost.ID" class="mt-4">
      <template #header>
        <div class="flex items-center justify-between">
          <span>{{ currentHost.name }} - 资源管理</span>
          <el-button icon="refresh" link @click="loadResource">刷新</el-button>
        </div>
      </template>
      <el-tabs v-model="activeTab" @tab-change="loadResource">
        <el-tab-pane label="容器" name="containers">
          <el-table :data="containers" size="small">
            <el-table-column label="容器名" min-width="160">
              <template #default="scope">{{ (scope.row.Names||[''])[0] }}</template>
            </el-table-column>
            <el-table-column label="镜像" prop="Image" min-width="160" show-overflow-tooltip />
            <el-table-column label="状态" prop="Status" width="140" />
            <el-table-column label="操作" width="280">
              <template #default="scope">
                <el-button type="primary" link @click="doAction(scope.row.Id,'start')">启动</el-button>
                <el-button type="warning" link @click="doAction(scope.row.Id,'stop')">停止</el-button>
                <el-button type="primary" link @click="doAction(scope.row.Id,'restart')">重启</el-button>
                <el-button type="primary" link @click="viewLogs(scope.row.Id)">日志</el-button>
                <el-button type="primary" link @click="openLiveLogs(scope.row.Id)">实时日志</el-button>
                <el-button type="primary" link @click="openContainerTerminal(scope.row.Id)">终端</el-button>
                <el-button type="danger" link @click="doAction(scope.row.Id,'remove')">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="镜像" name="images">
          <el-table :data="images" size="small">
            <el-table-column label="RepoTags" min-width="220">
              <template #default="scope">{{ (scope.row.RepoTags||[]).join(', ') }}</template>
            </el-table-column>
            <el-table-column label="ID" prop="Id" min-width="200" show-overflow-tooltip />
            <el-table-column label="大小(MB)" width="110">
              <template #default="scope">{{ (scope.row.Size/1024/1024).toFixed(2) }}</template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="网络" name="networks">
          <el-table :data="networks" size="small">
            <el-table-column label="名称" prop="Name" min-width="160" />
            <el-table-column label="ID" prop="Id" min-width="200" show-overflow-tooltip />
            <el-table-column label="驱动" prop="Driver" width="120" />
            <el-table-column label="范围" prop="Scope" width="100" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="数据卷" name="volumes">
          <el-table :data="volumes" size="small">
            <el-table-column label="名称" prop="Name" min-width="180" />
            <el-table-column label="驱动" prop="Driver" width="120" />
            <el-table-column label="挂载点" prop="Mountpoint" min-width="240" show-overflow-tooltip />
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 新增/编辑 -->
    <el-drawer destroy-on-close size="560" v-model="dialogVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增 Docker 接入点' : '编辑 Docker 接入点' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>
      <el-form :model="formData" label-position="top">
        <el-form-item label="名称"><el-input v-model="formData.name" placeholder="接入点名称" /></el-form-item>
        <el-row :gutter="20">
          <el-col :span="10">
            <el-form-item label="协议">
              <el-select v-model="formData.protocol" class="w-full">
                <el-option label="unix (本地socket)" value="unix" />
                <el-option label="tcp (远程)" value="tcp" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="14">
            <el-form-item label="主机/Socket路径">
              <el-input v-model="formData.host" :placeholder="formData.protocol==='unix' ? '/var/run/docker.sock' : '192.168.1.10'" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="端口(tcp)"><el-input v-model="formData.port" placeholder="2375" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="formData.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
    </el-drawer>

    <!-- 日志 -->
    <el-drawer destroy-on-close size="720" v-model="logsVisible" title="容器日志">
      <pre style="white-space:pre-wrap;font-size:12px;max-height:70vh;overflow:auto">{{ logsContent }}</pre>
    </el-drawer>

    <!-- 实时日志流 -->
    <el-drawer destroy-on-close size="1000" v-model="liveLogsVisible" title="实时日志">
      <div style="height:calc(100vh - 60px);background:#000;padding:8px">
        <XTerminal v-if="liveLogsVisible" :url="liveLogsUrl" :interactive="false" />
      </div>
    </el-drawer>

    <!-- 容器终端 -->
    <el-drawer destroy-on-close size="100%" v-model="execTermVisible" title="容器终端">
      <div style="height:calc(100vh - 60px);background:#000;padding:8px">
        <XTerminal v-if="execTermVisible" :url="execTermUrl" />
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
  import XTerminal from '@/components/terminal/xTerminal.vue'
  import {
    getDockerHostList, createDockerHost, updateDockerHost, deleteDockerHost, findDockerHost,
    testDockerHost, getContainers, containerAction, getContainerLogs, getImages, getNetworks, getVolumes
  } from '@/api/docker'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref } from 'vue'

  defineOptions({ name: 'Docker' })

  const page = ref(1); const pageSize = ref(10); const total = ref(0)
  const tableData = ref([]); const searchInfo = ref({})
  const currentHost = ref({}); const activeTab = ref('containers')
  const containers = ref([]); const images = ref([]); const networks = ref([]); const volumes = ref([])

  const getTableData = async () => {
    const res = await getDockerHostList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (res.code === 0) { tableData.value = res.data.list; total.value = res.data.total }
  }
  getTableData()

  const onSubmit = () => { page.value = 1; getTableData() }
  const onReset = () => { searchInfo.value = {}; getTableData() }

  const handleCurrentChange = (row) => { if (row) { currentHost.value = row; loadResource() } }

  const loadResource = async () => {
    if (!currentHost.value.ID) return
    const id = currentHost.value.ID
    if (activeTab.value === 'containers') { const r = await getContainers({ ID: id }); if (r.code === 0) containers.value = r.data || [] }
    if (activeTab.value === 'images') { const r = await getImages({ ID: id }); if (r.code === 0) images.value = r.data || [] }
    if (activeTab.value === 'networks') { const r = await getNetworks({ ID: id }); if (r.code === 0) networks.value = r.data || [] }
    if (activeTab.value === 'volumes') { const r = await getVolumes({ ID: id }); if (r.code === 0) volumes.value = r.data || [] }
  }

  const testHost = async (row) => {
    const r = await testDockerHost({ ID: row.ID })
    if (r.code === 0) { ElMessage.success('连接成功, 版本: ' + (r.data.version || '')); getTableData() }
  }

  const doAction = async (cid, action) => {
    const r = await containerAction({ hostId: currentHost.value.ID, containerId: cid, action })
    if (r.code === 0) { ElMessage.success('操作成功'); loadResource() }
  }

  const logsVisible = ref(false); const logsContent = ref('')
  const viewLogs = async (cid) => {
    const r = await getContainerLogs({ ID: currentHost.value.ID, containerId: cid })
    if (r.code === 0) { logsContent.value = r.data.logs || '(空)'; logsVisible.value = true }
  }

  // 实时日志流 (WebSocket)
  const liveLogsVisible = ref(false); const liveLogsUrl = ref('')
  const buildWsBase = () => {
    const apiBase = import.meta.env.VITE_BASE_API || ''
    if (apiBase.startsWith('http')) return apiBase.replace(/^http/, 'ws').replace(/\/$/, '')
    return (location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + apiBase
  }
  const openLiveLogs = (cid) => {
    const token = encodeURIComponent(localStorage.getItem('token') || '')
    liveLogsUrl.value = `${buildWsBase()}/docker/containerLogsStream?hostId=${currentHost.value.ID}&containerId=${cid}&token=${token}`
    liveLogsVisible.value = true
  }

  // 容器交互式终端
  const execTermVisible = ref(false); const execTermUrl = ref('')
  const openContainerTerminal = (cid) => {
    const token = encodeURIComponent(localStorage.getItem('token') || '')
    execTermUrl.value = `${buildWsBase()}/docker/containerTerminal?hostId=${currentHost.value.ID}&containerId=${cid}&token=${token}`
    execTermVisible.value = true
  }

  const type = ref(''); const dialogVisible = ref(false)
  const defaultForm = () => ({ name: '', protocol: 'unix', host: '/var/run/docker.sock', port: '2375', remark: '' })
  const formData = ref(defaultForm())

  const openDialog = () => { type.value = 'create'; formData.value = defaultForm(); dialogVisible.value = true }
  const editHost = async (row) => { const r = await findDockerHost({ ID: row.ID }); if (r.code === 0) { formData.value = r.data; type.value = 'update'; dialogVisible.value = true } }
  const closeDialog = () => { dialogVisible.value = false; formData.value = defaultForm() }
  const enterDialog = async () => {
    const res = type.value === 'update' ? await updateDockerHost(formData.value) : await createDockerHost(formData.value)
    if (res.code === 0) { ElMessage.success('保存成功'); closeDialog(); getTableData() }
  }
  const deleteRow = (row) => {
    ElMessageBox.confirm('确定删除该接入点?', '提示', { type: 'warning' }).then(async () => {
      const r = await deleteDockerHost({ ID: row.ID })
      if (r.code === 0) { ElMessage.success('删除成功'); getTableData() }
    })
  }
</script>
