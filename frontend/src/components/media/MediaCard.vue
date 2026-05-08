<template>
  <div
    ref="cardEl"
    class="media-card"
    @click="handleClick"
    role="button"
    tabindex="0"
    @keydown.enter="handleClick"
    @keydown.space.prevent="handleClick"
    :aria-label="`${media.title}${media.year ? ` (${media.year})` : ''}`"
  >
    <!-- Poster Container -->
    <div class="card-poster" :class="{ 'card-landscape': landscape }">
      <!-- Image (lazy-loaded) -->
      <img
        v-if="imageVisible && media.poster && !imageError"
        :src="media.poster"
        :alt="media.title"
        class="card-img"
        @load="imageLoaded = true"
        @error="imageError = true"
        :class="{ 'card-img-loaded': imageLoaded }"
      />

      <!-- Skeleton / Fallback -->
      <div
        v-if="!imageLoaded || !media.poster || imageError"
        class="card-placeholder"
        :class="{ 'skeleton': !imageVisible }"
      >
        <svg v-if="imageVisible" class="placeholder-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
        </svg>
      </div>

      <!-- Progress bar (continue watching) -->
      <div v-if="progress && progress > 0" class="card-progress">
        <div class="card-progress-bar" :style="{ width: `${Math.min(progress, 100)}%` }"></div>
      </div>

      <!-- Hover overlay -->
      <div class="card-overlay">
        <div class="overlay-play">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M8 5v14l11-7z"/>
          </svg>
        </div>
        <div class="overlay-info">
          <p class="overlay-title">{{ media.title }}</p>
          <div class="overlay-meta">
            <span v-if="media.year">{{ media.year }}</span>
            <span v-if="progress && progress > 0" class="overlay-progress-text">
              {{ Math.round(progress) }}%
            </span>
          </div>
        </div>
      </div>

      <!-- Type badge on top right -->
      <div class="card-badge" :class="`badge-${media.type}`">
        {{ typeShort(media.type) }}
      </div>
    </div>

    <!-- Card footer (title below) -->
    <div class="card-footer">
      <p class="card-title">{{ media.title }}</p>
      <p v-if="media.year" class="card-year">{{ media.year }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { MediaItem } from '@/types/media'

interface Props {
  media: MediaItem
  progress?: number // 0-100 percentage
  landscape?: boolean // for episode cards
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'click', media: MediaItem): void
}>()

const cardEl = ref<HTMLElement | null>(null)
const imageLoaded = ref(false)
const imageError = ref(false)
const imageVisible = ref(false)

let observer: IntersectionObserver | null = null

const handleClick = () => {
  emit('click', props.media)
}

const typeShort = (type: string) => {
  const map: Record<string, string> = { movie: 'F', series: 'S', episode: 'E' }
  return map[type] ?? ''
}

onMounted(() => {
  observer = new IntersectionObserver(
    (entries) => {
      if (entries[0]?.isIntersecting) {
        imageVisible.value = true
        observer?.disconnect()
      }
    },
    { rootMargin: '100px', threshold: 0 }
  )
  if (cardEl.value) observer.observe(cardEl.value)
})

onUnmounted(() => {
  observer?.disconnect()
})
</script>

<style scoped>
.media-card {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  cursor: pointer;
  outline: none;
  position: relative;
}

/* Poster */
.card-poster {
  position: relative;
  width: 100%;
  padding-bottom: 150%; /* 2:3 ratio */
  border-radius: 0.5rem;
  overflow: hidden;
  background: var(--color-cinema-dark-700);
  transition: transform var(--transition-base), box-shadow var(--transition-base);
}

.card-landscape {
  padding-bottom: 56.25%; /* 16:9 */
}

.media-card:hover .card-poster,
.media-card:focus .card-poster {
  transform: scale(1.05);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.7), 0 0 0 2px rgba(229, 9, 20, 0.4);
}

.card-img {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: 0;
  transition: opacity 0.4s ease;
}

.card-img-loaded {
  opacity: 1;
}

.card-placeholder {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--color-cinema-dark-700) 0%, var(--color-cinema-dark-600) 100%);
}

.placeholder-icon {
  width: 2.5rem;
  height: 2.5rem;
  color: rgba(255, 255, 255, 0.15);
}

/* Progress bar */
.card-progress {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: rgba(255, 255, 255, 0.15);
  z-index: 5;
}

.card-progress-bar {
  height: 100%;
  background: var(--color-accent);
  border-radius: 2px;
  transition: width 0.3s ease;
}

/* Hover overlay */
.card-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to top,
    rgba(0, 0, 0, 0.95) 0%,
    rgba(0, 0, 0, 0.6) 40%,
    transparent 70%
  );
  opacity: 0;
  transition: opacity var(--transition-base);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  z-index: 4;
}

.media-card:hover .card-overlay,
.media-card:focus .card-overlay {
  opacity: 1;
}

.overlay-play {
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  background: rgba(229, 9, 20, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  transform: scale(0.8);
  transition: transform var(--transition-fast);
  box-shadow: 0 4px 20px rgba(229, 9, 20, 0.6);
}

.media-card:hover .overlay-play {
  transform: scale(1);
}

.overlay-play svg {
  width: 1.5rem;
  height: 1.5rem;
  color: #fff;
  margin-left: 2px;
}

.overlay-info {
  position: absolute;
  bottom: 0.625rem;
  left: 0.625rem;
  right: 0.625rem;
}

.overlay-title {
  font-size: 0.8rem;
  font-weight: 600;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.overlay-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.7rem;
  color: var(--color-text-muted);
  margin-top: 0.125rem;
}

.overlay-progress-text {
  color: var(--color-accent);
  font-weight: 600;
}

/* Type badge */
.card-badge {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  width: 1.375rem;
  height: 1.375rem;
  border-radius: 0.25rem;
  font-size: 0.65rem;
  font-weight: 700;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 6;
  backdrop-filter: blur(4px);
}

.badge-movie {
  background: rgba(229, 9, 20, 0.8);
  color: #fff;
}

.badge-series {
  background: rgba(99, 102, 241, 0.8);
  color: #fff;
}

.badge-episode {
  background: rgba(34, 197, 94, 0.8);
  color: #fff;
}

/* Footer */
.card-footer {
  padding: 0 0.125rem;
}

.card-title {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.3;
}

.card-year {
  font-size: 0.75rem;
  color: var(--color-text-muted);
  margin-top: 0.125rem;
}
</style>
