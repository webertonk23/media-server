<template>
  <DefaultLayout>
    <div class="catalog-page">
      <div class="catalog-header">
        <h1 class="catalog-title">Recém Adicionados</h1>
        <p class="catalog-subtitle">Novidades na sua biblioteca</p>
      </div>

      <div v-if="loading && items.length === 0" class="catalog-skeleton">
        <div v-for="i in 20" :key="i" class="skeleton card-skeleton"></div>
      </div>

      <div v-else-if="items.length === 0" class="catalog-empty">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
        </svg>
        <p>Biblioteca vazia</p>
        <span>Adicione arquivos de mídia e faça um scan</span>
      </div>

      <div v-else class="catalog-grid">
        <div v-for="(item, index) in items" :key="item.id"
          class="catalog-card-wrap" :style="{ '--delay': `${Math.min(index * 0.03, 0.5)}s` }">
          <MediaCard :media="item" @click="handleClick" />
        </div>
      </div>

      <div ref="sentinelRef" class="sentinel">
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
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import MediaCard from '@/components/media/MediaCard.vue'

const router = useRouter()
const mediaStore = useMediaStore()
const loading = ref(false)
const items = computed(() => mediaStore.mediaItems.filter(i => i.type !== 'episode'))

const load = async () => {
  loading.value = true
  try { await mediaStore.fetchMedia({ page: 1, limit: 50 }) }
  catch (e) { console.error(e) }
  finally { loading.value = false }
}

const loadMore = async () => {
  if (loading.value || !mediaStore.pagination.hasMore) return
  try { await mediaStore.fetchMedia({ page: mediaStore.pagination.page + 1, limit: 50 }, true) }
  catch {}
}

const handleClick = (media: MediaItem) => router.push({ name: 'media-detail', params: { id: media.id } })

const { sentinelRef, isLoading: isLoadingMore } = useInfiniteScroll(loadMore, { rootMargin: '200px', threshold: 0 })
defineExpose({ sentinelRef })
onMounted(load)
</script>

<style scoped>
.catalog-page { padding: 2rem; min-height: 100vh; }
.catalog-header { margin-bottom: 2rem; }
.catalog-title { font-family: var(--font-display); font-size: 2rem; font-weight: 700; margin-bottom: 0.375rem; }
.catalog-subtitle { font-size: 0.9rem; color: var(--color-text-muted); }
.catalog-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); gap: 1.25rem; }
.catalog-card-wrap { animation: fadeInUp 0.4s ease both; animation-delay: var(--delay, 0s); }
@keyframes fadeInUp { from { opacity: 0; transform: translateY(12px); } to { opacity: 1; transform: translateY(0); } }
.catalog-skeleton { display: grid; grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); gap: 1.25rem; }
.card-skeleton { height: 225px; border-radius: 0.5rem; }
.catalog-empty { display: flex; flex-direction: column; align-items: center; gap: 0.75rem; padding: 5rem 2rem; text-align: center; }
.catalog-empty svg { width: 4rem; height: 4rem; color: rgba(255,255,255,0.15); }
.catalog-empty p { font-size: 1.25rem; font-weight: 600; color: var(--color-text-secondary); }
.catalog-empty span { font-size: 0.875rem; color: var(--color-text-muted); }
.sentinel { display: flex; justify-content: center; padding: 2rem; }
.spinner { width: 2rem; height: 2rem; border: 2px solid rgba(255,255,255,0.1); border-top-color: var(--color-accent); border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
@media (max-width: 640px) { .catalog-page { padding: 1rem; } .catalog-title { font-size: 1.5rem; } }
</style>
