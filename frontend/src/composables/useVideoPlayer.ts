/**
 * Video Player Composable
 * 
 * Provides video player logic and controls for HTML5 video playback.
 * Handles video events, player controls, and integrates with playerStore
 * for state management and progress tracking.
 * 
 * **Validates: Requirements 5.1, 5.2, 5.3, 6.1, 6.2, 6.5**
 */

import { ref, onMounted, onUnmounted, type Ref } from 'vue';
import { usePlayerStore } from '@/stores/playerStore';
import type { VideoPlayerControls } from '@/types/player';

/**
 * Progress save interval in milliseconds (10 seconds)
 * Progress is automatically saved every 10 seconds during playback
 */
const PROGRESS_SAVE_INTERVAL = 10000;

/**
 * Composable for managing HTML5 video player
 * 
 * Provides controls for video playback and automatically syncs state
 * with the player store. Handles all video events and implements
 * automatic progress tracking.
 * 
 * **Validates: Requirements 5.1, 5.2, 5.3, 6.1, 6.2, 6.5**
 * 
 * @param videoElement - Ref to the HTML5 video element
 * @returns Object containing player controls, state, and error handling
 * 
 * @example
 * ```vue
 * <script setup lang="ts">
 * import { ref } from 'vue';
 * import { useVideoPlayer } from '@/composables/useVideoPlayer';
 * 
 * const videoRef = ref<HTMLVideoElement | null>(null);
 * const { controls, isLoading, error } = useVideoPlayer(videoRef);
 * 
 * // Use controls in template
 * const handlePlay = () => controls.play();
 * const handleSeek = (time: number) => controls.seek(time);
 * </script>
 * 
 * <template>
 *   <video ref="videoRef" />
 *   <button @click="handlePlay">Play</button>
 * </template>
 * ```
 */
