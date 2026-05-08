<template>
  <div class="video-player-container" @mousemove="handleMouseMove" @mouseleave="handleMouseLeave">
    <!-- Video Element -->
    <video
      ref="videoRef"
      class="video-element"
      :src="streamUrl"
      @click="togglePlayPause"
    />

    <!-- Loading Spinner -->
    <div v-if="isLoading" class="loading-overlay">
      <LoadingSpinner />
    </div>

    <!-- Error Message -->
    <ErrorMessage
      v-if="error"
      :error="error"
      title="Video Playback Error"
      :onRetry="handleRetry"
    />

    <!-- Player Controls -->
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
      @seek="controls.seek"
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
import LoadingSpinner from '@/components/common/LoadingSpinner.vue';
import ErrorMessage from '@/components/common/ErrorMessage.vue';
import PlayerControls from './PlayerControls.vue';

/**
 * VideoPlayer Component
 * 
 * Main video player component that integrates HTML5 video element with
 * custom controls, keyboard shortcuts, and progress tracking.
 * 
 * **Validates: Requirements 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.4, 9.4, 17.1, 17.2, 17.3, 17.4, 17.5**
 */

interface Props {
  /** ULID of the media item to play */
  mediaId: string;
}

const props = defineProps<Props>();

// Store and composables
const playerStore = usePlayerStore();
const videoRef = ref<HTMLVideoElement | null>(null);
const { controls, isLoading, error, toggleMute } = useVideoPlayer(videoRef);

// Enable keyboard shortcuts
useKeyboardShortcuts(controls);

// Player state
const playerState = computed(() => playerStore.playerState);

// Stream URL
const streamUrl = computed(() => getStreamUrl(props.mediaId));

// Controls visibility
const showControls = ref(true);
let hideControlsTimer: ReturnType<typeof setTimeout> | null = null;

/**
 * Handle mouse move to show controls
 * 
 * Shows controls when user moves mouse and sets timer to hide them
 * after 3 seconds of inactivity.
 * 
 * **Validates: Requirements 5.4, 5.5**
 */
const handleMouseMove = (): void => {
  showControls.value = true;

  // Clear existing timer
  if (hideControlsTimer) {
    clearTimeout(hideControlsTimer);
  }

  // Hide controls after 3 seconds of inactivity
  hideControlsTimer = setTimeout(() => {
    if (playerState.value.isPlaying) {
      showControls.value = false;
    }
  }, 3000);
};

/**
 * Handle mouse leave to hide controls
 * 
 * Hides controls when mouse leaves the player area during playback.
 * 
 * **Validates: Requirement 5.4**
 */
const handleMouseLeave = (): void => {
  if (playerState.value.isPlaying) {
    showControls.value = false;
  }
};

/**
 * Toggle play/pause on video click
 * 
 * **Validates: Requirement 5.3**
 */
const togglePlayPause = (): void => {
  if (playerState.value.isPlaying) {
    controls.pause();
  } else {
    controls.play();
  }
};

/**
 * Handle retry after error
 * 
 * Reloads the video element to retry playback.
 * 
 * **Validates: Requirement 9.4**
 */
const handleRetry = (): void => {
  if (videoRef.value) {
    videoRef.value.load();
  }
};

/**
 * Load saved progress and initialize player
 * 
 * Initializes the player store which fetches media details and saved progress.
 * The useVideoPlayer composable will automatically seek to the saved position
 * when the video metadata loads.
 * 
 * **Validates: Requirements 6.4, 17.1**
 */
onMounted(async () => {
  try {
    // Initialize player with media ID (this also loads saved progress)
    await playerStore.initializePlayer(props.mediaId);
    console.debug('[VideoPlayer] Player initialized with saved progress:', playerStore.savedProgress);
  } catch (err: any) {
    console.error('[VideoPlayer] Failed to initialize player:', err);
    // Don't show error to user - just start from beginning
  }
});
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

/* Fullscreen styles */
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
