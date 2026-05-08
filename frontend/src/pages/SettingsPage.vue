<template>
  <DefaultLayout>
    <div class="settings-page">
      <div class="settings-header">
        <h1 class="settings-title">Configurações</h1>
        <p class="settings-subtitle">Gerencie sua experiência no MediaServer</p>
      </div>

      <div class="settings-sections">
        <!-- Library -->
        <div class="settings-section">
          <h2 class="section-title">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
            </svg>
            Biblioteca
          </h2>
          <div class="settings-items">
            <div class="settings-item">
              <div class="item-info">
                <p class="item-label">Escanear biblioteca completa</p>
                <p class="item-desc">Varre todos os arquivos de mídia nas pastas configuradas</p>
              </div>
              <button class="btn-primary" :disabled="scanning" @click="scanAll">
                <svg v-if="scanning" class="spin-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M21 12a9 9 0 11-6.219-8.56"/>
                </svg>
                <span>{{ scanning ? 'Escaneando...' : 'Escanear tudo' }}</span>
              </button>
            </div>
            <div class="settings-item">
              <div class="item-info">
                <p class="item-label">Escanear apenas filmes</p>
                <p class="item-desc">Processa somente arquivos de filmes</p>
              </div>
              <button class="btn-ghost" :disabled="scanning" @click="scanMovies">Escanear filmes</button>
            </div>
            <div class="settings-item">
              <div class="item-info">
                <p class="item-label">Escanear apenas séries</p>
                <p class="item-desc">Processa somente arquivos de séries</p>
              </div>
              <button class="btn-ghost" :disabled="scanning" @click="scanSeries">Escanear séries</button>
            </div>
          </div>
        </div>

        <!-- About -->
        <div class="settings-section">
          <h2 class="section-title">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            Sobre
          </h2>
          <div class="settings-items">
            <div class="settings-item about-item">
              <div class="about-logo">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M8 5v14l11-7z"/>
                </svg>
              </div>
              <div class="about-info">
                <p class="about-name">MediaServer</p>
                <p class="about-desc">Servidor de mídia local de alta performance</p>
                <div class="about-features">
                  <span>Go backend</span>
                  <span>Vue 3</span>
                  <span>SQLite</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Scan notification -->
      <div v-if="scanMsg" class="scan-toast" :class="scanMsgType">
        {{ scanMsg }}
      </div>
    </div>
  </DefaultLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import apiClient from '@/services/api'
import DefaultLayout from '@/layouts/DefaultLayout.vue'

const scanning = ref(false)
const scanMsg = ref('')
const scanMsgType = ref<'success' | 'error'>('success')

const notify = (msg: string, type: 'success' | 'error' = 'success') => {
  scanMsg.value = msg
  scanMsgType.value = type
  setTimeout(() => { scanMsg.value = '' }, 3500)
}

const runScan = async (endpoint: string, label: string) => {
  if (scanning.value) return
  scanning.value = true
  try {
    await apiClient.post(endpoint)
    notify(`${label} concluído com sucesso!`)
  } catch (e: any) {
    notify(e.message || `Falha ao ${label.toLowerCase()}`, 'error')
  } finally {
    scanning.value = false
  }
}

const scanAll = () => runScan('/scan', 'Scan completo')
const scanMovies = () => runScan('/scan/movies', 'Scan de filmes')
const scanSeries = () => runScan('/scan/series', 'Scan de séries')
</script>

<style scoped>
.settings-page { padding: 2rem; min-height: 100vh; max-width: 860px; }
.settings-header { margin-bottom: 2.5rem; }
.settings-title { font-family: var(--font-display); font-size: 2rem; font-weight: 700; margin-bottom: 0.375rem; }
.settings-subtitle { font-size: 0.9rem; color: var(--color-text-muted); }

.settings-sections { display: flex; flex-direction: column; gap: 2rem; }

.settings-section {
  background: rgba(255,255,255,0.03);
  border: 1px solid rgba(255,255,255,0.07);
  border-radius: 0.875rem;
  overflow: hidden;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  font-size: 0.9375rem;
  font-weight: 600;
  color: var(--color-text-primary);
  padding: 1.125rem 1.5rem;
  border-bottom: 1px solid rgba(255,255,255,0.06);
  background: rgba(255,255,255,0.02);
}
.section-title svg { width: 1.125rem; height: 1.125rem; color: var(--color-accent); }

.settings-items { display: flex; flex-direction: column; }

.settings-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 2rem;
  padding: 1.125rem 1.5rem;
  border-bottom: 1px solid rgba(255,255,255,0.04);
}
.settings-item:last-child { border-bottom: none; }

.item-info { flex: 1; min-width: 0; }
.item-label { font-size: 0.9rem; font-weight: 500; color: var(--color-text-primary); margin-bottom: 0.2rem; }
.item-desc { font-size: 0.8rem; color: var(--color-text-muted); line-height: 1.4; }

.spin-icon { width: 1rem; height: 1rem; animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

/* About item */
.about-item { gap: 1.25rem; }
.about-logo {
  width: 3rem; height: 3rem;
  background: var(--color-accent);
  border-radius: 0.625rem;
  display: flex; align-items: center; justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 4px 16px var(--color-accent-glow);
}
.about-logo svg { width: 1.5rem; height: 1.5rem; color: #fff; }
.about-name { font-size: 1rem; font-weight: 700; color: var(--color-text-primary); margin-bottom: 0.2rem; }
.about-desc { font-size: 0.8rem; color: var(--color-text-muted); margin-bottom: 0.5rem; }
.about-features { display: flex; gap: 0.5rem; flex-wrap: wrap; }
.about-features span {
  font-size: 0.7rem; padding: 0.15rem 0.5rem;
  background: rgba(255,255,255,0.07);
  border-radius: 0.25rem;
  color: var(--color-text-secondary);
}

/* Toast notification */
.scan-toast {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  padding: 0.875rem 1.5rem;
  border-radius: 0.625rem;
  font-size: 0.875rem;
  font-weight: 500;
  backdrop-filter: blur(12px);
  z-index: 100;
  animation: slideIn 0.3s ease;
}
.scan-toast.success { background: rgba(34, 197, 94, 0.15); border: 1px solid rgba(34, 197, 94, 0.3); color: #4ade80; }
.scan-toast.error { background: rgba(229, 9, 20, 0.15); border: 1px solid rgba(229, 9, 20, 0.3); color: #ff6b6b; }
@keyframes slideIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }

@media (max-width: 640px) {
  .settings-page { padding: 1rem; }
  .settings-title { font-size: 1.5rem; }
  .settings-item { flex-direction: column; align-items: flex-start; gap: 0.75rem; }
  .scan-toast { left: 1rem; right: 1rem; bottom: 5rem; }
}
</style>