export function useVideoPlayer(videoElement: Ref<HTMLVideoElement | null>) {
  const playerStore = usePlayerStore();
  
  // Local state
  const isLoading = ref(true);
  const error = ref<string | null>(null);
  
  // Progress tracking timer
  let progressTimer: ReturnType<typeof setInterval> | null = null;

  /**
   * Play the video
   * 
   * Starts video playback and updates player state.
   * Handles play promise rejection gracefully.
   * 
   * **Validates: Requirement 5.3**
   */
  const play = async (): Promise<void> => {
    if (!videoElement.value) {
      console.warn('[useVideoPlayer] Cannot play: video element not available');
      return;
    }

    try {
      await videoElement.value.play();
      playerStore.updatePlayerState({ isPlaying: true });
      startProgressTracking();
    } catch (err: any) {
      console.error('[useVideoPlayer] Play failed:', err);
      error.value = 'Failed to play video';
    }
  };

  /**
   * Pause the video
   * 
   * Pauses video playback, updates player state, and saves progress immediately.
   * 
   * **Validates: Requirement 5.3**
   */
  const pause = (): void => {
    if (!videoElement.value) {
      console.warn('[useVideoPlayer] Cannot pause: video element not available');
      return;
    }

    videoElement.value.pause();
    playerStore.updatePlayerState({ isPlaying: false });
    stopProgressTracking();
    
    // Save progress immediately when pausing
    playerStore.saveProgressImmediate();
  };

  /**
   * Seek to a specific time position
   * 
   * Updates the video currentTime and saves progress immediately.
   * 
   * **Validates: Requirement 5.3**
   * 
   * @param time - Target time in seconds
   */
  const seek = (time: number): void => {
    if (!videoElement.value) {
      console.warn('[useVideoPlayer] Cannot seek: video element not available');
      return;
    }

    // Clamp time to valid range
    const clampedTime = Math.max(0, Math.min(time, videoElement.value.duration || 0));
    videoElement.value.currentTime = clampedTime;
    
    playerStore.updatePlayerState({ currentTime: clampedTime });
    
    // Save progress immediately after seeking
    playerStore.saveProgressImmediate();
  };

  /**
   * Set volume level
   * 
   * Updates the video volume and mute state.
   * 
   * **Validates: Requirement 5.3**
   * 
   * @param volume - Volume level (0.0 to 1.0)
   */
  const setVolume = (volume: number): void => {
    if (!videoElement.value) {
      console.warn('[useVideoPlayer] Cannot set volume: video element not available');
      return;
    }

    // Clamp volume to valid range
    const clampedVolume = Math.max(0, Math.min(1, volume));
    videoElement.value.volume = clampedVolume;
    
    // Unmute if setting volume > 0
    if (clampedVolume > 0 && videoElement.value.muted) {
      videoElement.value.muted = false;
      playerStore.updatePlayerState({ isMuted: false });
    }
    
    playerStore.updatePlayerState({ volume: clampedVolume });
  };

  /**
   * Toggle mute state
   * 
   * Toggles the video mute state.
   * 
   * **Validates: Requirement 5.3**
   */
  const toggleMute = (): void => {
    if (!videoElement.value) {
      console.warn('[useVideoPlayer] Cannot toggle mute: video element not available');
      return;
    }

    const newMutedState = !videoElement.value.muted;
    videoElement.value.muted = newMutedState;
    playerStore.updatePlayerState({ isMuted: newMutedState });
  };

  /**
   * Toggle fullscreen mode
   * 
   * Enters or exits fullscreen mode for the video element.
   * Handles browser compatibility for fullscreen API.
   * 
   * **Validates: Requirement 5.3**
   */
  const toggleFullscreen = async (): Promise<void> => {
    if (!videoElement.value) {
      console.warn('[useVideoPlayer] Cannot toggle fullscreen: video element not available');
      return;
    }

    try {
      if (!document.fullscreenElement) {
        // Enter fullscreen
        await videoElement.value.requestFullscreen();
        playerStore.updatePlayerState({ isFullscreen: true });
      } else {
        // Exit fullscreen
        await document.exitFullscreen();
        playerStore.updatePlayerState({ isFullscreen: false });
      }
    } catch (err: any) {
      console.error('[useVideoPlayer] Fullscreen toggle failed:', err);
      error.value = 'Failed to toggle fullscreen';
    }
  };

  /**
   * Start progress tracking interval
   * 
   * Starts a timer that saves progress every 10 seconds during playback.
   * 
   * **Validates: Requirement 6.1**
   */
  const startProgressTracking = (): void => {
    // Clear existing timer if any
    stopProgressTracking();
    
    // Save progress every 10 seconds
    progressTimer = setInterval(() => {
      if (videoElement.value && !videoElement.value.paused) {
        playerStore.saveProgress();
      }
    }, PROGRESS_SAVE_INTERVAL);
    
    console.debug('[useVideoPlayer] Progress tracking started');
  };

  /**
   * Stop progress tracking interval
   * 
   * Stops the progress tracking timer.
   * 
   * **Validates: Requirement 6.1**
   */
  const stopProgressTracking = (): void => {
    if (progressTimer) {
      clearInterval(progressTimer);
      progressTimer = null;
      console.debug('[useVideoPlayer] Progress tracking stopped');
    }
  };

  /**
   * Handle video timeupdate event
   * 
   * Updates player state with current time as video plays.
   * 
   * **Validates: Requirement 5.2**
   */
  const handleTimeUpdate = (): void => {
    if (!videoElement.value) return;
    
    playerStore.updatePlayerState({
      currentTime: videoElement.value.currentTime,
    });
  };

  /**
   * Handle video loadedmetadata event
   * 
   * Updates player state with video duration and seeks to saved progress.
   * 
   * **Validates: Requirements 5.2, 6.4**
   */
  const handleLoadedMetadata = (): void => {
    if (!videoElement.value) return;
    
    const duration = videoElement.value.duration;
    playerStore.updatePlayerState({ duration });
    
    // Seek to saved progress if available
    const savedProgress = playerStore.savedProgress;
    if (savedProgress && savedProgress.position > 0) {
      videoElement.value.currentTime = savedProgress.position;
      console.debug('[useVideoPlayer] Resumed from saved progress:', savedProgress.position);
    }
    
    isLoading.value = false;
  };

  /**
   * Handle video ended event
   * 
   * Marks video as finished and saves progress.
   * 
   * **Validates: Requirements 5.2, 6.5**
   */
  const handleEnded = (): void => {
    playerStore.updatePlayerState({ isPlaying: false });
    stopProgressTracking();
    
    // Save final progress with finished flag
    playerStore.saveProgressImmediate();
    
    console.debug('[useVideoPlayer] Video playback ended');
  };

  /**
   * Handle video error event
   * 
   * Captures video errors and updates error state.
   * 
   * **Validates: Requirement 5.2**
   */
  const handleError = (): void => {
    if (!videoElement.value) return;
    
    const videoError = videoElement.value.error;
    let errorMessage = 'Failed to load video';
    
    if (videoError) {
      switch (videoError.code) {
        case MediaError.MEDIA_ERR_ABORTED:
          errorMessage = 'Video loading aborted';
          break;
        case MediaError.MEDIA_ERR_NETWORK:
          errorMessage = 'Network error while loading video';
          break;
        case MediaError.MEDIA_ERR_DECODE:
          errorMessage = 'Video decoding failed';
          break;
        case MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED:
          errorMessage = 'Video format not supported';
          break;
        default:
          errorMessage = 'Unknown video error';
      }
    }
    
    error.value = errorMessage;
    isLoading.value = false;
    
    console.error('[useVideoPlayer] Video error:', errorMessage, videoError);
  };

  /**
   * Handle video play event
   * 
   * Updates player state when video starts playing.
   * 
   * **Validates: Requirement 5.2**
   */
  const handlePlay = (): void => {
    playerStore.updatePlayerState({ isPlaying: true });
    startProgressTracking();
  };

  /**
   * Handle video pause event
   * 
   * Updates player state when video is paused.
   * 
   * **Validates: Requirement 5.2**
   */
  const handlePause = (): void => {
    playerStore.updatePlayerState({ isPlaying: false });
    stopProgressTracking();
  };

  /**
   * Handle video waiting event
   * 
   * Updates loading state when video is buffering.
   * 
   * **Validates: Requirement 5.2**
   */
  const handleWaiting = (): void => {
    isLoading.value = true;
  };

  /**
   * Handle video canplay event
   * 
   * Updates loading state when video is ready to play.
   * 
   * **Validates: Requirement 5.2**
   */
  const handleCanPlay = (): void => {
    isLoading.value = false;
  };

  /**
   * Handle fullscreenchange event
   * 
   * Updates player state when fullscreen mode changes.
   * 
   * **Validates: Requirement 5.2**
   */
  const handleFullscreenChange = (): void => {
    const isFullscreen = !!document.fullscreenElement;
    playerStore.updatePlayerState({ isFullscreen });
  };

  /**
   * Attach event listeners to video element
   * 
   * Sets up all necessary video event listeners.
   * 
   * **Validates: Requirement 5.2**
   */
  const attachEventListeners = (): void => {
    if (!videoElement.value) return;
    
    const video = videoElement.value;
    
    // Playback events
    video.addEventListener('play', handlePlay);
    video.addEventListener('pause', handlePause);
    video.addEventListener('ended', handleEnded);
    video.addEventListener('timeupdate', handleTimeUpdate);
    
    // Loading events
    video.addEventListener('loadedmetadata', handleLoadedMetadata);
    video.addEventListener('waiting', handleWaiting);
    video.addEventListener('canplay', handleCanPlay);
    
    // Error events
    video.addEventListener('error', handleError);
    
    // Fullscreen events
    document.addEventListener('fullscreenchange', handleFullscreenChange);
    
    console.debug('[useVideoPlayer] Event listeners attached');
  };

  /**
   * Detach event listeners from video element
   * 
   * Removes all video event listeners on cleanup.
   * 
   * **Validates: Requirement 5.2**
   */
  const detachEventListeners = (): void => {
    if (!videoElement.value) return;
    
    const video = videoElement.value;
    
    // Playback events
    video.removeEventListener('play', handlePlay);
    video.removeEventListener('pause', handlePause);
    video.removeEventListener('ended', handleEnded);
    video.removeEventListener('timeupdate', handleTimeUpdate);
    
    // Loading events
    video.removeEventListener('loadedmetadata', handleLoadedMetadata);
    video.removeEventListener('waiting', handleWaiting);
    video.removeEventListener('canplay', handleCanPlay);
    
    // Error events
    video.removeEventListener('error', handleError);
    
    // Fullscreen events
    document.removeEventListener('fullscreenchange', handleFullscreenChange);
    
    console.debug('[useVideoPlayer] Event listeners detached');
  };

  /**
   * Initialize video player
   * 
   * Sets up event listeners and initializes player state.
   */
  onMounted(() => {
    if (videoElement.value) {
      attachEventListeners();
      
      // Initialize volume from player store
      const { volume, isMuted } = playerStore.playerState;
      videoElement.value.volume = volume;
      videoElement.value.muted = isMuted;
    }
  });

  /**
   * Cleanup video player
   * 
   * Removes event listeners, stops progress tracking, and saves final progress.
   */
  onUnmounted(() => {
    detachEventListeners();
    stopProgressTracking();
    
    // Save progress one final time before unmounting
    if (videoElement.value && !videoElement.value.paused) {
      playerStore.saveProgressImmediate();
    }
  });

  // Player controls interface
  const controls: VideoPlayerControls = {
    play,
    pause,
    seek,
    setVolume,
    toggleFullscreen,
  };

  return {
    controls,
    isLoading,
    error,
    toggleMute,
  };
}
