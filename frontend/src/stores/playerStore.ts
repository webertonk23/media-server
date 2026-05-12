import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { MediaItem } from '@/types/media';
import type { PlayerState, ProgressData, Track } from '@/types/player';
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
    audioTracks: [],
    subtitleTracks: [],
    selectedAudioIndex: 0,
    selectedSubtitleIndex: -1,
    subtitleFont: 'Inter',
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
      const audioTracks: Track[] = [
        { id: '1', index: 0, language: 'pt-BR', label: 'Português (Original)', codec: 'aac' },
        { id: '2', index: 1, language: 'en-US', label: 'English', codec: 'aac' }
      ];
      const subtitleTracks: Track[] = [
        { id: 's1', index: 0, language: 'pt-BR', label: 'Português', codec: 'vtt' },
        { id: 's2', index: 1, language: 'en-US', label: 'English', codec: 'vtt' }
      ];
      playerState.value = {
        ...playerState.value,
        currentTime: progress?.position || 0,
        duration: media.duration || progress?.duration || 0,
        audioTracks,
        subtitleTracks,
        selectedAudioIndex: 0,
        selectedSubtitleIndex: -1,
        isPlaying: false,
        isFullscreen: false,
      };
      console.debug('[Player Store] Initialized player:', {
        mediaId,
        title: media.title,
        duration: playerState.value.duration,
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
      position: Math.round(playerState.value.currentTime),
      duration: Math.round(playerState.value.duration),
      finished,
    };
    try {
      await progressService.saveProgress(currentMedia.value.id, progressData);
      savedProgress.value = progressData;
      lastSaveTime = Date.now();
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
      audioTracks: [],
      subtitleTracks: [],
      selectedAudioIndex: 0,
      selectedSubtitleIndex: -1,
      subtitleFont: 'Inter',
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
