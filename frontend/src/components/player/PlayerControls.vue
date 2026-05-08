<template>
  <div class="player-controls">
    <!-- Progress Bar -->
    <ProgressBar
      :current-time="currentTime"
      :duration="duration"
      @seek="handleSeek"
    />

    <!-- Control Buttons Row -->
    <div class="controls-row">
      <!-- Left Section: Play/Pause -->
      <div class="controls-section controls-left">
        <button
          class="control-button play-pause-button"
          @click="togglePlayPause"
          :aria-label="isPlaying ? 'Pause' : 'Play'"
        >
          <!-- Play Icon -->
          <svg v-if="!isPlaying" viewBox="0 0 24 24" fill="currentColor">
            <path d="M8 5v14l11-7z" />
          </svg>
          <!-- Pause Icon -->
          <svg v-else viewBox="0 0 24 24" fill="currentColor">
            <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
          </svg>
        </button>

        <!-- Time Display -->
        <div class="time-display">
          <span class="current-time">{{ formattedCurrentTime }}</span>
          <span class="time-separator">/</span>
          <span class="total-time">{{ formattedDuration }}</span>
        </div>
      </div>

      <!-- Right Section: Volume and Fullscreen -->
      <div class="controls-section controls-right">
        <!-- Volume Control -->
        <VolumeControl
          :volume="volume"
          :is-muted="isMuted"
          @volume="handleVolume"
          @toggle-mute="handleToggleMute"
        />

        <!-- Fullscreen Button -->
        <button
          class="control-button fullscreen-button"
          @click="handleToggleFullscreen"
          :aria-label="isFullscreen ? 'Exit Fullscreen' : 'Enter Fullscreen'"
        >
          <!-- Fullscreen Icon -->
          <svg v-if="!isFullscreen" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4"
            />
          </svg>
          <!-- Exit Fullscreen Icon -->
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 9V4.5M9 9H4.5M9 9L3.75 3.75M15 9h4.5M15 9V4.5M15 9l5.25-5.25M9 15v4.5M9 15H4.5M9 15l-5.25 5.25M15 15h4.5M15 15v4.5m0-4.5l5.25 5.25"
            />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { formatTime } from '@/utils/formatters';
import ProgressBar from './ProgressBar.vue';
import VolumeControl from './VolumeControl.vue';

/**
 * PlayerControls Component
 * 
 * Displays video player controls including play/pause, progress bar,
 * time display, volume control, and fullscreen toggle.
 * 
 * **Validates: Requirements 5.3, 15.2, 15.3**
 */

interface Props {
  /** Whether the video is currently playing */
  isPlaying: boolean;
  /** Current playback position in seconds */
  currentTime: number;
  /** Total video duration in seconds */
  duration: number;
  /** Volume level (0.0 to 1.0) */
  volume: number;
  /** Whether audio is muted */
  isMuted: boolean;
  /** Whether player is in fullscreen mode */
  isFullscreen: boolean;
}

interface Emits {
  (e: 'play'): void;
  (e: 'pause'): void;
  (e: 'seek', time: number): void;
  (e: 'volume', volume: number): void;
  (e: 'toggle-mute'): void;
  (e: 'toggle-fullscreen'): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

/**
 * Formatted current time display
 * 
 * **Validates: Requirement 5.3**
 */
const formattedCurrentTime = computed(() => formatTime(props.currentTime));

/**
 * Formatted duration display
 * 
 * **Validates: Requirement 5.3**
 */
const formattedDuration = computed(() => formatTime(props.duration));

/**
 * Toggle play/pause
 * 
 * **Validates: Requirement 5.3**
 */
const togglePlayPause = (): void => {
  if (props.isPlaying) {
    emit('pause');
  } else {
    emit('play');
  }
};

/**
 * Handle seek event from progress bar
 * 
 * **Validates: Requirement 5.3**
 */
const handleSeek = (time: number): void => {
  emit('seek', time);
};

/**
 * Handle volume change from volume control
 * 
 * **Validates: Requirement 5.3**
 */
const handleVolume = (volume: number): void => {
  emit('volume', volume);
};

/**
 * Handle mute toggle from volume control
 * 
 * **Validates: Requirement 5.3**
 */
const handleToggleMute = (): void => {
  emit('toggle-mute');
};

/**
 * Handle fullscreen toggle
 * 
 * **Validates: Requirement 5.3**
 */
const handleToggleFullscreen = (): void => {
  emit('toggle-fullscreen');
};
</script>

<style scoped>
.player-controls {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 1rem;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.9) 0%, rgba(0, 0, 0, 0.7) 50%, transparent 100%);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  z-index: 20;
  transition: opacity 0.3s ease;
}

.controls-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  margin-top: 0.75rem;
}

.controls-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.controls-left {
  flex: 0 0 auto;
}

.controls-right {
  flex: 0 0 auto;
}

.control-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  padding: 0.5rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 0.5rem;
  color: #ffffff;
  cursor: pointer;
  transition: all 0.2s ease;
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}

.control-button:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.4);
  transform: scale(1.05);
}

.control-button:active {
  transform: scale(0.95);
}

.control-button svg {
  width: 1.5rem;
  height: 1.5rem;
}

.play-pause-button {
  width: 3rem;
  height: 3rem;
}

.play-pause-button svg {
  width: 1.75rem;
  height: 1.75rem;
}

.time-display {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #ffffff;
  font-variant-numeric: tabular-nums;
  user-select: none;
}

.current-time {
  color: #ffffff;
}

.time-separator {
  color: rgba(255, 255, 255, 0.6);
}

.total-time {
  color: rgba(255, 255, 255, 0.8);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .player-controls {
    padding: 0.75rem;
  }

  .controls-row {
    gap: 0.75rem;
    margin-top: 0.5rem;
  }

  .controls-section {
    gap: 0.75rem;
  }

  .control-button {
    width: 2.25rem;
    height: 2.25rem;
    padding: 0.375rem;
  }

  .control-button svg {
    width: 1.25rem;
    height: 1.25rem;
  }

  .play-pause-button {
    width: 2.75rem;
    height: 2.75rem;
  }

  .play-pause-button svg {
    width: 1.5rem;
    height: 1.5rem;
  }

  .time-display {
    font-size: 0.75rem;
    gap: 0.375rem;
  }
}

/* Glassmorphism effect enhancement */
@supports (backdrop-filter: blur(8px)) or (-webkit-backdrop-filter: blur(8px)) {
  .player-controls {
    background: linear-gradient(
      to top,
      rgba(0, 0, 0, 0.8) 0%,
      rgba(0, 0, 0, 0.6) 50%,
      transparent 100%
    );
  }

  .control-button {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.15);
  }

  .control-button:hover {
    background: rgba(255, 255, 255, 0.15);
    border-color: rgba(255, 255, 255, 0.3);
  }
}
</style>
