<template>
  <DefaultLayout>
    <div class="logs-page">
      <div class="logs-header">
        <h1 class="logs-title">Logs do Sistema</h1>
        <div class="logs-actions">
          <button class="btn btn-ghost btn-sm" @click="fetchLogs">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="refresh-icon" :class="{ 'spinning': loading }">
              <path d="M23 4v6h-6M1 20v-6h6M3.51 9a9 9 0 0114.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0020.49 15" />
            </svg>
            Atualizar
          </button>
          <button class="btn btn-ghost btn-sm" @click="clearDisplay">Limpar Vista</button>
        </div>
      </div>

      <div class="logs-container" ref="logsContainer">
        <div v-if="logs.length === 0" class="logs-empty">
          Nenhum log disponível.
        </div>
        <div v-for="(log, index) in logs" :key="index" class="log-entry">
          <span class="log-time">{{ formatTime(log.timestamp) }}</span>
          <span class="log-message">{{ log.message }}</span>
        </div>
      </div>
    </div>
  </DefaultLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
import apiClient from '@/services/api'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

interface LogEntry {
  timestamp: string
  message: string
}

const logs = ref<LogEntry[]>([])
const loading = ref(false)
const logsContainer = ref<HTMLElement | null>(null)
let intervalId: number | null = null

const fetchLogs = async () => {
  loading.value = true
  try {
    const res = await apiClient.get('/logs')
    logs.value = res.data || []
    await nextTick()
    scrollToBottom()
  } catch (e) {
    console.error('Falha ao buscar logs', e)
  } finally {
    loading.value = false
  }
}

const formatTime = (ts: string) => {
  const date = new Date(ts)
  return date.toLocaleTimeString('pt-BR', { hour12: false }) + '.' + String(date.getMilliseconds()).padStart(3, '0')
}

const scrollToBottom = () => {
  if (logsContainer.value) {
    logsContainer.value.scrollTop = logsContainer.value.scrollHeight
  }
}

const clearDisplay = () => {
  logs.value = []
}

onMounted(() => {
  fetchLogs()
  intervalId = window.setInterval(fetchLogs, 3000)
})

onUnmounted(() => {
  if (intervalId) clearInterval(intervalId)
})
</script>

<style scoped>
.logs-page {
  padding: 2rem;
  height: calc(100vh - 64px);
  display: flex;
  flex-direction: column;
}

.logs-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.logs-title {
  font-family: var(--font-display);
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0;
}

.logs-actions {
  display: flex;
  gap: 0.75rem;
}

.logs-container {
  flex: 1;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.75rem;
  padding: 1rem;
  overflow-y: auto;
  font-family: 'Fira Code', 'Courier New', Courier, monospace;
  font-size: 0.85rem;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  scroll-behavior: smooth;
}

.log-entry {
  display: flex;
  gap: 1rem;
  line-height: 1.4;
  padding: 0.125rem 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
}

.log-time {
  color: var(--color-cinema-accent-secondary);
  flex-shrink: 0;
  white-space: nowrap;
}

.log-message {
  color: #e2e8f0;
  word-break: break-all;
  white-space: pre-wrap;
}

.logs-empty {
  color: var(--color-text-muted);
  text-align: center;
  padding: 3rem;
  font-style: italic;
}

.refresh-icon {
  width: 1rem;
  height: 1rem;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Custom scrollbar for logs */
.logs-container::-webkit-scrollbar {
  width: 8px;
}

.logs-container::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.02);
}

.logs-container::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.logs-container::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>
