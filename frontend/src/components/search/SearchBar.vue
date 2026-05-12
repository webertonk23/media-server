<template>
  <div class="search-bar">
    <div class="search-input-wrapper">
      <!-- Search Icon -->
      <svg
        class="search-icon"
        xmlns="http:
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
      <!-- Search Input -->
      <input
        v-model="searchQuery"
        type="text"
        class="search-input"
        placeholder="Search for movies, series..."
        @input="handleInput"
      />
      <!-- Clear Button -->
      <button
        v-if="searchQuery"
        class="clear-button"
        type="button"
        aria-label="Clear search"
        @click="clearSearch"
      >
        <svg
          xmlns="http:
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, watch } from 'vue'
import { useDebounce } from '@/composables/useDebounce'
interface Emits {
  (e: 'search', query: string): void
}
const emit = defineEmits<Emits>()
const searchQuery = ref('')
const { debouncedValue } = useDebounce(searchQuery, 300)
watch(debouncedValue, (newValue) => {
  emit('search', newValue)
})
const handleInput = () => {
}
const clearSearch = () => {
  searchQuery.value = ''
  emit('search', '')
}
</script>
<style scoped>
.search-bar {
  width: 100%;
  max-width: 600px;
  margin: 0 auto;
}
.search-input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  background: rgba(30, 30, 30, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.5rem;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}
.search-input-wrapper:focus-within {
  border-color: rgba(239, 68, 68, 0.5);
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}
.search-icon {
  position: absolute;
  left: 1rem;
  width: 1.25rem;
  height: 1.25rem;
  color: #9ca3af;
  pointer-events: none;
}
.search-input {
  width: 100%;
  padding: 0.875rem 3rem 0.875rem 3rem;
  background: transparent;
  border: none;
  color: #ffffff;
  font-size: 1rem;
  outline: none;
}
.search-input::placeholder {
  color: #6b7280;
}
.clear-button {
  position: absolute;
  right: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2rem;
  height: 2rem;
  padding: 0;
  background: transparent;
  border: none;
  border-radius: 0.375rem;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.2s ease;
}
.clear-button:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
}
.clear-button:active {
  transform: scale(0.95);
}
.clear-button svg {
  width: 1.25rem;
  height: 1.25rem;
}
/* Responsive adjustments */
@media (max-width: 768px) {
  .search-input {
    padding: 0.75rem 2.75rem 0.75rem 2.75rem;
    font-size: 0.875rem;
  }
  .search-icon {
    left: 0.75rem;
    width: 1.125rem;
    height: 1.125rem;
  }
  .clear-button {
    right: 0.5rem;
    width: 1.75rem;
    height: 1.75rem;
  }
  .clear-button svg {
    width: 1.125rem;
    height: 1.125rem;
  }
}
/* Dark cinematic theme enhancements */
.search-input-wrapper {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3),
              0 2px 4px -1px rgba(0, 0, 0, 0.2);
}
.search-input-wrapper:hover {
  border-color: rgba(255, 255, 255, 0.2);
}
</style>
