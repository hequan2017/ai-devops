<template>
  <div class="mt-2">
    <el-card>
      <template #header>
        <div class="text-lg font-bold">关于 AI运维平台</div>
      </template>
      <div class="leading-loose text-sm">
        <p>
          <strong>AI运维平台（AI-DevOps）</strong>
          是一款面向物理服务器资产的运维管理系统，聚焦戴尔（Dell）、华为（Huawei）等主流品牌服务器的全生命周期管理。
        </p>
        <p class="mt-3 font-bold">核心能力</p>
        <ul class="list-disc ml-6 mt-1">
          <li>服务器资产台账：厂商、型号、SN、机房机柜、硬件配置、负责人</li>
          <li>统一记录带外管理 IP（iDRAC / iBMC）与业务 IP</li>
          <li>运行状态可视化：运行中 / 已关机 / 维护中 / 故障</li>
          <li>基于角色（RBAC）与菜单的细粒度权限管控</li>
          <li>SQLite 文件数据库，开箱即用，便于容器化部署</li>
        </ul>
        <p class="mt-3 font-bold">技术栈</p>
        <p class="mt-1">
          后端 Go (Gin + GORM)，前端 Vue3 + Element Plus + Pinia + Vite。
        </p>
      </div>
    </el-card>
  </div>
</template>

<script setup>
  defineOptions({
    name: 'About'
  })
</script>
