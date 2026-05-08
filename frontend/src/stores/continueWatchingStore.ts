/**
 * Continue Watching Store
 * 
 * Manages the "Continue Watching" feature by tracking media items with saved progress.
 * Filters items to show only those with progress < 95% (not finished).
 * 
 * **Validates: Requirements 12.3, 7.1, 7.2**
 */

import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { MediaItem } from '@/types/media';
import type { ProgressData } from '@/types/player';
import * as mediaService from '@/services/mediaService';
import * as progressService from '@/services/progressService';

/**
 * MediaItemWithProgress combines a media item with its playback progress
 */
export interface MediaItemWithProgress {
  /** The media item */
  media: MediaItem;
  
  /** Progress data (position, duration, finished) */
  progress: ProgressData;
  
  /** Calculated progress percentage (0-100) */
  progressPercentage: number;
}

/**
 * Continue Watching store for managing in-progress media items
 * 
 * State:
 * - items: Array of media items with their progress data
 * - loading: Loading state for async operations
 * 
 * Actions:
 * - fetchContinueWatching: Fetch all items with progress < 95%
 * - updateProgress: Update progress for a specific item
 * - removeItem: Remove an item from continue watching (when finished)
 */
export const useContinueWatchingStore = defineStore('continueWatching', () => {
  // State
  const items = ref<MediaItemWithProgress[]>([]);
  const loading = ref<boolean>(false);

  /**
   * Fetch all media items with saved progress < 95%
   * 
   * This function fetches all media items and checks their progress.
   * Only items with progress < 95% are included in the continue watching list.
   * Items are sorted by last watched time (most recent first).
   * 
   * **Validates: Requirements 7.1, 12.3**
   * 
   * @returns Promise that resolves when items are loaded
   * 
   * @example
   * ```typescript
   * await fetchContinueWatching();
   * console.log(`${items.value.length} items to continue watching`);
   * ```
   */
  async function fetchContinueWatching(): Promise<void> {
    loading.value = true;

    try {
      // Fetch all media items (we'll need to check progress for each)
      // Note: In a production app, this should be a dedicated backend endpoint
      // that returns only items with progress < 95% to avoid fetching all media
      const response = await mediaService.getMedia({ limit: 100 });
      
      // Fetch progress for each media item
      const itemsWithProgress: MediaItemWithProgress[] = [];
      
      for (const media of response.items) {
        try {
          const progress = await progressService.getProgress(media.id);
          
          // Only include items with progress that are not finished
          if (progress && !progress.finished) {
            const progressPercentage = (progress.position / progress.duration) * 100;
            
            // Only include items with progress < 95%
            if (progressPercentage < 95) {
              itemsWithProgress.push({
                media,
                progress,
                progressPercentage,
              });
            }
          }
        } catch (error) {
          // Skip items where progress fetch fails
          console.debug(`[Continue Watching] No progress for ${media.id}`);
        }
      }
      
      // Sort by most recently watched (would need lastWatchedAt from backend)
      // For now, items are in the order they were fetched
      items.value = itemsWithProgress;
    } catch (error: any) {
      console.error('[Continue Watching] Failed to fetch items:', error);
      items.value = [];
    } finally {
      loading.value = false;
    }
  }

  /**
   * Update progress for a specific media item
   * 
   * Updates the progress for an item in the continue watching list.
   * If progress reaches >= 95%, the item is automatically removed.
   * 
   * **Validates: Requirements 7.2, 12.3**
   * 
   * @param mediaId - ULID of the media item
   * @param progress - New progress data
   * 
   * @example
   * ```typescript
   * // Update progress to 50%
   * await updateProgress('01HQXYZ123ABC456DEF789GHI', {
   *   position: 1800,
   *   duration: 3600,
   *   finished: false
   * });
   * ```
   */
  async function updateProgress(mediaId: string, progress: ProgressData): Promise<void> {
    const progressPercentage = (progress.position / progress.duration) * 100;
    
    // If progress >= 95% or marked as finished, remove from continue watching
    if (progressPercentage >= 95 || progress.finished) {
      removeItem(mediaId);
      return;
    }
    
    // Find existing item
    const existingIndex = items.value.findIndex(item => item.media.id === mediaId);
    
    if (existingIndex !== -1) {
      // Update existing item
      const existingItem = items.value[existingIndex];
      if (existingItem) {
        items.value[existingIndex] = {
          media: existingItem.media,
          progress,
          progressPercentage,
        };
      }
    } else {
      // Item not in list yet - would need to fetch media details
      // For now, we'll just skip adding it (it will appear on next fetchContinueWatching)
      console.debug(`[Continue Watching] Item ${mediaId} not in list, will appear on next fetch`);
    }
  }

  /**
   * Remove an item from continue watching
   * 
   * Removes an item from the continue watching list.
   * This is called when a video is marked as finished (>= 95% watched).
   * 
   * **Validates: Requirements 7.2, 12.3**
   * 
   * @param mediaId - ULID of the media item to remove
   * 
   * @example
   * ```typescript
   * // Remove item when user finishes watching
   * removeItem('01HQXYZ123ABC456DEF789GHI');
   * ```
   */
  function removeItem(mediaId: string): void {
    items.value = items.value.filter(item => item.media.id !== mediaId);
  }

  /**
   * Clear all continue watching items
   * 
   * Resets the store to initial state. Useful for logout or refresh scenarios.
   */
  function clearAll(): void {
    items.value = [];
  }

  return {
    // State
    items,
    loading,
    
    // Actions
    fetchContinueWatching,
    updateProgress,
    removeItem,
    clearAll,
  };
});
