import { ref, onMounted, onUnmounted, type Ref } from 'vue';
import { usePlayerStore } from '@/stores/playerStore';
import type { VideoPlayerControls } from '@/types/player';

const PROGRESS_SAVE_INTERVAL = 10000;

export function useVideoPlayer(videoElement: Ref<HTMLVideoElement | null>) {
  const playerStore = usePlayerStore();

  const isLoading = ref(true);
  const error = ref<string | null>(null);

  let progressTimer: ReturnType<typeof setInterval> | null = null;

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

  const pause = (): void => {
    if (!videoElement.value) {
      console.warn('[useVideoPlayer] Cannot pause: video element not available');
      return;
    }

    videoElement.value.pause();
    playerStore.updatePlayerState({ isPlaying: false });
    stopProgressTracking();

    playerStore.saveProgressImmediate();
  };

  const seek = (time: number): void => {
    if (!videoElement.value) {
      console.warn('[useVideoPlayer] Cannot seek: video element not available');
      return;
    }

    const duration = videoElement.value.duration !== Infinity ? videoElement.value.duration : playerStore.playerState.duration;
    const clampedTime = Math.max(0, Math.min(time, duration || 0));
    videoElement.value.currentTime = clampedTime;

    playerStore.updatePlayerState({ currentTime: clampedTime });

    playerStore.saveProgressImmediate();
  };

  const setVolume = (volume: number): void => {
    if (!videoElement.value) {
      console.warn('[useVideoPlayer] Cannot set volume: video element not available');
      return;
    }

    const clampedVolume = Math.max(0, Math.min(1, volume));
    videoElement.value.volume = clampedVolume;

    if (clampedVolume > 0 && videoElement.value.muted) {
      videoElement.value.muted = false;
      playerStore.updatePlayerState({ isMuted: false });
    }

    playerStore.updatePlayerState({ volume: clampedVolume });
  };

  const toggleMute = (): void => {
    if (!videoElement.value) {
      console.warn('[useVideoPlayer] Cannot toggle mute: video element not available');
      return;
    }

    const newMutedState = !videoElement.value.muted;
    videoElement.value.muted = newMutedState;
    playerStore.updatePlayerState({ isMuted: newMutedState });
  };

  const toggleFullscreen = async (): Promise<void> => {
    if (!videoElement.value) {
      console.warn('[useVideoPlayer] Cannot toggle fullscreen: video element not available');
      return;
    }

    try {
      if (!document.fullscreenElement) {
        await videoElement.value.requestFullscreen();
        playerStore.updatePlayerState({ isFullscreen: true });
      } else {
        await document.exitFullscreen();
        playerStore.updatePlayerState({ isFullscreen: false });
      }
    } catch (err: any) {
      console.error('[useVideoPlayer] Fullscreen toggle failed:', err);
      error.value = 'Failed to toggle fullscreen';
    }
  };

  const startProgressTracking = (): void => {
    stopProgressTracking();

    progressTimer = setInterval(() => {
      if (videoElement.value && !videoElement.value.paused) {
        playerStore.saveProgress();
      }
    }, PROGRESS_SAVE_INTERVAL);

    console.debug('[useVideoPlayer] Progress tracking started');
  };

  const stopProgressTracking = (): void => {
    if (progressTimer) {
      clearInterval(progressTimer);
      progressTimer = null;
      console.debug('[useVideoPlayer] Progress tracking stopped');
    }
  };

  const handleTimeUpdate = (): void => {
    if (!videoElement.value) return;

    let currentTime = videoElement.value.currentTime;
    try {
      const url = new URL(videoElement.value.src, window.location.origin);
      const startParam = url.searchParams.get('start');
      if (startParam) {
        currentTime += parseFloat(startParam);
      }
    } catch (e) {
      // Ignore URL parsing errors
    }

    playerStore.updatePlayerState({
      currentTime: currentTime,
    });
  };

  const handleLoadedMetadata = (): void => {
    if (!videoElement.value) return;

    const duration = videoElement.value.duration;
    if (duration && duration !== Infinity && duration > 0 && playerStore.playerState.duration === 0) {
      playerStore.updatePlayerState({ duration });
      console.debug('[useVideoPlayer] Duration updated from metadata:', duration);
    } else {
      console.debug('[useVideoPlayer] Keeping store duration:', playerStore.playerState.duration);
    }

    const savedProgress = playerStore.savedProgress;
    if (savedProgress && savedProgress.position > 0) {
      try {
        const url = new URL(videoElement.value.src, window.location.origin);
        if (!url.searchParams.has('start')) {
          videoElement.value.currentTime = savedProgress.position;
          console.debug('[useVideoPlayer] Resumed from saved progress:', savedProgress.position);
        } else {
          console.debug('[useVideoPlayer] Stream already started at offset, skipping initial seek');
        }
      } catch (e) {
        videoElement.value.currentTime = savedProgress.position;
      }
    }

    isLoading.value = false;
  };

  const handleEnded = (): void => {
    playerStore.updatePlayerState({ isPlaying: false });
    stopProgressTracking();

    playerStore.saveProgressImmediate();

    console.debug('[useVideoPlayer] Video playback ended');
  };

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

  const handlePlay = (): void => {
    playerStore.updatePlayerState({ isPlaying: true });
    startProgressTracking();
  };

  const handlePause = (): void => {
    playerStore.updatePlayerState({ isPlaying: false });
    stopProgressTracking();
  };

  const handleWaiting = (): void => {
    isLoading.value = true;
  };

  const handleCanPlay = (): void => {
    isLoading.value = false;
  };

  const handleFullscreenChange = (): void => {
    const isFullscreen = !!document.fullscreenElement;
    playerStore.updatePlayerState({ isFullscreen });
  };

  const attachEventListeners = (): void => {
    if (!videoElement.value) return;

    const video = videoElement.value;

    video.addEventListener('play', handlePlay);
    video.addEventListener('pause', handlePause);
    video.addEventListener('ended', handleEnded);
    video.addEventListener('timeupdate', handleTimeUpdate);

    video.addEventListener('loadedmetadata', handleLoadedMetadata);
    video.addEventListener('waiting', handleWaiting);
    video.addEventListener('canplay', handleCanPlay);

    video.addEventListener('error', handleError);

    document.addEventListener('fullscreenchange', handleFullscreenChange);

    console.debug('[useVideoPlayer] Event listeners attached');
  };

  const detachEventListeners = (): void => {
    if (!videoElement.value) return;

    const video = videoElement.value;

    video.removeEventListener('play', handlePlay);
    video.removeEventListener('pause', handlePause);
    video.removeEventListener('ended', handleEnded);
    video.removeEventListener('timeupdate', handleTimeUpdate);

    video.removeEventListener('loadedmetadata', handleLoadedMetadata);
    video.removeEventListener('waiting', handleWaiting);
    video.removeEventListener('canplay', handleCanPlay);

    video.removeEventListener('error', handleError);

    document.removeEventListener('fullscreenchange', handleFullscreenChange);

    console.debug('[useVideoPlayer] Event listeners detached');
  };

  onMounted(() => {
    if (videoElement.value) {
      attachEventListeners();

      const { volume, isMuted } = playerStore.playerState;
      videoElement.value.volume = volume;
      videoElement.value.muted = isMuted;
    }
  });

  onUnmounted(() => {
    detachEventListeners();
    stopProgressTracking();

    if (videoElement.value && !videoElement.value.paused) {
      playerStore.saveProgressImmediate();
    }
  });

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
