import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { MediaItem } from '@/types/media';
import type { PlayerState, ProgressData } from '@/types/player';
import * as mediaService from '@/services/mediaService';
import * as progressService from '@/services/progressService';

const PROGRESS_SAVE_INTERVAL = 10000;
const COMPLETION_THRESHOLD = 0.95;

export const usePlayerStore = defineStore('player', () => {
  const currentMedia = ref<MediaItem | null>(null);
  const playerState = ref<PlayerState>({
    isPlaying: false,
    currentTime: 0,
    duration: 0,
    volume: 1.0,
    isFullscreen: false,
    isMuted: false,
  });
  const savedProgress = ref<ProgressData | null>(null);

  let saveProgressTimer: ReturnType<typeof setTimeout> | null = null;
  let lastSaveTime = 0;

  async function initializePlayer(mediaId: string): Promise<void> {
    try {
      const media = await mediaService.getMediaById(mediaId);
      currentMedia.value = media;

      const progress = await progressService.getProgress(mediaId);
      savedProgress.value = progress;

      if (progress) {
        playerState.value.currentTime = progress.position;
        playerState.value.duration = progress.duration;
      } else {
        playerState.value.currentTime = 0;
        playerState.value.duration = 0;
      }

      playerState.value.isPlaying = false;
      playerState.value.isFullscreen = false;

      console.debug('[Player Store] Initialized player:', {
        mediaId,
        title: media.title,
        hasProgress: !!progress,
        startPosition: progress?.position || 0,
      });
    } catch (error: any) {
      console.error('[Player Store] Failed to initialize player:', error);
      throw error;
    }
  }

  function updatePlayerState(state: Partial<PlayerState>): void {
    playerState.value = {
      ...playerState.value,
      ...state,
    };

    if ('currentTime' in state || 'duration' in state) {
      saveProgress();
    }
  }

  function saveProgress(): void {
    if (!currentMedia.value || playerState.value.duration === 0) {
      return;
    }

    const now = Date.now();
    const timeSinceLastSave = now - lastSaveTime;

    if (saveProgressTimer) {
      clearTimeout(saveProgressTimer);
      saveProgressTimer = null;
    }

    if (timeSinceLastSave >= PROGRESS_SAVE_INTERVAL) {
      performSave();
    } else {
      const delay = PROGRESS_SAVE_INTERVAL - timeSinceLastSave;
      saveProgressTimer = setTimeout(() => {
        performSave();
        saveProgressTimer = null;
      }, delay);
    }
  }

  async function performSave(): Promise<void> {
    if (!currentMedia.value || playerState.value.duration === 0) {
      return;
    }

    const progressPercentage = playerState.value.currentTime / playerState.value.duration;
    const finished = progressPercentage >= COMPLETION_THRESHOLD;

    const progressData: ProgressData = {
      position: playerState.value.currentTime,
      duration: playerState.value.duration,
      finished,
    };

    try {
      await progressService.saveProgress(currentMedia.value.id, progressData);
      savedProgress.value = progressData;
      lastSaveTime = Date.now();

      console.debug('[Player Store] Progress saved:', {
        mediaId: currentMedia.value.id,
        position: progressData.position,
        duration: progressData.duration,
        percentage: Math.round(progressPercentage * 100),
        finished,
      });
    } catch (error: any) {
      console.error('[Player Store] Failed to save progress:', error);
    }
  }

  async function saveProgressImmediate(): Promise<void> {
    if (saveProgressTimer) {
      clearTimeout(saveProgressTimer);
      saveProgressTimer = null;
    }

    await performSave();
  }

  function clearPlayer(): void {
    if (saveProgressTimer) {
      clearTimeout(saveProgressTimer);
      saveProgressTimer = null;
    }

    currentMedia.value = null;
    savedProgress.value = null;
    playerState.value = {
      isPlaying: false,
      currentTime: 0,
      duration: 0,
      volume: 1.0,
      isFullscreen: false,
      isMuted: false,
    };
    lastSaveTime = 0;
  }

  return {
    currentMedia,
    playerState,
    savedProgress,
    initializePlayer,
    updatePlayerState,
    saveProgress,
    saveProgressImmediate,
    clearPlayer,
  };
});

