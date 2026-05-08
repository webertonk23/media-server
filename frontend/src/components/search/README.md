# Search Components

This directory contains the search functionality components for the Vue Media SPA.

## Components

### SearchBar.vue

A search input component with debouncing and clear functionality.

**Features:**
- Search icon indicator
- Debounced search with 300ms delay (Requirement 8.1)
- Clear button when input has value (Requirement 8.5)
- Emits `search` event after debounce completes

**Props:** None

**Events:**
- `search(query: string)` - Emitted after 300ms debounce when user types

**Usage:**
```vue
<template>
  <SearchBar @search="handleSearch" />
</template>

<script setup lang="ts">
const handleSearch = (query: string) => {
  // Perform search with the debounced query
  console.log('Search query:', query)
}
</script>
```

### SearchResults.vue

Displays search results in a grid format with loading and empty states.

**Features:**
- Displays results using MediaGrid component (Requirement 8.3)
- Shows LoadingSpinner while searching
- Shows "No results found" message when empty (Requirement 8.4)
- Shows initial state before any search
- Emits `media-click` event when user clicks a media item

**Props:**
- `results: MediaItem[]` - Array of search results
- `loading?: boolean` - Whether search is in progress (default: false)
- `hasSearched?: boolean` - Whether a search has been performed (default: false)

**Events:**
- `media-click(media: MediaItem)` - Emitted when user clicks a media card

**Usage:**
```vue
<template>
  <SearchResults
    :results="searchResults"
    :loading="isSearching"
    :has-searched="hasPerformedSearch"
    @media-click="handleMediaClick"
  />
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { MediaItem } from '@/types/media'

const searchResults = ref<MediaItem[]>([])
const isSearching = ref(false)
const hasPerformedSearch = ref(false)

const handleMediaClick = (media: MediaItem) => {
  // Navigate to media detail page
  router.push(`/media/${media.id}`)
}
</script>
```

## Complete Example

Here's a complete example of using both components together:

```vue
<template>
  <div class="search-page">
    <SearchBar @search="performSearch" />
    <SearchResults
      :results="results"
      :loading="loading"
      :has-searched="hasSearched"
      @media-click="navigateToMedia"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import SearchBar from '@/components/search/SearchBar.vue'
import SearchResults from '@/components/search/SearchResults.vue'
import { mediaService } from '@/services/mediaService'
import type { MediaItem } from '@/types/media'

const router = useRouter()
const results = ref<MediaItem[]>([])
const loading = ref(false)
const hasSearched = ref(false)

const performSearch = async (query: string) => {
  if (!query) {
    results.value = []
    hasSearched.value = false
    return
  }

  loading.value = true
  hasSearched.value = true

  try {
    const response = await mediaService.searchMedia(query)
    results.value = response.items
  } catch (error) {
    console.error('Search failed:', error)
    results.value = []
  } finally {
    loading.value = false
  }
}

const navigateToMedia = (media: MediaItem) => {
  router.push(`/media/${media.id}`)
}
</script>

<style scoped>
.search-page {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}
</style>
```

## Requirements Validation

- ✅ **Requirement 8.1**: Debouncing of 300ms for search input
- ✅ **Requirement 8.3**: Display search results in grid format
- ✅ **Requirement 8.4**: Show "No results found" message when empty
- ✅ **Requirement 8.5**: Clear button functionality

## Dependencies

- `useDebounce` composable from `@/composables/useDebounce.ts`
- `LoadingSpinner` component from `@/components/common/LoadingSpinner.vue`
- `MediaGrid` component from `@/components/media/MediaGrid.vue`
- `MediaItem` type from `@/types/media.ts`
