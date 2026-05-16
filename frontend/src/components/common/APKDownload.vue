<template>
  <div class="apk-download-card">
    <div v-if="loading" class="loading-state">
      <LoadingSpinner />
      <p>Verificando disponibilidade...</p>
    </div>

    <div v-else-if="error" class="error-state">
      <svg class="icon-error" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <circle cx="12" cy="12" r="10" stroke-width="2"/>
        <line x1="12" y1="8" x2="12" y2="12" stroke-width="2"/>
        <line x1="12" y1="16" x2="12.01" y2="16" stroke-width="2"/>
      </svg>
      <p>{{ error }}</p>
    </div>

    <div v-else-if="apkInfo.available" class="apk-available">
      <div class="apk-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path d="M17 17h2a2 2 0 0 0 2-2v-4a2 2 0 0 0-2-2H5a2 2 0 0 0-2 2v4a2 2 0 0 0 2 2h2" stroke-width="2"/>
          <path d="M17 9V7a2 2 0 0 0-2-2H9a2 2 0 0 0-2 2v2" stroke-width="2"/>
          <rect x="7" y="13" width="10" height="8" rx="2" stroke-width="2"/>
        </svg>
      </div>

      <div class="apk-info">
        <h3>App Android</h3>
        <div class="apk-details">
          <div class="detail-item">
            <span class="label">Arquivo:</span>
            <span class="value">{{ apkInfo.filename }}</span>
          </div>
          <div class="detail-item">
            <span class="label">Tamanho:</span>
            <span class="value">{{ formatSize(apkInfo.size) }}</span>
          </div>
          <div class="detail-item">
            <span class="label">Atualizado:</span>
            <span class="value">{{ formatDate(apkInfo.modified) }}</span>
          </div>
        </div>
      </div>

      <button @click="downloadAPK" class="download-button" :disabled="downloading">
        <svg v-if="!downloading" class="icon-download" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" stroke-width="2"/>
          <polyline points="7 10 12 15 17 10" stroke-width="2"/>
          <line x1="12" y1="15" x2="12" y2="3" stroke-width="2"/>
        </svg>
        <LoadingSpinner v-else class="spinner-small" />
        <span>{{ downloading ? 'Baixando...' : 'Baixar APK' }}</span>
      </button>

      <p class="install-note">
        <svg class="icon-info" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <circle cx="12" cy="12" r="10" stroke-width="2"/>
          <line x1="12" y1="16" x2="12" y2="12" stroke-width="2"/>
          <line x1="12" y1="8" x2="12.01" y2="8" stroke-width="2"/>
        </svg>
        Você precisará habilitar "Fontes desconhecidas" nas configurações do Android para instalar o app.
      </p>
    </div>

    <div v-else class="apk-unavailable">
      <svg class="icon-unavailable" viewBox="0 0 24 24" fill="none" stroke="currentColor">
        <circle cx="12" cy="12" r="10" stroke-width="2"/>
        <line x1="15" y1="9" x2="9" y2="15" stroke-width="2"/>
        <line x1="9" y1="9" x2="15" y2="15" stroke-width="2"/>
      </svg>
      <h3>APK não disponível</h3>
      <p>O aplicativo Android não está disponível para download no momento.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '@/services/api'
import LoadingSpinner from './LoadingSpinner.vue'

interface APKInfo {
  available: boolean
  filename?: string
  size?: number
  modified?: string
}

const apkInfo = ref<APKInfo>({ available: false })
const loading = ref(true)
const downloading = ref(false)
const error = ref<string | null>(null)

const fetchAPKInfo = async () => {
  try {
    loading.value = true
    error.value = null
    const response = await api.get<APKInfo>('/apk/info')
    apkInfo.value = response.data
  } catch (err: any) {
    if (err.response?.status === 404) {
      apkInfo.value = { available: false }
    } else {
      error.value = 'Erro ao verificar disponibilidade do APK'
      console.error('Erro ao obter informações do APK:', err)
    }
  } finally {
    loading.value = false
  }
}

const downloadAPK = () => {
  downloading.value = true
  
  // Cria um link temporário para download
  const link = document.createElement('a')
  link.href = `${api.defaults.baseURL}/apk/download`
  link.download = apkInfo.value.filename || 'media-server.apk'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  
  // Reset após um delay
  setTimeout(() => {
    downloading.value = false
  }, 2000)
}

const formatSize = (bytes?: number): string => {
  if (!bytes) return 'N/A'
  const mb = bytes / (1024 * 1024)
  return `${mb.toFixed(2)} MB`
}

const formatDate = (date?: string): string => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleDateString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchAPKInfo()
})
</script>

<style scoped>
.apk-download-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 16px;
  padding: 2rem;
  color: white;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}

.loading-state,
.error-state,
.apk-unavailable {
  text-align: center;
  padding: 2rem 0;
}

.loading-state p,
.error-state p,
.apk-unavailable p {
  margin-top: 1rem;
  opacity: 0.9;
}

.icon-error,
.icon-unavailable {
  width: 64px;
  height: 64px;
  margin: 0 auto;
  opacity: 0.8;
}

.apk-available {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.apk-icon {
  display: flex;
  justify-content: center;
}

.apk-icon svg {
  width: 80px;
  height: 80px;
  filter: drop-shadow(0 4px 8px rgba(0, 0, 0, 0.2));
}

.apk-info {
  text-align: center;
}

.apk-info h3 {
  font-size: 1.75rem;
  font-weight: 700;
  margin-bottom: 1rem;
}

.apk-details {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 1rem;
  backdrop-filter: blur(10px);
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.95rem;
}

.detail-item .label {
  opacity: 0.8;
  font-weight: 500;
}

.detail-item .value {
  font-weight: 600;
}

.download-button {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  background: white;
  color: #667eea;
  border: none;
  border-radius: 12px;
  padding: 1rem 2rem;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.download-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.download-button:active:not(:disabled) {
  transform: translateY(0);
}

.download-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.icon-download {
  width: 24px;
  height: 24px;
}

.spinner-small {
  width: 20px;
  height: 20px;
}

.install-note {
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 0.75rem;
  font-size: 0.875rem;
  line-height: 1.5;
  opacity: 0.9;
}

.icon-info {
  width: 20px;
  height: 20px;
  flex-shrink: 0;
  margin-top: 2px;
}

@media (max-width: 640px) {
  .apk-download-card {
    padding: 1.5rem;
  }

  .apk-icon svg {
    width: 64px;
    height: 64px;
  }

  .apk-info h3 {
    font-size: 1.5rem;
  }

  .download-button {
    padding: 0.875rem 1.5rem;
    font-size: 1rem;
  }
}
</style>
