<template>
  <div class="video-player-container" @mousemove="handleMouseMove" @mouseleave="handleMouseLeave">
    <video
      ref="videoRef"
      class="video-element"
      :src="streamUrl"
      @click="togglePlayPause"
    />

    <div v-if="isLoading" class="loading-overlay">
      <LoadingSpinner />
    </div>

    <ErrorMessage
      v-if="error"
      :error="error"
      title="Video Playback Error"
      :onRetry="handleRetry"
    />

    <PlayerControls
      v-show="showControls && !error"
      :is-playing="playerState.isPlaying"
      :current-time="playerState.currentTime"
      :duration="playerState.duration"
      :volume="playerState.volume"
      :is-muted="playerState.isMuted"
      :is-fullscreen="playerState.isFullscreen"
      @play="controls.play"
      @pause="controls.pause"
      @seek="handleSeek"
      @volume="controls.setVolume"
      @toggle-mute="toggleMute"
      @toggle-fullscreen="controls.toggleFullscreen"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { usePlayerStore } from '@/stores/playerStore';
import { useVideoPlayer } from '@/composables/useVideoPlayer';
import { useKeyboardShortcuts } from '@/composables/useKeyboardShortcuts';
import { getStreamUrl } from '@/services/streamService';
import type { MediaFile } from '@/types/media';
import LoadingSpinner from '@/components/common/LoadingSpinner.vue';
import ErrorMessage from '@/components/common/ErrorMessage.vue';
import PlayerControls from './PlayerControls.vue';

interface Props {
  mediaId: string;
}

const props = defineProps<Props>();

// Store and composables
const playerStore = usePlayerStore();
const videoRef = ref<HTMLVideoElement | null>(null);
const { controls, isLoading, error, toggleMute } = useVideoPlayer(videoRef);

useKeyboardShortcuts(controls);

const playerState = computed(() => playerStore.playerState);

const streamUrl = ref(getStreamUrl(props.mediaId));
const needsTranscoding = computed(() => {
  const media = playerStore.currentMedia;
  if (!media) return false;
  return !media.files?.some((f: MediaFile) => f.path.toLowerCase().endsWith('.mp4') && f.status === 'completed');
});

const showControls = ref(true);
let hideControlsTimer: ReturnType<typeof setTimeout> | null = null;

const handleMouseMove = (): void => {
  showControls.value = true;

  if (hideControlsTimer) {
    clearTimeout(hideControlsTimer);
  }

  hideControlsTimer = setTimeout(() => {
    if (playerState.value.isPlaying) {
      showControls.value = false;
    }
  }, 3000);
};

const handleMouseLeave = (): void => {
  if (playerState.value.isPlaying) {
    showControls.value = false;
  }
};

const togglePlayPause = (): void => {
  if (playerState.value.isPlaying) {
    controls.pause();
  } else {
    controls.play();
  }
};

const handleRetry = (): void => {
  if (videoRef.value) {
    videoRef.value.load();
  }
};

onMounted(async () => {
  try {

    await playerStore.initializePlayer(props.mediaId);

    if (needsTranscoding.value && playerStore.savedProgress) {
      const position = playerStore.savedProgress.position;
      if (position > 0) {
        console.debug('[VideoPlayer] Starting on-the-fly stream from:', position);
        streamUrl.value = getStreamUrl(props.mediaId, position);
      }
    }
    
    console.debug('[VideoPlayer] Player initialized. Needs transcoding:', needsTranscoding.value);
  } catch (err: any) {
    console.error('[VideoPlayer] Failed to initialize player:', err);
  }
});

const handleSeek = (time: number) => {
  if (needsTranscoding.value) {
    console.debug('[VideoPlayer] Reloading stream for seek to:', time);
    streamUrl.value = getStreamUrl(props.mediaId, time);
    setTimeout(() => {
      if (videoRef.value) {
        videoRef.value.play().catch(e => console.warn('[VideoPlayer] Play after seek failed:', e));
      }
    }, 100);
  } else {
    controls.seek(time);
  }
};
</script>

<style scoped>
.video-player-container {
  position: relative;
  width: 100%;
  height: 100%;
  background: #000000;
  overflow: hidden;
  cursor: default;
}

.video-element {
  width: 100%;
  height: 100%;
  object-fit: contain;
  cursor: pointer;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.8);
  z-index: 10;
}

.video-player-container:fullscreen {
  width: 100vw;
  height: 100vh;
}

.video-player-container:-webkit-full-screen {
  width: 100vw;
  height: 100vh;
}

.video-player-container:-moz-full-screen {
  width: 100vw;
  height: 100vh;
}

.video-player-container:-ms-fullscreen {
  width: 100vw;
  height: 100vh;
}
</style>
