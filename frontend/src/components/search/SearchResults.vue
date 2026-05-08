<template>
  <div class="search-results">
    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <LoadingSpinner size="large" />
      <p class="loading-text">Searching...</p>
    </div>

    <!-- No Results State (Requirement 8.4) -->
    <div v-else-if="!loading && results.length === 0 && hasSearched" class="empty-state">
      <svg
        class="empty-icon"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
      <h3 class="empty-title">No results found</h3>
      <p class="empty-message">Try searching with different keywords</p>
    </div>

    <!-- Search Results Grid (Requirement 8.3) -->
    <div v-else-if="results.length > 0" class="results-container">
      <h2 class="results-title">Search Results ({{ results.length }})</h2>
      <MediaGrid
        :items="results"
        :loading="false"
        @media-click="handleMediaClick"
      />
    </div>

    <!-- Initial State (before any search) -->
    <div v-else class="initial-state">
      <svg
        class="initial-icon"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
        />
      </svg>
      <p class="initial-message">Start typing to search for movies and series</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { MediaItem } from '@/types/media'
import MediaGrid from '@/components/media/MediaGrid.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'

interface Props {
  results: MediaItem[]
  loading?: boolean
  hasSearched?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  hasSearched: false
})

interface Emits {
  (e: 'media-click', media: MediaItem): void
}

const emit = defineEmits<Emits>()

// Handle media card click
const handleMediaClick = (media: MediaItem) => {
  emit('media-click', media)
}
</script>

<style scoped>
.search-results {
  width: 100%;
  min-height: 400px;
}

/* Loading State */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  gap: 1.5rem;
}

.loading-text {
  font-size: 1.125rem;
  color: #9ca3af;
  margin: 0;
}

/* Empty State (No Results) */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  padding: 2rem;
  text-align: center;
}

.empty-icon {
  width: 5rem;
  height: 5rem;
  color: #4a4a4a;
  opacity: 0.6;
  margin-bottom: 1.5rem;
}

.empty-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0 0 0.5rem 0;
}

.empty-message {
  font-size: 1rem;
  color: #9ca3af;
  margin: 0;
}

/* Initial State (Before Search) */
.initial-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 400px;
  padding: 2rem;
  text-align: center;
}

.initial-icon {
  width: 4rem;
  height: 4rem;
  color: #6b7280;
  opacity: 0.5;
  margin-bottom: 1rem;
}

.initial-message {
  font-size: 1.125rem;
  color: #6b7280;
  margin: 0;
}

/* Results Container */
.results-container {
  width: 100%;
}

.results-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #ffffff;
  margin: 0 0 2rem 0;
  padding: 0 0.5rem;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .empty-icon {
    width: 4rem;
    height: 4rem;
    margin-bottom: 1rem;
  }

  .empty-title {
    font-size: 1.25rem;
  }

  .empty-message {
    font-size: 0.875rem;
  }

  .initial-icon {
    width: 3rem;
    height: 3rem;
  }

  .initial-message {
    font-size: 1rem;
  }

  .results-title {
    font-size: 1.25rem;
    margin-bottom: 1.5rem;
  }

  .loading-text {
    font-size: 1rem;
  }
}

/* Cinematic theme enhancements */
.empty-state,
.initial-state {
  background: linear-gradient(
    180deg,
    rgba(30, 30, 30, 0.3) 0%,
    rgba(20, 20, 20, 0.5) 100%
  );
  border-radius: 1rem;
}

.results-title {
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
}
</style>
