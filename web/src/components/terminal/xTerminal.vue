<template>
  <div ref="el" class="xterm-box"></div>
</template>

<script setup>
  import { Terminal } from '@xterm/xterm'
  import { FitAddon } from '@xterm/addon-fit'
  import '@xterm/xterm/css/xterm.css'
  import { onMounted, onBeforeUnmount, ref, watch } from 'vue'

  const props = defineProps({
    url: { type: String, default: '' },
    // interactive=true: 双向终端(SSH)；false: 只读日志流
    interactive: { type: Boolean, default: true }
  })

  const el = ref()
  let term, ws, fit

  const connect = (u) => {
    if (!u) return
    if (ws) { try { ws.close() } catch (e) {} }
    ws = new WebSocket(u)
    ws.onmessage = (e) => term.write(typeof e.data === 'string' ? e.data : '')
    ws.onclose = () => term.write('\r\n\r\n[连接已断开]\r\n')
    ws.onerror = () => term.write('\r\n[连接错误]\r\n')
    if (props.interactive) {
      term.onData((d) => ws && ws.readyState === 1 && ws.send(d))
    }
  }

  onMounted(() => {
    term = new Terminal({ fontSize: 14, cursorBlink: true, scrollback: 5000 })
    fit = new FitAddon()
    term.loadAddon(fit)
    term.open(el.value)
    setTimeout(() => fit && fit.fit(), 60)
    if (props.url) connect(props.url)
  })

  watch(() => props.url, (u) => connect(u))

  onBeforeUnmount(() => {
    if (ws) { try { ws.close() } catch (e) {} }
    if (term) term.dispose()
  })

  defineExpose({ fit: () => fit && fit.fit() })
</script>

<style scoped>
  .xterm-box {
    height: 100%;
    width: 100%;
    background: #000;
    padding: 4px;
  }
  .xterm-box :deep(.xterm) {
    height: 100%;
  }
</style>
