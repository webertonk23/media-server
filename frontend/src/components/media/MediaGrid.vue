<template>
  <div class="media-grid-container">
    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <LoadingSpinner size="large" />
    </div>
    <!-- Media Grid -->
    <div v-else-if="items.length > 0" class="media-grid">
      <MediaCard
        v-for="item in items"
        :key="item.id"
        :media="item"
        @click="handleMediaClick"
      />
    </div>
    <!-- Empty State -->
    <div v-else class="empty-state">
      <svg
        class="empty-icon"
        xmlns="http:
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z"
        />
      </svg>
      <p class="empty-message">No media items found</p>
    </div>
  </div>
</template>
<script setup lang="ts">
import type { MediaItem } from '@/types/media'
import MediaCard from './MediaCard.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
interface Props {
  items: MediaItem[]
  loading?: boolean
}
const props = withDefaults(defineProps<Props>(), {
  loading: false
})
interface Emits {
  (e: 'media-click', media: MediaItem): void
}
const emit = defineEmits<Emits>()
const handleMediaClick = (media: MediaItem) => {
  emit('media-click', media)
}
</script>
<style scoped>
.media-grid-container {
  width: 100%;
  min-height: 200px;
}
.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}
.media-grid {
  display: grid;
  gap: 1.5rem;
  width: 100%;
  /* Responsive grid columns based on screen size */
  /* Mobile: 2 columns */
  grid-template-columns: repeat(2, 1fr);
}
/* Tablet: 3 columns (768px - 1023px) */
@media (min-width: 768px) {
  .media-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: 1.75rem;
  }
}
/* Desktop: 4 columns (1024px - 1279px) */
@media (min-width: 1024px) {
  .media-grid {
    grid-template-columns: repeat(4, 1fr);
    gap: 2rem;
  }
}
/* Large Desktop: 5 columns (1280px - 1919px) */
@media (min-width: 1280px) {
  .media-grid {
    grid-template-columns: repeat(5, 1fr);
  }
}
/* Ultrawide: 6 columns (>= 1920px) */
@media (min-width: 1920px) {
  .media-grid {
    grid-template-columns: repeat(6, 1fr);
    gap: 2.5rem;
  }
}
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  padding: 2rem;
}
.empty-icon {
  width: 4rem;
  height: 4rem;
  color: #4a4a4a;
  opacity: 0.5;
  margin-bottom: 1rem;
}
.empty-message {
  font-size: 1.125rem;
  color: #6a6a6a;
  margin: 0;
}
/* Responsive empty state */
@media (max-width: 768px) {
  .empty-icon {
    width: 3rem;
    height: 3rem;
  }
  .empty-message {
    font-size: 1rem;
  }
}
</style>
