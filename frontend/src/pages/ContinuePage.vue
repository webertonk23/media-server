<template>
  <DefaultLayout>
    <div class="catalog-page">
      <div class="catalog-header">
        <h1 class="catalog-title">Continuar Assistindo</h1>
        <p class="catalog-subtitle">Retome de onde parou</p>
      </div>
      <div v-if="loading" class="catalog-skeleton">
        <div v-for="i in 10" :key="i" class="skeleton card-skeleton"></div>
      </div>
      <div v-else-if="items.length === 0" class="catalog-empty">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
        </svg>
        <p>Nenhum item em progresso</p>
        <span>Comece a assistir algo para aparecer aqui</span>
        <router-link to="/" class="btn-primary">Explorar biblioteca</router-link>
      </div>
      <div v-else class="catalog-grid">
        <div v-for="(item, index) in items" :key="item.media.id"
          class="catalog-card-wrap" :style="{ '--delay': `${Math.min(index * 0.03, 0.5)}s` }">
          <MediaCard
            :media="item.media"
            :progress="item.progressPercentage"
            @click="handleClick"
          />
        </div>
      </div>
    </div>
  </DefaultLayout>
</template>
<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useContinueWatchingStore } from '@/stores/continueWatchingStore'
import type { MediaItem } from '@/types/media'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import MediaCard from '@/components/media/MediaCard.vue'
const router = useRouter()
const store = useContinueWatchingStore()
const { items, loading } = store
const handleClick = (media: MediaItem) => {
  router.push({ name: 'media-detail', params: { id: media.id } })
}
onMounted(() => store.fetchContinueWatching())
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
.catalog-empty span { font-size: 0.875rem; color: var(--color-text-muted); margin-bottom: 0.5rem; }
@media (max-width: 640px) { .catalog-page { padding: 1rem; } .catalog-title { font-size: 1.5rem; } }
</style>
