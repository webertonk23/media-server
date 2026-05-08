/**
 * Player Store
 * 
 * Manages video player state and progress tracking.
 * Handles player initialization, state updates, and automatic progress saving.
 * 
 * **Validates: Requirements 12.4, 5.2, 6.1, 6.2**
 */

import { defineStore } from 'pinia';
import { ref } from 'vue';
import type { MediaItem } from '@/types/media';
import type { PlayerState, ProgressData } from '@/types/player';
import * as mediaService from '@/services/mediaService';
import * as progressService from '@/services/progressService';

/**
 * Progress save interval in milliseconds (10 seconds)
 * Progress is automatically saved every 10 seconds during playback
 */
const PROGRESS_SAVE_INTERVAL = 10000;

/**
 * Completion threshold percentage (95%)
 * Videos are marked as finished when reaching this percentage
 */
const COMPLETION_THRESHOLD = 0.95;

/**
 * Player store for managing video player state and progress
 * 
 * State:
 * - currentMedia: Currently loaded media item
 * - playerState: Current player state (playing, time, volume, etc.)
 * - savedProgress: Last saved progress data
 * 
 * Actions:
 * - initializePlayer: Load media and saved progress
 * - updatePlayerState: Update player state
 * - saveProgress: Save current progress with throttling
 */
export const usePlayerStore = defineStore('player', () => {
  // State
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

  // Internal state for throttling
  let saveProgressTimer: ReturnType<typeof setTimeout> | null = null;
  let lastSaveTime = 0;

  /**
   * Initialize player with media and load saved progress
   * 
   * Fetches media details and any saved progress from the backend.
   * Sets up the player state with the loaded media and progress position.
   * 
   * **Validates: Requirements 12.4, 5.2**
   * 
   * @param mediaId - ULID of the media item to play
   * @returns Promise that resolves when player is initialized
   * @throws ApiError if media not found or request fails
   * 
   * @example
   * ```typescript
   * // Initialize player for a specific media item
   * await initializePlayer('01HQXYZ123ABC456DEF789GHI');
   * 
   * // Player is now ready with media loaded and progress restored
   * console.log(currentMedia.value?.title);
   * console.log(savedProgress.value?.position);
   * ```
   */
  async function initializePlayer(mediaId: string): Promise<void> {
    try {
      // Fetch media details
      const media = await mediaService.getMediaById(mediaId);
      currentMedia.value = media;

      // Fetch saved progress
      const progress = await progressService.getProgress(mediaId);
      savedProgress.value = progress;

      // Initialize player state with saved progress if available
      if (progress) {
        playerState.value.currentTime = progress.position;
        playerState.value.duration = progress.duration;
      } else {
        // Reset to beginning if no progress
        playerState.value.currentTime = 0;
        playerState.value.duration = 0;
      }

      // Reset playback state
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

  /**
   * Update player state with partial state changes
   * 
   * Merges the provided state changes into the current player state.
   * Automatically triggers progress save when currentTime is updated.
   * 
   * **Validates: Requirements 12.4, 5.2**
   * 
   * @param state - Partial player state to update
   * 
   * @example
   * ```typescript
   * // Update playback state
   * updatePlayerState({ isPlaying: true });
   * 
   * // Update current time (triggers throttled progress save)
   * updatePlayerState({ currentTime: 120 });
   * 
   * // Update multiple properties
   * updatePlayerState({
   *   currentTime: 150,
   *   volume: 0.8,
   *   isMuted: false
   * });
   * ```
   */
  function updatePlayerState(state: Partial<PlayerState>): void {
    // Merge state changes
    playerState.value = {
      ...playerState.value,
      ...state,
    };

    // Trigger throttled progress save if currentTime or duration changed
    if ('currentTime' in state || 'duration' in state) {
      saveProgress();
    }
  }

  /**
   * Save current progress with throttling
   * 
   * Saves the current playback position to the backend.
   * Implements throttling to avoid excessive API calls (max once per 10 seconds).
   * Automatically marks video as finished when reaching 95% completion.
   * 
   * **Validates: Requirements 6.1, 6.2**
   * 
   * @example
   * ```typescript
   * // Save progress (will be throttled to once per 10 seconds)
   * saveProgress();
   * 
   * // Progress is automatically saved during playback
   * // No need to call manually in most cases
   * ```
   */
  function saveProgress(): void {
    if (!currentMedia.value || playerState.value.duration === 0) {
      return;
    }

    const now = Date.now();
    const timeSinceLastSave = now - lastSaveTime;

    // Clear existing timer
    if (saveProgressTimer) {
      clearTimeout(saveProgressTimer);
      saveProgressTimer = null;
    }

    // If enough time has passed, save immediately
    if (timeSinceLastSave >= PROGRESS_SAVE_INTERVAL) {
      performSave();
    } else {
      // Otherwise, schedule save for later
      const delay = PROGRESS_SAVE_INTERVAL - timeSinceLastSave;
      saveProgressTimer = setTimeout(() => {
        performSave();
        saveProgressTimer = null;
      }, delay);
    }
  }

  /**
   * Perform the actual progress save operation
   * 
   * Internal helper that executes the progress save API call.
   * Calculates completion percentage and marks as finished if >= 95%.
   */
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
      // Error is already logged in progressService, no need to throw
    }
  }

  /**
   * Force immediate progress save
   * 
   * Bypasses throttling and saves progress immediately.
   * Useful when user pauses, seeks, or exits the player.
   * 
   * @returns Promise that resolves when progress is saved
   * 
   * @example
   * ```typescript
   * // Save immediately when user pauses
   * await saveProgressImmediate();
   * 
   * // Save before navigating away
   * await saveProgressImmediate();
   * ```
   */
  async function saveProgressImmediate(): Promise<void> {
    // Clear any pending throttled save
    if (saveProgressTimer) {
      clearTimeout(saveProgressTimer);
      saveProgressTimer = null;
    }

    await performSave();
  }

  /**
   * Clear player state and cancel pending saves
   * 
   * Resets the store to initial state and cancels any pending progress saves.
   * Should be called when navigating away from the player.
   * 
   * @example
   * ```typescript
   * // Clean up when leaving player page
   * clearPlayer();
   * ```
   */
  function clearPlayer(): void {
    // Cancel pending save
    if (saveProgressTimer) {
      clearTimeout(saveProgressTimer);
      saveProgressTimer = null;
    }

    // Reset state
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
    // State
    currentMedia,
    playerState,
    savedProgress,

    // Actions
    initializePlayer,
    updatePlayerState,
    saveProgress,
    saveProgressImmediate,
    clearPlayer,
  };
});
