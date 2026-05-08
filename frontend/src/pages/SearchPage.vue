<template>
  <DefaultLayout>
    <div class="search-page">
      <!-- Search header -->
      <div class="search-header">
        <h1 class="search-title">Buscar</h1>
        <div class="search-bar-wrap">
          <div class="search-input-wrap">
            <svg class="s-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
            <input
              ref="inputRef"
              v-model="query"
              type="text"
              placeholder="Buscar filmes, séries..."
              class="search-input"
              @input="debouncedSearch"
              @keyup.enter="doSearch"
              autofocus
            />
            <button v-if="query" class="s-clear" @click="clearSearch">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Searching indicator -->
      <div v-if="loading" class="search-loading">
        <div v-for="i in 8" :key="i" class="skeleton search-skel"></div>
      </div>

      <!-- Results -->
      <div v-else-if="results.length > 0" class="search-results">
        <p class="results-count">
          {{ results.length }} resultado{{ results.length !== 1 ? 's' : '' }} para
          <strong>"{{ lastQuery }}"</strong>
        </p>
        <div class="results-grid">
          <div
            v-for="(item, index) in results"
            :key="item.id"
            class="result-card-wrap"
            :style="{ '--delay': `${Math.min(index * 0.03, 0.4)}s` }"
          >
            <MediaCard :media="item" @click="handleClick" />
          </div>
        </div>
      </div>

      <!-- Empty result after search -->
      <div v-else-if="hasSearched && !loading" class="search-empty">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <p>Nenhum resultado para <strong>"{{ lastQuery }}"</strong></p>
        <span>Tente um termo diferente ou verifique a ortografia</span>
      </div>

      <!-- Initial state (nothing typed yet) -->
      <div v-else class="search-initial">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1">
          <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <p>Digite para buscar na sua biblioteca</p>
      </div>
    </div>
  </DefaultLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMediaStore } from '@/stores/mediaStore'
import type { MediaItem } from '@/types/media'

import DefaultLayout from '@/layouts/DefaultLayout.vue'
import MediaCard from '@/components/media/MediaCard.vue'

const router = useRouter()
const route = useRoute()
const mediaStore = useMediaStore()

const query = ref('')
const lastQuery = ref('')
const loading = ref(false)
const hasSearched = ref(false)
const inputRef = ref<HTMLInputElement | null>(null)

const results = computed(() =>
  hasSearched.value && lastQuery.value ? mediaStore.mediaItems : []
)

let debounceTimer: ReturnType<typeof setTimeout> | null = null

const debouncedSearch = () => {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(doSearch, 350)
}

const doSearch = async () => {
  const q = query.value.trim()
  if (!q) {
    hasSearched.value = false
    mediaStore.clearAll()
    return
  }
  hasSearched.value = true
  lastQuery.value = q
  loading.value = true
  try {
    await mediaStore.searchMedia(q)
  } catch (err) {
    console.error('[SearchPage] search failed:', err)
  } finally {
    loading.value = false
  }
}

const clearSearch = () => {
  query.value = ''
  hasSearched.value = false
  mediaStore.clearAll()
  inputRef.value?.focus()
}

const handleClick = (media: MediaItem) => {
  router.push({ name: 'media-detail', params: { id: media.id } })
}

onMounted(() => {
  // Pre-fill from URL query param
  const q = route.query.q as string
  if (q) {
    query.value = q
    doSearch()
  }
  inputRef.value?.focus()
})
</script>

<style scoped>
.search-page {
  padding: 2rem;
  min-height: 100vh;
}

.search-header {
  margin-bottom: 2.5rem;
  max-width: 640px;
}

.search-title {
  font-family: var(--font-display);
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 1.25rem;
}

.search-bar-wrap {
  position: relative;
}

.search-input-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.s-icon {
  position: absolute;
  left: 1rem;
  width: 1.125rem;
  height: 1.125rem;
  color: var(--color-text-muted);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 0.875rem 3rem 0.875rem 3rem;
  background: rgba(255,255,255,0.06);
  border: 1px solid rgba(255,255,255,0.12);
  border-radius: 0.625rem;
  color: var(--color-text-primary);
  font-size: 1rem;
  transition: all var(--transition-fast);
}

.search-input::placeholder { color: var(--color-text-muted); }

.search-input:focus {
  outline: none;
  background: rgba(255,255,255,0.09);
  border-color: rgba(229, 9, 20, 0.45);
  box-shadow: 0 0 0 3px rgba(229, 9, 20, 0.1);
}

.s-clear {
  position: absolute;
  right: 1rem;
  width: 1.125rem;
  height: 1.125rem;
  color: var(--color-text-muted);
  transition: color var(--transition-fast);
}
.s-clear:hover { color: var(--color-text-primary); }
.s-clear svg { width: 100%; height: 100%; }

/* Loading skeletons */
.search-loading {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 1.25rem;
}

.search-skel {
  height: 225px;
  border-radius: 0.5rem;
}

/* Results */
.search-results {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.results-count {
  font-size: 0.875rem;
  color: var(--color-text-muted);
  margin-bottom: 1.5rem;
}

.results-count strong { color: var(--color-text-primary); }

.results-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 1.25rem;
}

.result-card-wrap {
  animation: fadeInUp 0.35s ease both;
  animation-delay: var(--delay, 0s);
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Empty / initial states */
.search-empty,
.search-initial {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.875rem;
  padding: 6rem 2rem;
  text-align: center;
  color: var(--color-text-muted);
}

.search-empty svg,
.search-initial svg {
  width: 5rem;
  height: 5rem;
  opacity: 0.2;
  margin-bottom: 0.5rem;
}

.search-empty p { font-size: 1.125rem; color: var(--color-text-secondary); }
.search-empty strong { color: var(--color-text-primary); }
.search-empty span { font-size: 0.875rem; }
.search-initial p { font-size: 1rem; color: var(--color-text-secondary); }

@media (max-width: 640px) {
  .search-page { padding: 1rem; }
  .search-title { font-size: 1.5rem; }
  .results-grid,
  .search-loading {
    grid-template-columns: repeat(auto-fill, minmax(130px, 1fr));
    gap: 0.875rem;
  }
}
</style>
