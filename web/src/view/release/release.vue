<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="项目"><el-input v-model="searchInfo.project" clearable /></el-form-item>
        <el-form-item label="环境">
          <el-select v-model="searchInfo.environment" clearable style="width:120px">
            <el-option label="dev" value="dev" />
            <el-option label="staging" value="staging" />
            <el-option label="prod" value="prod" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable style="width:130px">
            <el-option v-for="s in statusOptions" :key="s.value" :label="s.label" :value="s.value" />
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
        <el-button type="primary" icon="plus" @click="openDialog">发起发版</el-button>
      </div>
      <el-table :data="tableData" row-key="ID">
        <el-table-column label="项目" prop="project" min-width="120" show-overflow-tooltip />
        <el-table-column label="版本" prop="version" width="110" />
        <el-table-column label="分支" prop="branch" width="120" show-overflow-tooltip />
        <el-table-column label="环境" prop="environment" width="90">
          <template #default="scope">
            <el-tag size="small">{{ scope.row.environment }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="scope">
            <el-tag :type="statusTag(scope.row.status)" effect="dark">{{ statusLabel(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="170">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" min-width="320">
          <template #default="scope">
            <el-button v-if="scope.row.status==='draft'" type="primary" link @click="submit(scope.row)">提交审批</el-button>
            <el-button v-if="scope.row.status==='pending'" type="success" link @click="approve(scope.row)">通过</el-button>
            <el-button v-if="scope.row.status==='pending'" type="warning" link @click="reject(scope.row)">拒绝</el-button>
            <el-button v-if="scope.row.status==='approved'" type="primary" link @click="execute(scope.row)">执行发布</el-button>
            <el-button type="primary" link @click="viewResult(scope.row)">结果</el-button>
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

    <el-drawer destroy-on-close size="640" v-model="dialogVisible" :show-close="false" :before-close="closeDialog">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">{{ type === 'create' ? '发起发版' : '编辑发版' }}</span>
          <div>
            <el-button type="primary" @click="enterDialog">确 定</el-button>
            <el-button @click="closeDialog">取 消</el-button>
          </div>
        </div>
      </template>
      <el-form :model="formData" label-position="top">
        <el-row :gutter="20">
          <el-col :span="12"><el-form-item label="项目"><el-input v-model="formData.project" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="版本号"><el-input v-model="formData.version" placeholder="v1.0.0" /></el-form-item></el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12"><el-form-item label="分支"><el-input v-model="formData.branch" placeholder="main" /></el-form-item></el-col>
          <el-col :span="12">
            <el-form-item label="环境">
              <el-select v-model="formData.environment" class="w-full">
                <el-option label="dev" value="dev" />
                <el-option label="staging" value="staging" />
                <el-option label="prod" value="prod" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="发版说明"><el-input v-model="formData.description" type="textarea" :rows="2" /></el-form-item>
        <el-form-item label="发布脚本(执行阶段运行)">
          <el-input v-model="formData.script" type="textarea" :rows="5" placeholder="# 例: git pull && go build -o app . && systemctl restart app" />
        </el-form-item>
      </el-form>
    </el-drawer>

    <el-drawer destroy-on-close size="720" v-model="resultVisible" title="执行结果">
      <pre style="white-space:pre-wrap;font-size:12px;max-height:70vh;overflow:auto">{{ resultContent }}</pre>
    </el-drawer>
  </div>
</template>

<script setup>
  import {
    getReleaseList, createRelease, updateRelease, deleteRelease, findRelease,
    submitRelease, approveRelease, rejectRelease, executeRelease
  } from '@/api/release'
  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref } from 'vue'

  defineOptions({ name: 'Release' })

  const statusOptions = [
    { value: 'draft', label: '草稿' },
    { value: 'pending', label: '待审批' },
    { value: 'approved', label: '已通过' },
    { value: 'rejected', label: '已拒绝' },
    { value: 'executing', label: '执行中' },
    { value: 'done', label: '已完成' },
    { value: 'failed', label: '失败' }
  ]
  const statusLabel = (v) => statusOptions.find((s) => s.value === v)?.label || v
  const statusTag = (v) => ({
    draft: 'info', pending: 'warning', approved: 'success', rejected: 'danger',
    executing: '', done: 'success', failed: 'danger'
  }[v] || 'info')

  const page = ref(1); const pageSize = ref(10); const total = ref(0)
  const tableData = ref([]); const searchInfo = ref({})

  const getTableData = async () => {
    const res = await getReleaseList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (res.code === 0) { tableData.value = res.data.list; total.value = res.data.total }
  }
  getTableData()

  const onSubmit = () => { page.value = 1; getTableData() }
  const onReset = () => { searchInfo.value = {}; getTableData() }

  const submit = async (row) => { const r = await submitRelease({ ID: row.ID }); if (r.code === 0) { ElMessage.success('已提交审批'); getTableData() } }
  const approve = async (row) => { const r = await approveRelease({ ID: row.ID }); if (r.code === 0) { ElMessage.success('已通过'); getTableData() } }
  const reject = (row) => {
    ElMessageBox.prompt('请输入拒绝原因', '拒绝', { type: 'warning' }).then(async ({ value }) => {
      const r = await rejectRelease({ ID: row.ID, reason: value })
      if (r.code === 0) { ElMessage.success('已拒绝'); getTableData() }
    })
  }
  const execute = (row) => {
    ElMessageBox.confirm('确定执行发布? 将运行脚本。', '确认', { type: 'warning' }).then(async () => {
      const r = await executeRelease({ ID: row.ID })
      if (r.code === 0) { ElMessage.success('执行完成'); getTableData() }
    })
  }

  const resultVisible = ref(false); const resultContent = ref('')
  const viewResult = (row) => { resultContent.value = row.result || '(无)'; resultVisible.value = true }

  const type = ref(''); const dialogVisible = ref(false)
  const defaultForm = () => ({ project: '', version: '', branch: 'main', environment: 'prod', description: '', script: '' })
  const formData = ref(defaultForm())
  const openDialog = () => { type.value = 'create'; formData.value = defaultForm(); dialogVisible.value = true }
  const editRow = async (row) => { const r = await findRelease({ ID: row.ID }); if (r.code === 0) { formData.value = r.data; type.value = 'update'; dialogVisible.value = true } }
  const closeDialog = () => { dialogVisible.value = false; formData.value = defaultForm() }
  const enterDialog = async () => {
    const res = type.value === 'update' ? await updateRelease(formData.value) : await createRelease(formData.value)
    if (res.code === 0) { ElMessage.success('保存成功'); closeDialog(); getTableData() }
  }
  const deleteRow = (row) => {
    ElMessageBox.confirm('确定删除该发版记录?', '提示', { type: 'warning' }).then(async () => {
      const r = await deleteRelease({ ID: row.ID })
      if (r.code === 0) { ElMessage.success('删除成功'); getTableData() }
    })
  }
</script>
