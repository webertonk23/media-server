import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { MediaItem } from '@/types/media';
import type { ProgressData } from '@/types/player';
import type { ContinueWatchingResponse } from '@/types/api';
import * as progressService from '@/services/progressService';

export interface MediaItemWithProgress {
  media: MediaItem;
  progress: ProgressData;
  progressPercentage: number;
}

export const useContinueWatchingStore = defineStore('continueWatching', () => {
  const items = ref<MediaItemWithProgress[]>([]);
  const loading = ref<boolean>(false);

  async function fetchContinueWatching(): Promise<void> {
    loading.value = true;

    try {
      const response: ContinueWatchingResponse[] = await progressService.getContinueWatching();

      const itemsWithProgress: MediaItemWithProgress[] = response.map(item => {
        const progress: ProgressData = {
          position: item.position,
          duration: item.duration,
          finished: item.finished
        };

        return {
          media: item.media,
          progress,
          progressPercentage: (progress.position / progress.duration) * 100
        };
      });

      items.value = itemsWithProgress;
    } catch (error: any) {
      console.error('[Continue Watching] Failed to fetch items:', error);
      items.value = [];
    } finally {
      loading.value = false;
    }
  }

  async function updateProgress(mediaId: string, progress: ProgressData): Promise<void> {
    const progressPercentage = (progress.position / progress.duration) * 100;

    if (progressPercentage >= 95 || progress.finished) {
      removeItem(mediaId);
      return;
    }

    const existingIndex = items.value.findIndex(item => item.media.id === mediaId);

    if (existingIndex !== -1) {
      const existingItem = items.value[existingIndex];
      if (existingItem) {
        items.value[existingIndex] = {
          media: existingItem.media,
          progress,
          progressPercentage,
        };
      }
    } else {
      console.debug(`[Continue Watching] Item ${mediaId} not in list, will appear on next fetch`);
    }
  }

  function removeItem(mediaId: string): void {
    items.value = items.value.filter(item => item.media.id !== mediaId);
  }

  function clearAll(): void {
    items.value = [];
  }

  return {
    items,
    loading,

    fetchContinueWatching,
    updateProgress,
    removeItem,
    clearAll,
  };
});
