<template>
  <PlayerLayout>
    <div v-if="loading" class="loading-container">
      <LoadingSpinner size="large" />
      <p class="loading-text">Carregando player...</p>
    </div>
    <ErrorMessage
      v-else-if="error"
      :error="error"
      title="Falha ao carregar o player"
      :onRetry="handleRetry"
    />
    <div v-else-if="media" class="player-container">
      <VideoPlayer :media-id="props.id" />
    </div>
    <div v-else class="not-found-container">
      <h1 class="not-found-title">Mídia não encontrada</h1>
      <p class="not-found-message">A mídia solicitada não foi encontrada.</p>
      <button class="back-button" @click="goBack">
        Voltar
      </button>
    </div>
  </PlayerLayout>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import { usePlayerStore } from '@/stores/playerStore'
import PlayerLayout from '@/layouts/PlayerLayout.vue'
import VideoPlayer from '@/components/player/VideoPlayer.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import ErrorMessage from '@/components/common/ErrorMessage.vue'
interface Props {
  id: string
}
const props = defineProps<Props>()
const router = useRouter()
const playerStore = usePlayerStore()
const loading = ref(false)
const error = ref<string | null>(null)
const media = computed(() => playerStore.currentMedia)
const initializePlayer = async () => {
  loading.value = true
  error.value = null
  console.log('Initializing player for media:', props.id)
  try {
    await playerStore.initializePlayer(props.id)
  } catch (err: any) {
    error.value = err.message || 'Failed to initialize player'
    console.error('[PlayerPage] Failed to initialize player:', err)
  } finally {
    loading.value = false
  }
}
const handleRetry = () => {
  initializePlayer()
}
const goBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push({ name: 'home' })
  }
}
const cleanup = async () => {
  try {
    await playerStore.saveProgressImmediate()
    console.debug('[PlayerPage] Progress saved on unmount')
  } catch (err) {
    console.error('[PlayerPage] Failed to save progress on unmount:', err)
  } finally {
    playerStore.clearPlayer()
  }
}
onMounted(() => {
  initializePlayer()
})
onBeforeUnmount(() => {
  cleanup()
})
</script>
<style scoped>
.player-container {
  width: 100vw;
  height: 100vh;
  background: #000000;
}
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  gap: 1.5rem;
  background: #000000;
}
.loading-text {
  font-size: 1.125rem;
  color: #9ca3af;
  margin: 0;
}
.not-found-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 2rem;
  text-align: center;
  background: #000000;
}
.not-found-title {
  font-size: 2rem;
  font-weight: 700;
  color: #ffffff;
  margin: 0 0 1rem 0;
}
.not-found-message {
  font-size: 1.125rem;
  color: #9ca3af;
  margin: 0 0 2rem 0;
}
.back-button {
  padding: 0.875rem 2rem;
  font-size: 1rem;
  font-weight: 600;
  background: rgba(239, 68, 68, 0.9);
  color: #ffffff;
  border: none;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
}
.back-button:hover {
  background: rgba(239, 68, 68, 1);
  box-shadow: 0 6px 20px rgba(239, 68, 68, 0.6);
  transform: translateY(-2px);
}
.back-button:active {
  transform: translateY(0);
}
/* Responsive adjustments */
@media (max-width: 768px) {
  .not-found-title {
    font-size: 1.5rem;
  }
  .not-found-message {
    font-size: 1rem;
  }
  .back-button {
    padding: 0.75rem 1.5rem;
    font-size: 0.875rem;
  }
}
</style>
