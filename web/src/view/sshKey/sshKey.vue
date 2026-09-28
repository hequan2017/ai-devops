<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" @keyup.enter="onSubmit">
        <el-form-item label="名称"><el-input v-model="searchInfo.keyword" clearable /></el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog">新增密钥</el-button>
      </div>
      <el-table :data="tableData" row-key="ID">
        <el-table-column label="名称" prop="name" min-width="160" show-overflow-tooltip />
        <el-table-column label="私钥" min-width="200">
          <template #default="scope">
            <span v-if="!scope.row._show">••••••（已隐藏）</span>
            <pre v-else style="white-space:pre-wrap;font-size:11px;max-height:120px;overflow:auto">{{ maskKey(scope.row.privateKey) }}</pre>
          </template>
        </el-table-column>
        <el-table-column label="备注" prop="remark" min-width="140" show-overflow-tooltip />
        <el-table-column label="创建时间" width="170">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" width="220">
          <template #default="scope">
            <el-button type="primary" link @click="scope.row._show = !scope.row._show">{{ scope.row._show ? '隐藏' : '查看' }}</el-button>
            <el-button type="primary" link icon="edit" @click="editRow(scope.row)">编辑</el-button>
            <el-button type="primary" link icon="delete" @click="delRow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination layout="total, prev, pager, next" :current-page="page" :page-size="20" :total="total"
          @current-change="(v)=>{page=v;getTableData()}" />
      </div>
    </div>

    <el-drawer v-model="dialogVisible" size="560" :title="type==='create'?'新增密钥':'编辑密钥'" destroy-on-close :before-close="()=>dialogVisible=false">
      <el-form :model="formData" label-position="top">
        <el-form-item label="名称"><el-input v-model="formData.name" /></el-form-item>
        <el-form-item label="私钥(PEM)"><el-input v-model="formData.privateKey" type="textarea" :rows="8" placeholder="-----BEGIN OPENSSH PRIVATE KEY (REDACTED)-----..." /></el-form-item>
        <el-form-item label="备注"><el-input v-model="formData.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <div class="mt-3"><el-button type="primary" @click="enterDialog">保存</el-button></div>
    </el-drawer>
  </div>
</template>

<script setup>
  import { createSshKey, updateSshKey, deleteSshKey, findSshKey, getSshKeyList } from '@/api/sshKey'
  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref } from 'vue'

  defineOptions({ name: 'SshKey' })

  const page = ref(1); const total = ref(0); const tableData = ref([]); const searchInfo = ref({})
  const getTableData = async () => {
    const res = await getSshKeyList({ page: page.value, pageSize: 20, ...searchInfo.value })
    if (res.code === 0) { tableData.value = (res.data.list || []).map((i) => ({ ...i, _show: false })); total.value = res.data.total }
  }
  getTableData()
  const onSubmit = () => { page.value = 1; getTableData() }
  const onReset = () => { searchInfo.value = {}; getTableData() }

  const maskKey = (k) => (k || '').slice(0, 60) + ((k || '').length > 60 ? '\n...' : '')

  const type = ref(''); const dialogVisible = ref(false)
  const formData = ref({ name: '', privateKey: '', remark: '' })
  const openDialog = () => { type.value = 'create'; formData.value = { name: '', privateKey: '', remark: '' }; dialogVisible.value = true }
  const editRow = async (row) => { const r = await findSshKey({ ID: row.ID }); if (r.code === 0) { formData.value = r.data; type.value = 'update'; dialogVisible.value = true } }
  const enterDialog = async () => {
    const res = type.value === 'update' ? await updateSshKey(formData.value) : await createSshKey(formData.value)
    if (res.code === 0) { ElMessage.success('保存成功'); dialogVisible.value = false; getTableData() }
  }
  const delRow = (row) => {
    ElMessageBox.confirm('确定删除该密�?', '提示', { type: 'warning' }).then(async () => {
      const r = await deleteSshKey({ ID: row.ID })
      if (r.code === 0) { ElMessage.success('删除成功'); getTableData() }
    })
  }
</script>
