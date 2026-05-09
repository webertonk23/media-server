<template>
  <DefaultLayout>
    <div class="monitor-page">
      <div class="monitor-header">
        <h1 class="monitor-title">Monitor de Transcodificação</h1>
        <p class="monitor-subtitle">Acompanhe a conversão de mídia em tempo real</p>
      </div>
      
      <div class="status-grid">
        <div class="status-box processing-box">
          <div class="box-header">
            <h3>Processando</h3>
            <span class="badge">{{ status.processing.length }}</span>
          </div>
          <ul v-if="status.processing.length > 0" class="file-list">
            <li v-for="item in status.processing" :key="item.id" class="file-item">
              <span class="file-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="spin-icon">
                  <path d="M21 12a9 9 0 11-6.219-8.56"/>
                </svg>
              </span>
              <span class="file-path">{{ formatPath(item.path) }}</span>
            </li>
          </ul>
          <div v-else class="empty-state">
            <p>Nenhum arquivo em processamento no momento.</p>
          </div>
        </div>
        
        <div class="status-box pending-box">
          <div class="box-header">
            <h3>Fila de Espera</h3>
            <span class="badge">{{ status.pending.length }}</span>
          </div>
          <ul v-if="status.pending.length > 0" class="file-list">
            <li v-for="item in status.pending" :key="item.id" class="file-item">
              <span class="file-icon">⏳</span>
              <span class="file-path">{{ formatPath(item.path) }}</span>
            </li>
          </ul>
          <div v-else class="empty-state">
            <p>A fila de conversão está vazia.</p>
          </div>
        </div>
      </div>
    </div>
  </DefaultLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import apiClient from '@/services/api'
import type { MediaFile } from '@/types/media'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

const status = ref<{ pending: MediaFile[], processing: MediaFile[] }>({ 
  pending: [], 
  processing: [] 
})

let intervalId: number | null = null

const fetchStatus = async () => {
  try {
    const res = await apiClient.get('/transcode/status')
    status.value = {
      pending: res.data?.pending || [],
      processing: res.data?.processing || []
    }
  } catch (e) {
    console.error(e)
  }
}

const formatPath = (path: string) => {
  const parts = path.split('/')
  return parts[parts.length - 1]
}

onMounted(() => {
  fetchStatus()
  intervalId = window.setInterval(fetchStatus, 5000)
})

onUnmounted(() => {
  if (intervalId) clearInterval(intervalId)
})
</script>

<style scoped>
.monitor-page { padding: 2rem; min-height: 100vh; max-width: 1000px; margin: 0 auto; }
.monitor-header { margin-bottom: 2.5rem; }
.monitor-title { font-family: var(--font-display); font-size: 2rem; font-weight: 700; margin-bottom: 0.375rem; }
.monitor-subtitle { font-size: 0.9rem; color: var(--color-text-muted); }

.status-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; }

.status-box {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.07);
  border-radius: 0.875rem;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.box-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.125rem 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  background: rgba(255, 255, 255, 0.02);
}

.box-header h3 { font-size: 1rem; font-weight: 600; color: var(--color-text-primary); margin: 0; }
.badge { background: var(--color-accent); color: #fff; padding: 0.2rem 0.6rem; border-radius: 1rem; font-size: 0.75rem; font-weight: 700; }

.file-list { list-style: none; padding: 0; margin: 0; max-height: 400px; overflow-y: auto; }
.file-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}
.file-item:last-child { border-bottom: none; }

.file-icon { display: flex; align-items: center; justify-content: center; width: 1.5rem; height: 1.5rem; color: var(--color-accent); }
.file-path { font-size: 0.85rem; color: var(--color-text-secondary); word-break: break-all; }

.empty-state { padding: 3rem 1.5rem; text-align: center; color: var(--color-text-muted); font-size: 0.9rem; font-style: italic; }

.spin-icon { width: 1.2rem; height: 1.2rem; animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

@media (max-width: 768px) {
  .status-grid { grid-template-columns: 1fr; }
}
</style>
