<template>
  <div class="h-full">
    <div class="mb-2 text-xs tracking-wide text-black/55 dark:text-white/55">运维概览</div>
    <div class="grid grid-cols-2 gap-2 sm:grid-cols-3">
      <button
        v-for="item in shortcuts"
        :key="item.path"
        class="group flex w-full items-center gap-3 rounded-lg border border-black/10 bg-white/70 p-3 text-left transition-all duration-200 hover:border-[var(--el-color-primary)] hover:shadow-sm dark:border-white/10 dark:bg-white/[0.02]"
        type="button"
        @click="toPath(item)"
      >
        <span class="flex h-9 w-9 shrink-0 items-center justify-center rounded-md bg-slate-100 text-slate-700 transition-colors group-hover:bg-[var(--el-color-primary)] group-hover:text-white dark:bg-slate-800 dark:text-slate-200">
          <el-icon><component :is="item.icon" /></el-icon>
        </span>
        <span class="min-w-0">
          <span class="block text-lg font-bold leading-tight" :class="item.warn ? 'text-red-500' : ''">{{ item.count }}</span>
          <span class="block text-xs text-black/60 dark:text-white/60">{{ item.title }}</span>
        </span>
      </button>
    </div>
  </div>
</template>

<script setup>
  import { Monitor, Box, Share, Promotion, WarnTriangleFilled } from '@element-plus/icons-vue'
  import { useRouter } from 'vue-router'
  import { ref, computed, onMounted } from 'vue'
  import { getOverview } from '@/api/server'

  const router = useRouter()
  const toPath = (item) => router.push({ name: item.path })

  const ov = ref({})
  onMounted(async () => {
    const r = await getOverview()
    if (r.code === 0) ov.value = r.data
  })

  const shortcuts = computed(() => [
    { icon: Monitor, title: '服务器', path: 'server', count: ov.value.serverTotal ?? '-' },
    { icon: Box, title: 'Docker', path: 'docker', count: ov.value.dockerTotal ?? '-' },
    { icon: Share, title: 'K8s集群', path: 'k8s', count: ov.value.k8sTotal ?? '-' },
    { icon: Promotion, title: '待审批', path: 'release', count: ov.value.releasePending ?? '-' },
    { icon: WarnTriangleFilled, title: '告警未处理', path: 'alert', count: ov.value.alertUnresolved ?? '-', warn: true }
  ])
</script>

<style scoped lang="scss"></style>
