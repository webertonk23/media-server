<template>
  <DefaultLayout>
    <div v-if="loading" class="detail-loading">
      <div class="spinner"></div>
    </div>
    <ErrorMessage
      v-else-if="error"
      :error="error"
      title="Falha ao carregar detalhes"
      :on-retry="loadMediaDetails"
    />
    <MediaDetails
      v-else-if="media"
      :media="media"
      :progress="progress"
      :seasons="seasons"
      :episodes="episodes"
      :selected-season-id="selectedSeasonId"
      @play="handlePlay"
      @continue="handleContinue"
      @back="goBack"
      @season-select="handleSeasonSelect"
    />
    <div v-else class="not-found">
      <h1>Mídia não encontrada</h1>
      <button class="btn-primary" @click="goBack">Voltar</button>
    </div>
  </DefaultLayout>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMediaStore } from '@/stores/mediaStore'
import * as progressService from '@/services/progressService'
import { getSeriesSeasons, getSeasonEpisodes } from '@/services/mediaService'
import type { ProgressData } from '@/types/player'
import type { Season, Episode } from '@/types/media'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import MediaDetails from '@/components/media/MediaDetails.vue'
import ErrorMessage from '@/components/common/ErrorMessage.vue'
interface Props { id: string }
const props = defineProps<Props>()
const router = useRouter()
const mediaStore = useMediaStore()
const loading = ref(false)
const error = ref<string | null>(null)
const progress = ref<ProgressData | null>(null)
const media = computed(() => mediaStore.currentMedia)
const seasons = ref<Season[]>([])
const episodes = ref<Episode[]>([])
const selectedSeasonId = ref<string | null>(null)
const loadMediaDetails = async () => {
  loading.value = true
  error.value = null
  try {
    await mediaStore.fetchMediaById(props.id)
    if (mediaStore.currentMedia) {
      try {
        progress.value = await progressService.getProgress(props.id)
      } catch {
        progress.value = null
      }
      if (mediaStore.currentMedia.type === 'series') {
        try {
          const fetchedSeasons = await getSeriesSeasons(props.id)
          seasons.value = fetchedSeasons || []
          if (seasons.value.length > 0) {
            await handleSeasonSelect(seasons.value[0]?.id!)
          }
        } catch (e) {
          console.error('Failed to load seasons', e)
        }
      }
    }
  } catch (err: any) {
    error.value = err.message || 'Falha ao carregar detalhes'
  } finally {
    loading.value = false
  }
}
const handleSeasonSelect = async (seasonId: string) => {
  selectedSeasonId.value = seasonId
  try {
    episodes.value = await getSeasonEpisodes(seasonId)
  } catch (e) {
    console.error('Failed to load episodes', e)
    episodes.value = []
  }
}
const handlePlay = (item?: any) => {
  if (item && item.type === 'episode') {
    router.push({ name: 'player', params: { id: item.id } })
  } else {
    router.push({ name: 'player', params: { id: props.id } })
  }
}
const handleContinue = () => router.push({ name: 'player', params: { id: props.id } })
const goBack = () => window.history.length > 1 ? router.back() : router.push({ name: 'home' })
onMounted(loadMediaDetails)
</script>
<style scoped>
.detail-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 80vh;
}
.spinner {
  width: 2.5rem;
  height: 2.5rem;
  border: 2px solid rgba(255,255,255,0.1);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
.not-found {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  gap: 1.5rem;
  text-align: center;
}
.not-found h1 {
  font-size: 1.75rem;
  font-weight: 600;
}
</style>
