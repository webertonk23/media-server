<template>
  <div class="media-row">
    <!-- Row Header -->
    <div class="row-header">
      <h2 class="row-title">{{ title }}</h2>
      <router-link v-if="viewAllPath" :to="viewAllPath" class="view-all-link">
        Ver todos
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="9 18 15 12 9 6"/>
        </svg>
      </router-link>
    </div>

    <!-- Scrollable Container -->
    <div class="row-wrapper">
      <!-- Left Arrow -->
      <button
        v-if="showLeftBtn"
        class="scroll-btn scroll-btn-left"
        @click="scrollLeft"
        aria-label="Rolar para esquerda"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <polyline points="15 18 9 12 15 6"/>
        </svg>
      </button>

      <!-- Cards Container -->
      <div ref="scrollRef" class="row-scroll" @scroll="updateBtns">
        <div class="row-cards">
          <div
            v-for="(item, index) in items"
            :key="item.id"
            class="row-card-wrapper"
            :style="{ '--delay': `${index * 0.04}s` }"
          >
            <MediaCard
              :media="item"
              :progress="getProgress(item.id)"
              @click="$emit('media-click', item)"
            />
          </div>
        </div>
      </div>

      <!-- Right Arrow -->
      <button
        v-if="showRightBtn"
        class="scroll-btn scroll-btn-right"
        @click="scrollRight"
        aria-label="Rolar para direita"
      >
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
          <polyline points="9 18 15 12 9 6"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { MediaItem } from '@/types/media'
import MediaCard from './MediaCard.vue'

interface Props {
  title: string
  items: MediaItem[]
  viewAllPath?: string
  progressMap?: Record<string, number> // mediaId -> percentage
}

const props = defineProps<Props>()

defineEmits<{
  (e: 'media-click', media: MediaItem): void
}>()

const scrollRef = ref<HTMLElement | null>(null)
const showLeftBtn = ref(false)
const showRightBtn = ref(false)

const getProgress = (id: string) => props.progressMap?.[id]

const updateBtns = () => {
  const el = scrollRef.value
  if (!el) return
  showLeftBtn.value = el.scrollLeft > 10
  showRightBtn.value = el.scrollLeft < el.scrollWidth - el.clientWidth - 10
}

const scrollLeft = () => {
  scrollRef.value?.scrollBy({ left: -(scrollRef.value.clientWidth * 0.75), behavior: 'smooth' })
}

const scrollRight = () => {
  scrollRef.value?.scrollBy({ left: scrollRef.value.clientWidth * 0.75, behavior: 'smooth' })
}

let resizeObs: ResizeObserver | null = null

onMounted(() => {
  updateBtns()
  if (scrollRef.value) {
    resizeObs = new ResizeObserver(updateBtns)
    resizeObs.observe(scrollRef.value)
  }
  window.addEventListener('resize', updateBtns)
})

onUnmounted(() => {
  resizeObs?.disconnect()
  window.removeEventListener('resize', updateBtns)
})
</script>

<style scoped>
.media-row {
  margin-bottom: 2.5rem;
}

/* Header */
.row-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 2rem;
  margin-bottom: 1rem;
}

.row-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-text-primary);
  letter-spacing: -0.01em;
}

.view-all-link {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-accent);
  text-decoration: none;
  transition: opacity var(--transition-fast);
  white-space: nowrap;
}

.view-all-link:hover {
  opacity: 0.8;
}

.view-all-link svg {
  width: 0.875rem;
  height: 0.875rem;
  transition: transform var(--transition-fast);
}

.view-all-link:hover svg {
  transform: translateX(3px);
}

/* Wrapper */
.row-wrapper {
  position: relative;
}

/* Scroll container */
.row-scroll {
  overflow-x: auto;
  overflow-y: visible;
  scrollbar-width: none;
  padding: 0.5rem 2rem 1rem;
}

.row-scroll::-webkit-scrollbar { display: none; }

.row-cards {
  display: flex;
  gap: 0.875rem;
}

/* Individual card wrapper */
.row-card-wrapper {
  flex-shrink: 0;
  width: 160px;
  animation: fadeIn 0.4s ease both;
  animation-delay: var(--delay, 0s);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Scroll Buttons */
.scroll-btn {
  position: absolute;
  top: 50%;
  transform: translateY(-60%);
  z-index: 20;
  width: 2.75rem;
  height: 2.75rem;
  border-radius: 50%;
  background: rgba(10, 10, 15, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.12);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
  backdrop-filter: blur(8px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
}

.scroll-btn:hover {
  background: var(--color-accent);
  border-color: var(--color-accent);
  box-shadow: 0 4px 20px var(--color-accent-glow);
  transform: translateY(-60%) scale(1.1);
}

.scroll-btn svg {
  width: 1.125rem;
  height: 1.125rem;
}

.scroll-btn-left { left: 0.5rem; }
.scroll-btn-right { right: 0.5rem; }

/* Responsive */
@media (min-width: 768px) {
  .row-card-wrapper {
    width: 180px;
  }
}

@media (min-width: 1280px) {
  .row-card-wrapper {
    width: 200px;
  }

  .row-cards {
    gap: 1rem;
  }
}

@media (max-width: 640px) {
  .row-header {
    padding: 0 1rem;
  }

  .row-scroll {
    padding: 0.5rem 1rem 1rem;
  }

  .row-card-wrapper {
    width: 130px;
  }

  .row-cards {
    gap: 0.625rem;
  }

  .scroll-btn {
    display: none;
  }
}
</style>
