<template>
  <div class="modal-overlay" v-if="show" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h3>Explorar Pastas</h3>
        <span class="current-path">{{ currentPath }}</span>
      </div>
      <div class="file-list">
        <button class="file-item folder-up" @click="navigate('..')">
          <span class="icon">⬆️</span> .. (Voltar)
        </button>
        <button class="file-item" v-for="dir in dirs" :key="dir.path" @click="navigate(dir.path)">
          <span class="icon">📁</span> {{ dir.name }}
        </button>
      </div>
      <div class="modal-actions">
        <button class="btn-ghost" @click="$emit('close')">Cancelar</button>
        <button class="btn-primary" @click="$emit('select', currentPath)">Selecionar Pasta</button>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import apiClient from '@/services/api'
defineProps<{ show: boolean }>()
const emit = defineEmits(['close', 'select'])
const currentPath = ref('/')
const dirs = ref<{name: string, path: string}[]>([])
const navigate = async (path: string) => {
  if (path === '..') {
    const parent = currentPath.value.substring(0, currentPath.value.lastIndexOf('/')) || '/'
    currentPath.value = parent
  } else {
    currentPath.value = path
  }
  try {
    const res = await apiClient.get(`/files/list?path=${currentPath.value}`)
    dirs.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}
onMounted(() => navigate('/'))
</script>
<style scoped>
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.7); backdrop-filter: blur(8px); display: flex; align-items: center; justify-content: center; z-index: 1000; }
.modal-content { background: rgba(20, 20, 25, 0.95); padding: 1.5rem; border-radius: 1rem; width: 450px; max-height: 80vh; display: flex; flex-direction: column; border: 1px solid rgba(255, 255, 255, 0.1); box-shadow: 0 20px 40px rgba(0,0,0,0.5); }
.modal-header { margin-bottom: 1rem; border-bottom: 1px solid rgba(255, 255, 255, 0.1); padding-bottom: 1rem; }
.modal-header h3 { font-size: 1.2rem; font-weight: 600; margin: 0 0 0.5rem 0; color: #fff; }
.current-path { font-size: 0.8rem; color: var(--color-text-muted); word-break: break-all; font-family: monospace; background: rgba(0,0,0,0.3); padding: 0.3rem 0.5rem; border-radius: 0.3rem; display: block; }
.file-list { flex: 1; overflow-y: auto; margin-bottom: 1.5rem; padding-right: 0.5rem; }
.file-list::-webkit-scrollbar { width: 6px; }
.file-list::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.2); border-radius: 3px; }
.file-item { width: 100%; display: flex; align-items: center; gap: 0.75rem; padding: 0.75rem; border: none; background: transparent; color: #e0e0e0; cursor: pointer; text-align: left; border-radius: 0.5rem; transition: background 0.2s; font-size: 0.9rem; }
.file-item:hover { background: rgba(255, 255, 255, 0.05); color: #fff; }
.icon { font-size: 1.2rem; opacity: 0.8; }
.modal-actions { display: flex; justify-content: flex-end; gap: 0.75rem; border-top: 1px solid rgba(255,255,255,0.1); padding-top: 1rem; }
.btn-primary { background: var(--color-accent); color: white; padding: 0.6rem 1.25rem; border-radius: 0.5rem; border: none; cursor: pointer; font-weight: 600; transition: transform 0.2s, box-shadow 0.2s; }
.btn-primary:hover { transform: translateY(-1px); box-shadow: 0 4px 12px var(--color-accent-glow); }
.btn-ghost { background: transparent; color: #aaa; padding: 0.6rem 1.25rem; border-radius: 0.5rem; border: 1px solid rgba(255,255,255,0.2); cursor: pointer; font-weight: 500; transition: all 0.2s; }
.btn-ghost:hover { background: rgba(255,255,255,0.05); color: #fff; border-color: rgba(255,255,255,0.4); }
</style>
