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
    if (!videoElement.value) return;
    try {
      await videoElement.value.play();
      playerStore.updatePlayerState({ isPlaying: true });
      startProgressTracking();
    } catch (err: any) {
      error.value = 'Failed to play video';
    }
  };
  const pause = (): void => {
    if (!videoElement.value) return;
    videoElement.value.pause();
    playerStore.updatePlayerState({ isPlaying: false });
    stopProgressTracking();
    playerStore.saveProgressImmediate();
  };
  const seek = (time: number): void => {
    if (!videoElement.value) return;
    const duration = videoElement.value.duration !== Infinity ? videoElement.value.duration : playerStore.playerState.duration;
    const clampedTime = Math.max(0, Math.min(time, duration || 0));
    videoElement.value.currentTime = clampedTime;
    playerStore.updatePlayerState({ currentTime: clampedTime });
    playerStore.saveProgressImmediate();
  };
  const setVolume = (volume: number): void => {
    if (!videoElement.value) return;
    const clampedVolume = Math.max(0, Math.min(1, volume));
    videoElement.value.volume = clampedVolume;
    if (clampedVolume > 0 && videoElement.value.muted) {
      videoElement.value.muted = false;
      playerStore.updatePlayerState({ isMuted: false });
    }
    playerStore.updatePlayerState({ volume: clampedVolume });
  };
  const toggleMute = (): void => {
    if (!videoElement.value) return;
    const newMutedState = !videoElement.value.muted;
    videoElement.value.muted = newMutedState;
    playerStore.updatePlayerState({ isMuted: newMutedState });
  };
  const toggleFullscreen = async (): Promise<void> => {
    if (!videoElement.value) return;
    try {
      if (!document.fullscreenElement) {
        await videoElement.value.requestFullscreen();
        playerStore.updatePlayerState({ isFullscreen: true });
      } else {
        await document.exitFullscreen();
        playerStore.updatePlayerState({ isFullscreen: false });
      }
    } catch (err: any) {
      error.value = 'Failed to toggle fullscreen';
    }
  };
  const setAudioTrack = (index: number): void => {
    playerStore.updatePlayerState({ selectedAudioIndex: index });
  };
  const setSubtitleTrack = (index: number): void => {
    playerStore.updatePlayerState({ selectedSubtitleIndex: index });
    if (!videoElement.value || !videoElement.value.textTracks) return;
    const tracks = videoElement.value.textTracks;
    Array.from(tracks).forEach((track, i) => {
      track.mode = i === index ? 'showing' : 'disabled';
    });
  };
  const setSubtitleFont = (font: string): void => {
    playerStore.updatePlayerState({ subtitleFont: font });
  };
  const startProgressTracking = (): void => {
    stopProgressTracking();
    progressTimer = setInterval(() => {
      if (videoElement.value && !videoElement.value.paused) {
        playerStore.saveProgress();
      }
    }, PROGRESS_SAVE_INTERVAL);
  };
  const stopProgressTracking = (): void => {
    if (progressTimer) {
      clearInterval(progressTimer);
      progressTimer = null;
    }
  };
  const handleTimeUpdate = (): void => {
    if (!videoElement.value) return;
    let currentTime = videoElement.value.currentTime;
    try {
      const url = new URL(videoElement.value.src, window.location.origin);
      const startParam = url.searchParams.get('start');
      if (startParam) currentTime += parseFloat(startParam);
    } catch (e) {}
    playerStore.updatePlayerState({ currentTime });
  };
  const handleLoadedMetadata = (): void => {
    if (!videoElement.value) return;
    const duration = videoElement.value.duration;
    if (duration && duration !== Infinity && duration > 0 && playerStore.playerState.duration === 0) {
      playerStore.updatePlayerState({ duration });
    }
    const savedProgress = playerStore.savedProgress;
    if (savedProgress && savedProgress.position > 0) {
      try {
        const url = new URL(videoElement.value.src, window.location.origin);
        if (!url.searchParams.has('start')) {
          videoElement.value.currentTime = savedProgress.position;
        }
      } catch (e) {
        videoElement.value.currentTime = savedProgress.position;
      }
    }
    isLoading.value = false;
    setSubtitleTrack(playerStore.playerState.selectedSubtitleIndex);
  };
  const handleEnded = (): void => {
    playerStore.updatePlayerState({ isPlaying: false });
    stopProgressTracking();
    playerStore.saveProgressImmediate();
  };
  const handleError = (): void => {
    if (!videoElement.value) return;
    const videoError = videoElement.value.error;
    let errorMessage = 'Failed to load video';
    if (videoError) {
      switch (videoError.code) {
        case MediaError.MEDIA_ERR_ABORTED: errorMessage = 'Video loading aborted'; break;
        case MediaError.MEDIA_ERR_NETWORK: errorMessage = 'Network error while loading video'; break;
        case MediaError.MEDIA_ERR_DECODE: errorMessage = 'Video decoding failed'; break;
        case MediaError.MEDIA_ERR_SRC_NOT_SUPPORTED: errorMessage = 'Video format not supported'; break;
        default: errorMessage = 'Unknown video error';
      }
    }
    error.value = errorMessage;
    isLoading.value = false;
  };
  const handlePlay = (): void => {
    playerStore.updatePlayerState({ isPlaying: true });
    startProgressTracking();
  };
  const handlePause = (): void => {
    playerStore.updatePlayerState({ isPlaying: false });
    stopProgressTracking();
  };
  const handleWaiting = (): void => { isLoading.value = true; };
  const handleCanPlay = (): void => { isLoading.value = false; };
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
    setAudioTrack,
    setSubtitleTrack,
    setSubtitleFont,
  };
  return {
    controls,
    isLoading,
    error,
    toggleMute,
  };
}
