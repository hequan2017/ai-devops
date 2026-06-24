<template>
  <div>
    <el-tabs v-model="active">
      <!-- 告警规则 -->
      <el-tab-pane label="告警规则" name="rule">
        <div class="gva-table-box">
          <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openRule">新增规则</el-button>
          </div>
          <el-table :data="rules">
            <el-table-column label="名称" prop="name" min-width="140" show-overflow-tooltip />
            <el-table-column label="指标" prop="metric" width="90" />
            <el-table-column label="条件" width="140">
              <template #default="scope">{{ scope.row.operator }} {{ scope.row.threshold }}%</template>
            </el-table-column>
            <el-table-column label="级别" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.level === 'critical' ? 'danger' : 'warning'">{{ scope.row.level }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="启用" width="80">
              <template #default="scope">
                <el-tag :type="scope.row.enabled ? 'success' : 'info'">{{ scope.row.enabled ? '是' : '否' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" fixed="right" width="160">
              <template #default="scope">
                <el-button type="primary" link icon="edit" @click="editRule(scope.row)">编辑</el-button>
                <el-button type="primary" link icon="delete" @click="delRule(scope.row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination layout="total, prev, pager, next" :current-page="rPage" :page-size="20" :total="rTotal"
              @current-change="(v)=>{rPage=v;loadRules()}" />
          </div>
        </div>
      </el-tab-pane>

      <!-- 告警记录 -->
      <el-tab-pane label="告警记录" name="record">
        <div class="gva-search-box">
          <el-form :inline="true" :model="recSearch">
            <el-form-item label="级别">
              <el-select v-model="recSearch.level" clearable style="width:130px">
                <el-option label="warn" value="warn" />
                <el-option label="critical" value="critical" />
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" icon="search" @click="loadRecords">查询</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div class="gva-table-box">
          <el-table :data="records">
            <el-table-column label="时间" width="170">
              <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
            </el-table-column>
            <el-table-column label="服务器" prop="serverName" min-width="120" />
            <el-table-column label="规则" prop="ruleName" min-width="140" show-overflow-tooltip />
            <el-table-column label="信息" prop="message" min-width="220" show-overflow-tooltip />
            <el-table-column label="级别" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.level === 'critical' ? 'danger' : 'warning'">{{ scope.row.level }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.resolved ? 'success' : 'danger'">{{ scope.row.resolved ? '已处理' : '未处理' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" fixed="right" width="100">
              <template #default="scope">
                <el-button v-if="!scope.row.resolved" type="primary" link @click="resolve(scope.row)">处理</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="gva-pagination">
            <el-pagination layout="total, prev, pager, next" :current-page="dPage" :page-size="20" :total="dTotal"
              @current-change="(v)=>{dPage=v;loadRecords()}" />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 规则编辑 -->
    <el-drawer v-model="ruleVisible" size="520" :title="ruleType==='create'?'新增规则':'编辑规则'" destroy-on-close :before-close="()=>ruleVisible=false">
      <el-form :model="ruleForm" label-position="top">
        <el-form-item label="规则名称"><el-input v-model="ruleForm.name" /></el-form-item>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="指标">
              <el-select v-model="ruleForm.metric" class="w-full">
                <el-option label="CPU" value="cpu" />
                <el-option label="内存" value="mem" />
                <el-option label="磁盘" value="disk" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="操作符">
              <el-select v-model="ruleForm.operator" class="w-full">
                <el-option label=">" value=">" />
                <el-option label=">=" value=">=" />
                <el-option label="<" value="<" />
                <el-option label="<=" value="<=" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="阈值(%)"><el-input-number v-model="ruleForm.threshold" :min="0" :max="100" class="w-full" /></el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="级别">
              <el-select v-model="ruleForm.level" class="w-full">
                <el-option label="warn" value="warn" />
                <el-option label="critical" value="critical" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="启用"><el-switch v-model="ruleForm.enabled" /></el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="通知Webhook"><el-input v-model="ruleForm.notifyWebhook" placeholder="钉钉/企微/飞书 webhook 地址（留空则不通知）" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="ruleForm.remark" type="textarea" :rows="2" /></el-form-item>
      </el-form>
      <div class="mt-3">
        <el-button type="primary" @click="saveRule">保存</el-button>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
  import { getAlertRuleList, createAlertRule, updateAlertRule, deleteAlertRule, findAlertRule, getAlertRecordList, resolveAlertRecord } from '@/api/alert'
  import { formatDate } from '@/utils/format'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { ref } from 'vue'

  defineOptions({ name: 'Alert' })

  const active = ref('rule')

  // 规则
  const rules = ref([]); const rPage = ref(1); const rTotal = ref(0)
  const loadRules = async () => {
    const res = await getAlertRuleList({ page: rPage.value, pageSize: 20 })
    if (res.code === 0) { rules.value = res.data.list; rTotal.value = res.data.total }
  }
  loadRules()

  const ruleVisible = ref(false); const ruleType = ref('create')
  const ruleForm = ref({ name: '', metric: 'cpu', operator: '>', threshold: 90, level: 'warn', enabled: true, notifyWebhook: '', remark: '' })
  const openRule = () => { ruleType.value = 'create'; ruleForm.value = { name: '', metric: 'cpu', operator: '>', threshold: 90, level: 'warn', enabled: true, notifyWebhook: '', remark: '' }; ruleVisible.value = true }
  const editRule = async (row) => { const r = await findAlertRule({ ID: row.ID }); if (r.code === 0) { ruleForm.value = r.data; ruleType.value = 'update'; ruleVisible.value = true } }
  const saveRule = async () => {
    const res = ruleType.value === 'update' ? await updateAlertRule(ruleForm.value) : await createAlertRule(ruleForm.value)
    if (res.code === 0) { ElMessage.success('保存成功'); ruleVisible.value = false; loadRules() }
  }
  const delRule = (row) => {
    ElMessageBox.confirm('确定删除该规则?', '提示', { type: 'warning' }).then(async () => {
      const r = await deleteAlertRule({ ID: row.ID })
      if (r.code === 0) { ElMessage.success('删除成功'); loadRules() }
    })
  }

  // 记录
  const records = ref([]); const dPage = ref(1); const dTotal = ref(0); const recSearch = ref({})
  const loadRecords = async () => {
    const res = await getAlertRecordList({ page: dPage.value, pageSize: 20, ...recSearch.value })
    if (res.code === 0) { records.value = res.data.list; dTotal.value = res.data.total }
  }
  loadRecords()
  const resolve = (row) => {
    ElMessageBox.confirm('标记为已处理?', '提示', { type: 'warning' }).then(async () => {
      const r = await resolveAlertRecord({ ID: row.ID })
      if (r.code === 0) { ElMessage.success('已处理'); loadRecords() }
    })
  }
</script>
