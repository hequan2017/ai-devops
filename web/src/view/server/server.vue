<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="关键字" prop="keyword">
          <el-input v-model="searchInfo.keyword" placeholder="名称/SN/业务IP/管理IP" clearable />
        </el-form-item>
        <el-form-item label="厂商" prop="manufacturer">
          <el-select v-model="searchInfo.manufacturer" placeholder="厂商" clearable>
            <el-option v-for="item in manufacturerOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="searchInfo.status" placeholder="状态" clearable>
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
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
        <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
        <el-button type="warning" icon="tools" style="margin-left: 10px" :disabled="!multipleSelection.length" @click="openExec">批量执行</el-button>
        <el-button icon="delete" style="margin-left: 10px" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>
      <el-table ref="multipleTable" style="width: 100%" :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="创建日期" width="170">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="服务器名称" prop="name" min-width="130" show-overflow-tooltip />
        <el-table-column align="left" label="厂商" width="90">
          <template #default="scope">
            <el-tag :type="manufacturerTagType(scope.row.manufacturer)" effect="light">{{ manufacturerLabel(scope.row.manufacturer) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="型号" prop="model" min-width="120" show-overflow-tooltip />
        <el-table-column align="left" label="序列号" prop="serialNumber" min-width="140" show-overflow-tooltip />
        <el-table-column align="left" label="业务IP" prop="hostIp" width="120" />
        <el-table-column align="left" label="机房" prop="dataCenter" width="100" show-overflow-tooltip />
        <el-table-column align="left" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="statusTagType(scope.row.status)" effect="dark">{{ statusLabel(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="操作" fixed="right" min-width="280">
          <template #default="scope">
            <el-button type="primary" link icon="view" @click="getDetails(scope.row)">详情</el-button>
            <el-button v-if="scope.row.sshUser" type="primary" link icon="monitor" @click="openTerminal(scope.row)">终端</el-button>
            <el-button type="primary" link icon="edit" @click="updateServerFunc(scope.row)">变更</el-button>
            <el-dropdown class="ml-2" @command="(cmd)=>handleIpmi(scope.row, cmd)">
              <el-button type="primary" link>电源<el-icon class="el-icon--right"><ArrowDown /></el-icon></el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="status">查询状态</el-dropdown-item>
                  <el-dropdown-item command="on">开机</el-dropdown-item>
                  <el-dropdown-item command="soft">软关机</el-dropdown-item>
                  <el-dropdown-item command="off">强制关机</el-dropdown-item>
                  <el-dropdown-item command="reset">硬重启</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, sizes, prev, pager, next, jumper"
          :current-page="page" :page-size="pageSize" :page-sizes="[10,30,50,100]" :total="total"
          @current-change="handleCurrentChange" @size-change="handleSizeChange" />
      </div>
    </div>

    <el-drawer destroy-on-close size="680" v-model="dialogFormVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增服务器' : '编辑服务器' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>
      <el-form :model="formData" label-position="top" ref="elFormRef" :rules="rule">
        <el-form-item label="服务器名称:" prop="name"><el-input v-model="formData.name" clearable placeholder="服务器名称/主机名" /></el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="厂商:" prop="manufacturer">
              <el-select v-model="formData.manufacturer" class="w-full">
                <el-option v-for="item in manufacturerOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态:" prop="status">
              <el-select v-model="formData.status" class="w-full">
                <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="型号:" prop="model"><el-input v-model="formData.model" clearable placeholder="如 PowerEdge R740 / 2288H V5" /></el-form-item>
        <el-row :gutter="20">
          <el-col :span="12"><el-form-item label="序列号(SN):" prop="serialNumber"><el-input v-model="formData.serialNumber" clearable /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="业务IP:" prop="hostIp"><el-input v-model="formData.hostIp" clearable /></el-form-item></el-col>
        </el-row>
        <el-form-item label="带外管理IP:" prop="manageIp"><el-input v-model="formData.manageIp" clearable placeholder="iDRAC / iBMC 管理IP" /></el-form-item>
        <el-row :gutter="20">
          <el-col :span="8"><el-form-item label="机房:"><el-input v-model="formData.dataCenter" clearable /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="机柜:"><el-input v-model="formData.cabinet" clearable /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="U位:"><el-input v-model="formData.uPosition" clearable /></el-form-item></el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8"><el-form-item label="CPU:"><el-input v-model="formData.cpu" clearable /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="内存:"><el-input v-model="formData.memory" clearable /></el-form-item></el-col>
          <el-col :span="8"><el-form-item label="磁盘:"><el-input v-model="formData.disk" clearable /></el-form-item></el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12"><el-form-item label="操作系统:"><el-input v-model="formData.os" clearable /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="负责人:"><el-input v-model="formData.owner" clearable /></el-form-item></el-col>
        </el-row>
        <el-form-item label="备注:"><el-input v-model="formData.remark" type="textarea" :rows="2" clearable /></el-form-item>

        <el-divider content-position="left">IPMI / BMC 带外管理</el-divider>
        <el-row :gutter="20">
          <el-col :span="12"><el-form-item label="IPMI地址:"><el-input v-model="formData.ipmiIp" clearable placeholder="BMC IP" /></el-form-item></el-col>
          <el-col :span="6"><el-form-item label="用户名:"><el-input v-model="formData.ipmiUser" clearable /></el-form-item></el-col>
          <el-col :span="6"><el-form-item label="密码:"><el-input v-model="formData.ipmiPassword" show-password clearable /></el-form-item></el-col>
        </el-row>
      </el-form>
    </el-drawer>

    <el-drawer destroy-on-close size="560" v-model="detailShow" :show-close="true" :before-close="closeDetailShow">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="服务器名称">{{ detailForm.name }}</el-descriptions-item>
        <el-descriptions-item label="厂商">{{ manufacturerLabel(detailForm.manufacturer) }}</el-descriptions-item>
        <el-descriptions-item label="型号">{{ detailForm.model }}</el-descriptions-item>
        <el-descriptions-item label="序列号">{{ detailForm.serialNumber }}</el-descriptions-item>
        <el-descriptions-item label="业务IP">{{ detailForm.hostIp }}</el-descriptions-item>
        <el-descriptions-item label="带外管理IP">{{ detailForm.manageIp }}</el-descriptions-item>
        <el-descriptions-item label="机房/机柜/U位">{{ detailForm.dataCenter }} / {{ detailForm.cabinet }} / {{ detailForm.uPosition }}</el-descriptions-item>
        <el-descriptions-item label="CPU/内存/磁盘">{{ detailForm.cpu }} / {{ detailForm.memory }} / {{ detailForm.disk }}</el-descriptions-item>
        <el-descriptions-item label="操作系统">{{ detailForm.os }}</el-descriptions-item>
        <el-descriptions-item label="状态"><el-tag :type="statusTagType(detailForm.status)" effect="dark">{{ statusLabel(detailForm.status) }}</el-tag></el-descriptions-item>
        <el-descriptions-item label="负责人">{{ detailForm.owner }}</el-descriptions-item>
        <el-descriptions-item label="IPMI地址">{{ detailForm.ipmiIp }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ detailForm.remark }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>

    <!-- SSH 终端 -->
    <el-drawer v-model="termVisible" size="100%" :title="'SSH 终端 - ' + termName" destroy-on-close>
      <div style="height:calc(100vh - 60px);background:#000;padding:8px">
        <XTerminal v-if="termVisible" :url="termUrl" />
      </div>
    </el-drawer>

    <!-- 批量执行命令 -->
    <el-drawer v-model="execVisible" size="640" title="批量执行命令" destroy-on-close>
      <el-input v-model="execCmdText" type="textarea" :rows="3" placeholder="如：uptime && free -h" />
      <div class="mt-3">
        <el-button type="primary" :loading="execLoading" @click="handleExec">执行（{{ multipleSelection.length }}台）</el-button>
      </div>
      <el-divider />
      <div v-for="r in execResults" :key="r.serverId" class="mb-3">
        <div class="font-bold">
          {{ r.serverName }}
          <span v-if="r.error" class="text-red-500">[{{ r.error }}]</span>
        </div>
        <pre style="white-space:pre-wrap;font-size:12px;background:#f5f5f5;padding:8px;border-radius:4px;max-height:200px;overflow:auto">{{ r.output }}</pre>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
  import { ArrowDown } from '@element-plus/icons-vue'
  import XTerminal from '@/components/terminal/xTerminal.vue'
  import {
    createServer, deleteServer, deleteServerByIds, updateServer, findServer, getServerList, ipmiPower, execCmd
  } from '@/api/server'
  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref, reactive } from 'vue'

  defineOptions({ name: 'Server' })

  const manufacturerOptions = [
    { value: 'dell', label: '戴尔(Dell)' },
    { value: 'huawei', label: '华为(Huawei)' }
  ]
  const statusOptions = [
    { value: 'online', label: '运行中' },
    { value: 'offline', label: '已关机' },
    { value: 'maintenance', label: '维护中' },
    { value: 'fault', label: '故障' }
  ]
  const manufacturerLabel = (v) => manufacturerOptions.find((i) => i.value === v)?.label || v || '-'
  const manufacturerTagType = (v) => (v === 'huawei' ? 'danger' : 'primary')
  const statusLabel = (v) => statusOptions.find((i) => i.value === v)?.label || v || '-'
  const statusTagType = (v) => ({ online: 'success', offline: 'info', maintenance: 'warning', fault: 'danger' }[v] || 'info')

  const defaultForm = () => ({
    name: '', manufacturer: '', model: '', serialNumber: '',
    manageIp: '', hostIp: '', dataCenter: '', cabinet: '', uPosition: '',
    cpu: '', memory: '', disk: '', os: '',
    status: 'online', owner: '', remark: '',
    ipmiIp: '', ipmiUser: '', ipmiPassword: ''
  })
  const formData = ref(defaultForm())
  const rule = reactive({
    name: [{ required: true, message: '请输入服务器名称', trigger: ['input', 'blur'] }],
    manufacturer: [{ required: true, message: '请选择厂商', trigger: 'change' }],
    status: [{ required: true, message: '请选择状态', trigger: 'change' }]
  })
  const elFormRef = ref()
  const elSearchFormRef = ref()

  const page = ref(1); const total = ref(0); const pageSize = ref(10)
  const tableData = ref([]); const searchInfo = ref({})
  const onReset = () => { searchInfo.value = {}; getTableData() }
  const onSubmit = () => { page.value = 1; getTableData() }
  const handleSizeChange = (val) => { pageSize.value = val; getTableData() }
  const handleCurrentChange = (val) => { page.value = val; getTableData() }
  const getTableData = async () => {
    const table = await getServerList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (table.code === 0) { tableData.value = table.data.list; total.value = table.data.total; page.value = table.data.page; pageSize.value = table.data.pageSize }
  }
  getTableData()

  const multipleSelection = ref([])
  const handleSelectionChange = (val) => { multipleSelection.value = val }

  const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除该服务器吗?', '提示', { type: 'warning' }).then(async () => {
      const res = await deleteServer({ ID: row.ID })
      if (res.code === 0) { ElMessage.success('删除成功'); if (tableData.value.length === 1 && page.value > 1) page.value--; getTableData() }
    })
  }
  const onDelete = () => {
    ElMessageBox.confirm('确定要删除选中的服务器吗?', '提示', { type: 'warning' }).then(async () => {
      const ids = multipleSelection.value.map((item) => item.ID)
      const res = await deleteServerByIds({ ids })
      if (res.code === 0) { ElMessage.success('删除成功'); if (tableData.value.length === ids.length && page.value > 1) page.value--; getTableData() }
    })
  }

  // IPMI 电源控制
  const handleIpmi = (row, action) => {
    const labels = { status: '查询状态', on: '开机', soft: '软关机', off: '强制关机', reset: '硬重启' }
    ElMessageBox.confirm(`确定对「${row.name}」执行 ${labels[action]} ?`, 'IPMI 操作', { type: action === 'status' ? 'info' : 'warning' })
      .then(async () => {
        const res = await ipmiPower({ ID: row.ID, action })
        if (res.code === 0) { ElMessage.success(res.data.output || '操作成功') }
      }).catch(() => {})
  }

  const type = ref(''); const dialogFormVisible = ref(false)
  const openDialog = () => { type.value = 'create'; formData.value = defaultForm(); dialogFormVisible.value = true }
  const updateServerFunc = async (row) => {
    const res = await findServer({ ID: row.ID }); type.value = 'update'
    if (res.code === 0) { formData.value = res.data; dialogFormVisible.value = true }
  }
  const closeDialog = () => { dialogFormVisible.value = false; formData.value = defaultForm() }
  const enterDialog = () => {
    elFormRef.value?.validate(async (valid) => {
      if (!valid) return
      const res = type.value === 'update' ? await updateServer(formData.value) : await createServer(formData.value)
      if (res.code === 0) { ElMessage.success(type.value === 'update' ? '更新成功' : '创建成功'); closeDialog(); getTableData() }
    })
  }

  const detailShow = ref(false); const detailForm = ref({})
  const getDetails = async (row) => {
    const res = await findServer({ ID: row.ID })
    if (res.code === 0) { detailForm.value = res.data; detailShow.value = true }
  }
  const closeDetailShow = () => { detailShow.value = false; detailForm.value = {} }

  // WebSSH 终端
  const termVisible = ref(false)
  const termName = ref('')
  const termUrl = ref('')
  const buildWsBase = () => {
    const apiBase = import.meta.env.VITE_BASE_API || ''
    if (apiBase.startsWith('http')) {
      return apiBase.replace(/^http/, 'ws').replace(/\/$/, '')
    }
    return (location.protocol === 'https:' ? 'wss://' : 'ws://') + location.host + apiBase
  }
  const openTerminal = (row) => {
    const token = encodeURIComponent(localStorage.getItem('token') || '')
    termName.value = row.name
    termUrl.value = `${buildWsBase()}/server/terminal?id=${row.ID}&token=${token}`
    termVisible.value = true
  }

  // 批量执行命令
  const execVisible = ref(false)
  const execCmdText = ref('')
  const execResults = ref([])
  const execLoading = ref(false)
  const openExec = () => {
    execCmdText.value = ''
    execResults.value = []
    execVisible.value = true
  }
  const handleExec = async () => {
    if (!execCmdText.value) {
      ElMessage.warning('请输入命令')
      return
    }
    const ids = multipleSelection.value.map((i) => i.ID)
    execLoading.value = true
    const res = await execCmd({ ids, cmd: execCmdText.value })
    execLoading.value = false
    if (res.code === 0) {
      execResults.value = res.data || []
    }
  }
</script>

<style></style>
