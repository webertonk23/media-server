/**
 * Media Store
 * 
 * Centralized state management for media items using Pinia.
 * Handles fetching, searching, and caching of media data from the backend API.
 * 
 * **Validates: Requirements 12.2, 2.1, 4.1, 8.2**
 */

import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { MediaItem, PaginatedResponse } from '@/types/media';
import type { SearchParams } from '@/types/api';
import * as mediaService from '@/services/mediaService';

/**
 * Pagination state interface
 */
interface PaginationState {
  /** Current page number */
  page: number;
  /** Items per page */
  limit: number;
  /** Total number of items */
  total: number;
  /** Whether there are more pages */
  hasMore: boolean;
}

/**
 * Media store for managing media items state
 * 
 * State:
 * - mediaItems: Array of all loaded media items
 * - currentMedia: Currently selected media item
 * - loading: Loading state for async operations
 * - error: Error message if operation fails
 * - pagination: Pagination metadata
 * 
 * Actions:
 * - fetchMedia: Fetch paginated media items with optional filters
 * - fetchMediaById: Fetch a specific media item by ID
 * - searchMedia: Search media items by query string
 * 
 * Getters:
 * - getMediaById: Get a media item from cache by ID
 */
export const useMediaStore = defineStore('media', () => {
  // State
  const mediaItems = ref<MediaItem[]>([]);
  const currentMedia = ref<MediaItem | null>(null);
  const loading = ref<boolean>(false);
  const error = ref<string | null>(null);
  const pagination = ref<PaginationState>({
    page: 1,
    limit: 20,
    total: 0,
    hasMore: false,
  });

  // Getters
  /**
   * Get a media item from the cache by ID
   * 
   * **Validates: Requirement 12.2**
   * 
   * @param id - ULID of the media item
   * @returns MediaItem if found in cache, undefined otherwise
   */
  const getMediaById = computed(() => {
    return (id: string): MediaItem | undefined => {
      return mediaItems.value.find(item => item.id === id);
    };
  });

  // Actions
  /**
   * Fetch paginated media items with optional filtering
   * 
   * Fetches media from the backend API and updates the store state.
   * Supports pagination, search, and type filtering.
   * 
   * **Validates: Requirements 2.1, 12.2**
   * 
   * @param params - Search and pagination parameters
   * @param append - If true, appends results to existing items (for infinite scroll)
   * @returns Promise resolving to paginated response
   * 
   * @example
   * ```typescript
   * // Fetch first page
   * await fetchMedia({ page: 1, limit: 20 });
   * 
   * // Search for movies
   * await fetchMedia({ search: 'Matrix', type: 'movie' });
   * 
   * // Load more (infinite scroll)
   * await fetchMedia({ page: 2 }, true);
   * ```
   */
  async function fetchMedia(params: SearchParams = {}, append: boolean = false): Promise<PaginatedResponse<MediaItem>> {
    loading.value = true;
    error.value = null;

    try {
      const response = await mediaService.getMedia(params);

      // Update media items - either replace or append
      if (append) {
        mediaItems.value = [...mediaItems.value, ...response.items];
      } else {
        mediaItems.value = response.items;
      }

      // Update pagination state
      pagination.value = {
        page: response.page,
        limit: response.limit,
        total: response.total,
        hasMore: response.page * response.limit < response.total,
      };

      return response;
    } catch (err: any) {
      error.value = err.message || 'Erro ao carregar mídia';
      throw err;
    } finally {
      loading.value = false;
    }
  }

  /**
   * Fetch a specific media item by ID
   * 
   * Fetches detailed information for a single media item and sets it as currentMedia.
   * First checks the cache (mediaItems) before making an API call.
   * 
   * **Validates: Requirements 4.1, 12.2**
   * 
   * @param id - ULID of the media item
   * @returns Promise resolving to the media item
   * @throws ApiError if media not found (404) or request fails
   * 
   * @example
   * ```typescript
   * try {
   *   const media = await fetchMediaById('01HQXYZ123ABC456DEF789GHI');
   *   console.log(media.title);
   * } catch (err) {
   *   console.error('Media not found');
   * }
   * ```
   */
  async function fetchMediaById(id: string): Promise<MediaItem> {
    loading.value = true;
    error.value = null;

    try {
      // Check cache first
      const cached = getMediaById.value(id);
      if (cached) {
        currentMedia.value = cached;
        loading.value = false;
        return cached;
      }

      // Fetch from API
      const media = await mediaService.getMediaById(id);
      
      // Update current media
      currentMedia.value = media;

      // Add to cache if not already present
      const existingIndex = mediaItems.value.findIndex(item => item.id === id);
      if (existingIndex === -1) {
        mediaItems.value.push(media);
      } else {
        mediaItems.value[existingIndex] = media;
      }

      return media;
    } catch (err: any) {
      error.value = err.message || 'Erro ao carregar detalhes da mídia';
      currentMedia.value = null;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  /**
   * Search media items by query string
   * 
   * Searches across all media types using the search parameter.
   * Replaces existing mediaItems with search results.
   * 
   * **Validates: Requirements 8.2, 12.2**
   * 
   * @param query - Search query string to match against titles
   * @param page - Page number for pagination (default: 1)
   * @returns Promise resolving to paginated search results
   * 
   * @example
   * ```typescript
   * // Search for "Star Wars"
   * await searchMedia('Star Wars');
   * 
   * // Load next page of results
   * await searchMedia('Star Wars', 2);
   * ```
   */
  async function searchMedia(query: string, page: number = 1): Promise<PaginatedResponse<MediaItem>> {
    loading.value = true;
    error.value = null;

    try {
      const response = await mediaService.searchMedia(query, page);

      // Replace media items with search results
      if (page === 1) {
        mediaItems.value = response.items;
      } else {
        // Append for pagination
        mediaItems.value = [...mediaItems.value, ...response.items];
      }

      // Update pagination state
      pagination.value = {
        page: response.page,
        limit: response.limit,
        total: response.total,
        hasMore: response.page * response.limit < response.total,
      };

      return response;
    } catch (err: any) {
      error.value = err.message || 'Erro ao buscar mídia';
      throw err;
    } finally {
      loading.value = false;
    }
  }

  /**
   * Clear current media selection
   * 
   * Resets currentMedia to null. Useful when navigating away from detail pages.
   */
  function clearCurrentMedia(): void {
    currentMedia.value = null;
  }

  /**
   * Clear all media items and reset state
   * 
   * Resets the store to initial state. Useful for logout or refresh scenarios.
   */
  function clearAll(): void {
    mediaItems.value = [];
    currentMedia.value = null;
    error.value = null;
    pagination.value = {
      page: 1,
      limit: 20,
      total: 0,
      hasMore: false,
    };
  }

  return {
    // State
    mediaItems,
    currentMedia,
    loading,
    error,
    pagination,
    
    // Getters
    getMediaById,
    
    // Actions
    fetchMedia,
    fetchMediaById,
    searchMedia,
    clearCurrentMedia,
    clearAll,
  };
});
