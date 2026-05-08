<template>
  <DefaultLayout>
    <!-- Loading state -->
    <div v-if="loading && mediaItems.length === 0" class="home-loading">
      <div v-for="i in 3" :key="i" class="loading-section">
        <div class="skeleton loading-title"></div>
        <div class="loading-cards">
          <div v-for="j in 6" :key="j" class="skeleton loading-card"></div>
        </div>
      </div>
    </div>

    <!-- Error state -->
    <div v-else-if="error && mediaItems.length === 0" class="home-error">
      <div class="error-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="8" x2="12" y2="12"/>
          <line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
      </div>
      <h2>Erro ao carregar mídia</h2>
      <p>{{ error }}</p>
      <button class="btn-primary" @click="loadInitialData">Tentar novamente</button>
    </div>

    <!-- Content -->
    <div v-else class="home-page">
      <!-- Hero Banner -->
      <HeroBanner
        v-if="featuredMedia"
        :media="featuredMedia"
        :all-media="heroMedia"
        @play="handlePlay"
        @more-info="handleMoreInfo"
      />

      <!-- Continue Watching -->
      <section v-if="continueWatchingItems.length > 0" class="section">
        <MediaRow
          title="Continuar Assistindo"
          :items="continueWatchingItems.map(i => i.media)"
          :progress-map="continueWatchingProgressMap"
          view-all-path="/continue"
          @media-click="handleMediaClick"
        />
      </section>

      <!-- Recently Added -->
      <section v-if="recentMedia.length > 0" class="section">
        <MediaRow
          title="Recém Adicionados"
          :items="recentMedia"
          view-all-path="/recent"
          @media-click="handleMediaClick"
        />
      </section>

      <!-- Movies -->
      <section v-if="movies.length > 0" class="section">
        <MediaRow
          title="Filmes"
          :items="movies"
          view-all-path="/movies"
          @media-click="handleMediaClick"
        />
      </section>

      <!-- Series -->
      <section v-if="series.length > 0" class="section">
        <MediaRow
          title="Séries"
          :items="series"
          view-all-path="/series"
          @media-click="handleMediaClick"
        />
      </section>

      <!-- Infinite scroll sentinel -->
      <div ref="sentinelRef" class="sentinel">
        <div v-if="isLoadingMore" class="sentinel-loader">
          <div class="spinner"></div>
        </div>
      </div>
    </div>
  </DefaultLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMediaStore } from '@/stores/mediaStore'
import { useContinueWatchingStore } from '@/stores/continueWatchingStore'
import { useInfiniteScroll } from '@/composables/useInfiniteScroll'
import type { MediaItem } from '@/types/media'

import DefaultLayout from '@/layouts/DefaultLayout.vue'
import HeroBanner from '@/components/media/HeroBanner.vue'
import MediaRow from '@/components/media/MediaRow.vue'

const router = useRouter()
const mediaStore = useMediaStore()
const continueWatchingStore = useContinueWatchingStore()

const loading = ref(false)
const error = ref<string | null>(null)

const mediaItems = computed(() => mediaStore.mediaItems)
const continueWatchingItems = computed(() => continueWatchingStore.items)

// Map mediaId -> progress percentage for continue watching
const continueWatchingProgressMap = computed(() => {
  const map: Record<string, number> = {}
  continueWatchingStore.items.forEach(item => {
    map[item.media.id] = item.progressPercentage
  })
  return map
})

// Featured media for hero (first 5 items)
const heroMedia = computed(() => mediaItems.value.slice(0, 5))
const featuredMedia = computed(() => mediaItems.value[0] || null)

// Recently added (up to 20)
const recentMedia = computed(() => mediaItems.value.slice(0, 20))

// Movies only
const movies = computed(() =>
  mediaItems.value.filter(item => item.type === 'movie').slice(0, 20)
)

// Series only
const series = computed(() =>
  mediaItems.value.filter(item => item.type === 'series').slice(0, 20)
)

const loadInitialData = async () => {
  loading.value = true
  error.value = null

  try {
    await Promise.all([
      mediaStore.fetchMedia({ page: 1, limit: 50 }),
      continueWatchingStore.fetchContinueWatching().catch(() => {}),
    ])
  } catch (err: any) {
    error.value = err.message || 'Falha ao carregar mídia'
    console.error('[HomePage] Failed to load:', err)
  } finally {
    loading.value = false
  }
}

const loadMoreMedia = async () => {
  if (loading.value || !mediaStore.pagination.hasMore) return
  try {
    await mediaStore.fetchMedia({ page: mediaStore.pagination.page + 1, limit: 50 }, true)
  } catch (err: any) {
    console.error('[HomePage] Load more failed:', err)
  }
}

const handlePlay = (media: MediaItem) => {
  router.push({ name: 'player', params: { id: media.id } })
}

const handleMoreInfo = (media: MediaItem) => {
  router.push({ name: 'media-detail', params: { id: media.id } })
}

const handleMediaClick = (media: MediaItem) => {
  router.push({ name: 'media-detail', params: { id: media.id } })
}

const { sentinelRef, isLoading: isLoadingMore } = useInfiniteScroll(loadMoreMedia, {
  rootMargin: '200px',
  threshold: 0,
})

defineExpose({ sentinelRef })

onMounted(() => {
  loadInitialData()
})
</script>

<style scoped>
.home-page {
  min-height: 100vh;
  padding-bottom: 3rem;
}

/* Loading skeleton */
.home-loading {
  padding: 2rem;
  display: flex;
  flex-direction: column;
  gap: 3rem;
}

.loading-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.loading-title {
  height: 1.5rem;
  width: 12rem;
  border-radius: 0.375rem;
}

.loading-cards {
  display: flex;
  gap: 1rem;
}

.loading-card {
  width: 160px;
  height: 240px;
  border-radius: 0.5rem;
  flex-shrink: 0;
}

/* Error */
.home-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 60vh;
  gap: 1rem;
  text-align: center;
  padding: 2rem;
}

.error-icon svg {
  width: 4rem;
  height: 4rem;
  color: var(--color-accent);
}

.home-error h2 {
  font-size: 1.5rem;
  font-weight: 600;
}

.home-error p {
  color: var(--color-text-secondary);
  max-width: 400px;
}

/* Section */
.section {
  position: relative;
}

/* Sentinel */
.sentinel {
  display: flex;
  justify-content: center;
  padding: 2rem;
}

.sentinel-loader {
  display: flex;
  align-items: center;
  justify-content: center;
}

.spinner {
  width: 2rem;
  height: 2rem;
  border: 2px solid rgba(255, 255, 255, 0.1);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
