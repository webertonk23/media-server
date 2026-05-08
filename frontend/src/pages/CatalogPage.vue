<template>
  <DefaultLayout>
    <div class="catalog-page">
      <!-- Page Header -->
      <div class="catalog-header">
        <h1 class="catalog-title">{{ pageTitle }}</h1>
        <p class="catalog-subtitle">{{ pageSubtitle }}</p>
      </div>

      <!-- Loading skeleton -->
      <div v-if="loading && items.length === 0" class="catalog-skeleton">
        <div v-for="i in 20" :key="i" class="skeleton card-skeleton"></div>
      </div>

      <!-- Error -->
      <div v-else-if="error && items.length === 0" class="catalog-error">
        <p>{{ error }}</p>
        <button class="btn-primary" @click="load">Tentar novamente</button>
      </div>

      <!-- Grid -->
      <div v-else class="catalog-grid">
        <div
          v-for="(item, index) in items"
          :key="item.id"
          class="catalog-card-wrap"
          :style="{ '--delay': `${Math.min(index * 0.03, 0.5)}s` }"
        >
          <MediaCard
            :media="item"
            @click="handleClick"
          />
        </div>
      </div>

      <!-- Empty state -->
      <div v-if="!loading && items.length === 0 && !error" class="catalog-empty">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4"/>
        </svg>
        <p>Nenhum item encontrado</p>
        <span>Adicione arquivos de mídia e faça um scan da biblioteca</span>
        <button class="btn-primary" @click="triggerScan">Escanear biblioteca</button>
      </div>

      <!-- Load more sentinel -->
      <div ref="sentinelRef" class="load-more-sentinel">
        <div v-if="isLoadingMore" class="spinner"></div>
      </div>
    </div>
  </DefaultLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMediaStore } from '@/stores/mediaStore'
import { useInfiniteScroll } from '@/composables/useInfiniteScroll'
import type { MediaItem } from '@/types/media'
import apiClient from '@/services/api'

import DefaultLayout from '@/layouts/DefaultLayout.vue'
import MediaCard from '@/components/media/MediaCard.vue'

interface Props {
  type: 'movie' | 'series'
}

const props = defineProps<Props>()

const router = useRouter()
const mediaStore = useMediaStore()

const loading = ref(false)
const error = ref<string | null>(null)

const pageTitle = computed(() => props.type === 'movie' ? 'Filmes' : 'Séries')
const pageSubtitle = computed(() =>
  props.type === 'movie'
    ? 'Todos os filmes da sua biblioteca'
    : 'Todas as séries da sua biblioteca'
)

const items = computed(() =>
  mediaStore.mediaItems.filter(i => i.type === props.type)
)

const load = async () => {
  loading.value = true
  error.value = null
  try {
    await mediaStore.fetchMedia({ page: 1, limit: 50, type: props.type })
  } catch (err: any) {
    error.value = err.message || 'Erro ao carregar'
  } finally {
    loading.value = false
  }
}

const loadMore = async () => {
  if (loading.value || !mediaStore.pagination.hasMore) return
  try {
    await mediaStore.fetchMedia({
      page: mediaStore.pagination.page + 1,
      limit: 50,
      type: props.type,
    }, true)
  } catch {}
}

const handleClick = (media: MediaItem) => {
  router.push({ name: 'media-detail', params: { id: media.id } })
}

const triggerScan = async () => {
  try {
    await apiClient.post('/scan')
    await load()
  } catch {}
}

const { sentinelRef, isLoading: isLoadingMore } = useInfiniteScroll(loadMore, {
  rootMargin: '200px',
  threshold: 0,
})

defineExpose({ sentinelRef })

onMounted(load)
</script>

<style scoped>
.catalog-page {
  padding: 2rem;
  min-height: 100vh;
}

.catalog-header {
  margin-bottom: 2rem;
}

.catalog-title {
  font-family: var(--font-display);
  font-size: 2rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 0.375rem;
}

.catalog-subtitle {
  font-size: 0.9rem;
  color: var(--color-text-muted);
}

/* Grid */
.catalog-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 1.25rem;
}

.catalog-card-wrap {
  animation: fadeInUp 0.4s ease both;
  animation-delay: var(--delay, 0s);
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Skeleton */
.catalog-skeleton {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 1.25rem;
}

.card-skeleton {
  height: 225px;
  border-radius: 0.5rem;
}

/* Error */
.catalog-error {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 4rem 2rem;
  text-align: center;
  color: var(--color-text-secondary);
}

/* Empty */
.catalog-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding: 5rem 2rem;
  text-align: center;
}

.catalog-empty svg {
  width: 4rem;
  height: 4rem;
  color: rgba(255,255,255,0.15);
}

.catalog-empty p {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.catalog-empty span {
  font-size: 0.875rem;
  color: var(--color-text-muted);
  max-width: 320px;
  margin-bottom: 0.5rem;
}

/* Load more */
.load-more-sentinel {
  display: flex;
  justify-content: center;
  padding: 2rem;
}

.spinner {
  width: 2rem;
  height: 2rem;
  border: 2px solid rgba(255,255,255,0.1);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

@media (min-width: 640px) {
  .catalog-grid,
  .catalog-skeleton {
    grid-template-columns: repeat(auto-fill, minmax(165px, 1fr));
  }
}

@media (min-width: 1024px) {
  .catalog-grid,
  .catalog-skeleton {
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  }
}

@media (max-width: 640px) {
  .catalog-page { padding: 1rem; }
  .catalog-title { font-size: 1.5rem; }
}
</style>
