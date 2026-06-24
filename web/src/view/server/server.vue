<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="elSearchFormRef"
        :inline="true"
        :model="searchInfo"
        class="demo-form-inline"
        @keyup.enter="onSubmit"
      >
        <el-form-item label="关键字" prop="keyword">
          <el-input
            v-model="searchInfo.keyword"
            placeholder="名称/SN/业务IP/管理IP"
            clearable
          />
        </el-form-item>
        <el-form-item label="厂商" prop="manufacturer">
          <el-select
            v-model="searchInfo.manufacturer"
            placeholder="请选择厂商"
            clearable
          >
            <el-option
              v-for="item in manufacturerOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select
            v-model="searchInfo.status"
            placeholder="请选择状态"
            clearable
          >
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
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
        <el-button
          icon="delete"
          style="margin-left: 10px"
          :disabled="!multipleSelection.length"
          @click="onDelete"
        >删除</el-button>
      </div>
      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="创建日期" prop="createdAt" width="170">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column align="left" label="服务器名称" prop="name" min-width="140" show-overflow-tooltip />
        <el-table-column align="left" label="厂商" prop="manufacturer" width="90">
          <template #default="scope">
            <el-tag :type="manufacturerTagType(scope.row.manufacturer)" effect="light">
              {{ manufacturerLabel(scope.row.manufacturer) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="型号" prop="model" min-width="130" show-overflow-tooltip />
        <el-table-column align="left" label="序列号" prop="serialNumber" min-width="150" show-overflow-tooltip />
        <el-table-column align="left" label="业务IP" prop="hostIp" width="130" />
        <el-table-column align="left" label="机房" prop="dataCenter" width="110" show-overflow-tooltip />
        <el-table-column align="left" label="机柜" prop="cabinet" width="90" />
        <el-table-column align="left" label="状态" prop="status" width="100">
          <template #default="scope">
            <el-tag :type="statusTagType(scope.row.status)" effect="dark">
              {{ statusLabel(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column align="left" label="负责人" prop="owner" width="100" />
        <el-table-column align="left" label="操作" fixed="right" min-width="220">
          <template #default="scope">
            <el-button type="primary" link icon="view" @click="getDetails(scope.row)">详情</el-button>
            <el-button type="primary" link icon="edit" @click="updateServerFunc(scope.row)">变更</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <!-- 新增 / 编辑 -->
    <el-drawer
      destroy-on-close
      size="680"
      v-model="dialogFormVisible"
      :show-close="false"
      :before-close="closeDialog"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '新增服务器' : '编辑服务器' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>
      <el-form
        :model="formData"
        label-position="top"
        ref="elFormRef"
        :rules="rule"
      >
        <el-form-item label="服务器名称:" prop="name">
          <el-input v-model="formData.name" clearable placeholder="请输入服务器名称/主机名" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="厂商:" prop="manufacturer">
              <el-select v-model="formData.manufacturer" placeholder="请选择厂商" class="w-full">
                <el-option
                  v-for="item in manufacturerOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态:" prop="status">
              <el-select v-model="formData.status" placeholder="请选择状态" class="w-full">
                <el-option
                  v-for="item in statusOptions"
                  :key="item.value"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="型号:" prop="model">
          <el-input v-model="formData.model" clearable placeholder="如 PowerEdge R740 / 2288H V5" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="序列号(SN):" prop="serialNumber">
              <el-input v-model="formData.serialNumber" clearable placeholder="请输入序列号" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="业务IP:" prop="hostIp">
              <el-input v-model="formData.hostIp" clearable placeholder="请输入业务IP" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="带外管理IP:" prop="manageIp">
          <el-input v-model="formData.manageIp" clearable placeholder="iDRAC / iBMC 管理IP" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="机房:" prop="dataCenter">
              <el-input v-model="formData.dataCenter" clearable placeholder="机房" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="机柜:" prop="cabinet">
              <el-input v-model="formData.cabinet" clearable placeholder="机柜" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="U位:" prop="uPosition">
              <el-input v-model="formData.uPosition" clearable placeholder="U位" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="CPU:" prop="cpu">
              <el-input v-model="formData.cpu" clearable placeholder="如 2*Intel Xeon Gold" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="内存:" prop="memory">
              <el-input v-model="formData.memory" clearable placeholder="如 128GB" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="磁盘:" prop="disk">
              <el-input v-model="formData.disk" clearable placeholder="如 2*960GB SSD" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="操作系统:" prop="os">
              <el-input v-model="formData.os" clearable placeholder="如 CentOS 7.9" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="负责人:" prop="owner">
              <el-input v-model="formData.owner" clearable placeholder="负责人" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="备注:" prop="remark">
          <el-input v-model="formData.remark" type="textarea" :rows="3" clearable placeholder="备注" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <!-- 详情 -->
    <el-drawer
      destroy-on-close
      size="560"
      v-model="detailShow"
      :show-close="true"
      :before-close="closeDetailShow"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item label="服务器名称">{{ detailForm.name }}</el-descriptions-item>
        <el-descriptions-item label="厂商">{{ manufacturerLabel(detailForm.manufacturer) }}</el-descriptions-item>
        <el-descriptions-item label="型号">{{ detailForm.model }}</el-descriptions-item>
        <el-descriptions-item label="序列号(SN)">{{ detailForm.serialNumber }}</el-descriptions-item>
        <el-descriptions-item label="业务IP">{{ detailForm.hostIp }}</el-descriptions-item>
        <el-descriptions-item label="带外管理IP">{{ detailForm.manageIp }}</el-descriptions-item>
        <el-descriptions-item label="机房">{{ detailForm.dataCenter }}</el-descriptions-item>
        <el-descriptions-item label="机柜">{{ detailForm.cabinet }}</el-descriptions-item>
        <el-descriptions-item label="U位">{{ detailForm.uPosition }}</el-descriptions-item>
        <el-descriptions-item label="CPU">{{ detailForm.cpu }}</el-descriptions-item>
        <el-descriptions-item label="内存">{{ detailForm.memory }}</el-descriptions-item>
        <el-descriptions-item label="磁盘">{{ detailForm.disk }}</el-descriptions-item>
        <el-descriptions-item label="操作系统">{{ detailForm.os }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="statusTagType(detailForm.status)" effect="dark">{{ statusLabel(detailForm.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="负责人">{{ detailForm.owner }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ detailForm.remark }}</el-descriptions-item>
      </el-descriptions>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    createServer,
    deleteServer,
    deleteServerByIds,
    updateServer,
    findServer,
    getServerList
  } from '@/api/server'

  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref, reactive } from 'vue'

  defineOptions({ name: 'Server' })

  // 厂商选项
  const manufacturerOptions = [
    { value: 'dell', label: '戴尔(Dell)' },
    { value: 'huawei', label: '华为(Huawei)' }
  ]
  // 状态选项
  const statusOptions = [
    { value: 'online', label: '运行中' },
    { value: 'offline', label: '已关机' },
    { value: 'maintenance', label: '维护中' },
    { value: 'fault', label: '故障' }
  ]

  const manufacturerLabel = (v) => {
    const item = manufacturerOptions.find((i) => i.value === v)
    return item ? item.label : v || '-'
  }
  const manufacturerTagType = (v) => (v === 'huawei' ? 'danger' : 'primary')

  const statusLabel = (v) => {
    const item = statusOptions.find((i) => i.value === v)
    return item ? item.label : v || '-'
  }
  const statusTagType = (v) => {
    switch (v) {
      case 'online': return 'success'
      case 'offline': return 'info'
      case 'maintenance': return 'warning'
      case 'fault': return 'danger'
      default: return 'info'
    }
  }

  // 表单数据
  const defaultForm = () => ({
    name: '', manufacturer: '', model: '', serialNumber: '',
    manageIp: '', hostIp: '', dataCenter: '', cabinet: '',
    uPosition: '', cpu: '', memory: '', disk: '', os: '',
    status: 'online', owner: '', remark: ''
  })
  const formData = ref(defaultForm())

  // 验证规则
  const rule = reactive({
    name: [{ required: true, message: '请输入服务器名称', trigger: ['input', 'blur'] }],
    manufacturer: [{ required: true, message: '请选择厂商', trigger: 'change' }],
    status: [{ required: true, message: '请选择状态', trigger: 'change' }]
  })

  const elFormRef = ref()
  const elSearchFormRef = ref()

  // =========== 表格控制 ===========
  const page = ref(1)
  const total = ref(0)
  const pageSize = ref(10)
  const tableData = ref([])
  const searchInfo = ref({})

  const onReset = () => {
    searchInfo.value = {}
    getTableData()
  }

  const onSubmit = () => {
    page.value = 1
    getTableData()
  }

  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  const getTableData = async () => {
    const table = await getServerList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (table.code === 0) {
      tableData.value = table.data.list
      total.value = table.data.total
      page.value = table.data.page
      pageSize.value = table.data.pageSize
    }
  }

  getTableData()

  // 多选
  const multipleSelection = ref([])
  const handleSelectionChange = (val) => {
    multipleSelection.value = val
  }

  // 删除行
  const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除该服务器吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const res = await deleteServer({ ID: row.ID })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '删除成功' })
        if (tableData.value.length === 1 && page.value > 1) page.value--
        getTableData()
      }
    })
  }

  // 批量删除
  const onDelete = () => {
    ElMessageBox.confirm('确定要删除选中的服务器吗?', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }).then(async () => {
      const ids = multipleSelection.value.map((item) => item.ID)
      const res = await deleteServerByIds({ ids })
      if (res.code === 0) {
        ElMessage({ type: 'success', message: '删除成功' })
        if (tableData.value.length === ids.length && page.value > 1) page.value--
        getTableData()
      }
    })
  }

  // 新增 / 编辑控制
  const type = ref('')
  const dialogFormVisible = ref(false)

  const openDialog = () => {
    type.value = 'create'
    formData.value = defaultForm()
    dialogFormVisible.value = true
  }

  const updateServerFunc = async (row) => {
    const res = await findServer({ ID: row.ID })
    type.value = 'update'
    if (res.code === 0) {
      formData.value = res.data
      dialogFormVisible.value = true
    }
  }

  const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = defaultForm()
  }

  const enterDialog = () => {
    elFormRef.value?.validate(async (valid) => {
      if (!valid) return
      const res = type.value === 'update'
        ? await updateServer(formData.value)
        : await createServer(formData.value)
      if (res.code === 0) {
        ElMessage({ type: 'success', message: type.value === 'update' ? '更新成功' : '创建成功' })
        closeDialog()
        getTableData()
      }
    })
  }

  // 详情
  const detailShow = ref(false)
  const detailForm = ref({})

  const getDetails = async (row) => {
    const res = await findServer({ ID: row.ID })
    if (res.code === 0) {
      detailForm.value = res.data
      detailShow.value = true
    }
  }

  const closeDetailShow = () => {
    detailShow.value = false
    detailForm.value = {}
  }
</script>

<style></style>
