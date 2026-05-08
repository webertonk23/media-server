<template>
  <div class="volume-control">
    <!-- Mute/Unmute Button -->
    <button
      class="control-button volume-button"
      @click="handleToggleMute"
      :aria-label="isMuted ? 'Unmute' : 'Mute'"
    >
      <!-- Volume High Icon -->
      <svg v-if="!isMuted && volume > 0.5" viewBox="0 0 24 24" fill="currentColor">
        <path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM14 3.23v2.06c2.89.86 5 3.54 5 6.71s-2.11 5.85-5 6.71v2.06c4.01-.91 7-4.49 7-8.77s-2.99-7.86-7-8.77z" />
      </svg>
      <!-- Volume Medium Icon -->
      <svg v-else-if="!isMuted && volume > 0" viewBox="0 0 24 24" fill="currentColor">
        <path d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02z" />
      </svg>
      <!-- Volume Muted Icon -->
      <svg v-else viewBox="0 0 24 24" fill="currentColor">
        <path d="M16.5 12c0-1.77-1.02-3.29-2.5-4.03v2.21l2.45 2.45c.03-.2.05-.41.05-.63zm2.5 0c0 .94-.2 1.82-.54 2.64l1.51 1.51C20.63 14.91 21 13.5 21 12c0-4.28-2.99-7.86-7-8.77v2.06c2.89.86 5 3.54 5 6.71zM4.27 3L3 4.27 7.73 9H3v6h4l5 5v-6.73l4.25 4.25c-.67.52-1.42.93-2.25 1.18v2.06c1.38-.31 2.63-.95 3.69-1.81L19.73 21 21 19.73l-9-9L4.27 3zM12 4L9.91 6.09 12 8.18V4z" />
      </svg>
    </button>

    <!-- Volume Slider -->
    <div
      class="volume-slider-container"
      @mousedown="handleMouseDown"
      @mousemove="handleMouseMove"
      @mouseleave="handleMouseLeave"
      ref="sliderRef"
    >
      <div class="volume-track">
        <!-- Volume Fill -->
        <div class="volume-fill" :style="{ width: `${displayVolume}%` }" />

        <!-- Volume Scrubber -->
        <div
          class="volume-scrubber"
          :style="{ left: `${displayVolume}%` }"
        />
      </div>

      <!-- Volume Percentage Tooltip -->
      <div
        v-if="showTooltip"
        class="volume-tooltip"
        :style="{ left: `${tooltipPosition}%` }"
      >
        {{ tooltipVolume }}%
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

/**
 * VolumeControl Component
 * 
 * Provides volume slider with mute button, displays volume percentage,
 * and handles click/drag for volume adjustment.
 * 
 * **Validates: Requirements 5.3, 17.5**
 */

interface Props {
  /** Volume level (0.0 to 1.0) */
  volume: number;
  /** Whether audio is muted */
  isMuted: boolean;
}

interface Emits {
  (e: 'volume', volume: number): void;
  (e: 'toggle-mute'): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

// Refs
const sliderRef = ref<HTMLElement | null>(null);

// State
const isDragging = ref(false);
const showTooltip = ref(false);
const tooltipPosition = ref(0);
const tooltipVolume = ref(0);

/**
 * Display volume as percentage (0-100)
 * 
 * **Validates: Requirement 5.3**
 */
const displayVolume = computed(() => {
  if (props.isMuted) {
    return 0;
  }
  return props.volume * 100;
});

/**
 * Calculate volume from mouse position
 * 
 * @param clientX - Mouse X position
 * @returns Volume level (0.0 to 1.0)
 */
const calculateVolumeFromPosition = (clientX: number): number => {
  if (!sliderRef.value) {
    return props.volume;
  }

  const rect = sliderRef.value.getBoundingClientRect();
  const position = Math.max(0, Math.min(clientX - rect.left, rect.width));
  const percentage = position / rect.width;
  return Math.max(0, Math.min(1, percentage));
};

/**
 * Calculate percentage from mouse position
 * 
 * @param clientX - Mouse X position
 * @returns Percentage (0-100)
 */
const calculatePercentageFromPosition = (clientX: number): number => {
  if (!sliderRef.value) {
    return 0;
  }

  const rect = sliderRef.value.getBoundingClientRect();
  const position = Math.max(0, Math.min(clientX - rect.left, rect.width));
  return (position / rect.width) * 100;
};

/**
 * Handle mute toggle
 * 
 * **Validates: Requirement 5.3**
 */
const handleToggleMute = (): void => {
  emit('toggle-mute');
};

/**
 * Handle mouse down to start dragging
 * 
 * **Validates: Requirement 5.3**
 */
const handleMouseDown = (event: MouseEvent): void => {
  isDragging.value = true;
  handleVolumeChange(event);
};

/**
 * Handle mouse move for tooltip and dragging
 * 
 * **Validates: Requirement 5.3**
 */
const handleMouseMove = (event: MouseEvent): void => {
  // Update tooltip
  const percentage = calculatePercentageFromPosition(event.clientX);
  tooltipPosition.value = percentage;
  tooltipVolume.value = Math.round(percentage);
  showTooltip.value = true;

  // Handle dragging
  if (isDragging.value) {
    handleVolumeChange(event);
  }
};

/**
 * Handle mouse leave to hide tooltip
 */
const handleMouseLeave = (): void => {
  showTooltip.value = false;
};

/**
 * Handle volume change
 * 
 * **Validates: Requirement 5.3**
 */
const handleVolumeChange = (event: MouseEvent): void => {
  const volume = calculateVolumeFromPosition(event.clientX);
  emit('volume', volume);
};

/**
 * Handle global mouse up to stop dragging
 */
const handleGlobalMouseUp = (): void => {
  isDragging.value = false;
};

/**
 * Handle global mouse move for dragging outside the slider
 */
const handleGlobalMouseMove = (event: MouseEvent): void => {
  if (isDragging.value && sliderRef.value) {
    const volume = calculateVolumeFromPosition(event.clientX);
    emit('volume', volume);
  }
};

// Attach global event listeners
if (typeof document !== 'undefined') {
  document.addEventListener('mouseup', handleGlobalMouseUp);
  document.addEventListener('mousemove', handleGlobalMouseMove);
}
</script>

<style scoped>
.volume-control {
  display: flex;
  align-items: center;
  gap: 0.75rem;
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

.volume-slider-container {
  position: relative;
  width: 6rem;
  height: 1.5rem;
  display: flex;
  align-items: center;
  cursor: pointer;
  user-select: none;
}

.volume-track {
  position: relative;
  width: 100%;
  height: 0.25rem;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 0.125rem;
  overflow: visible;
  transition: height 0.2s ease;
}

.volume-slider-container:hover .volume-track {
  height: 0.375rem;
}

.volume-fill {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: linear-gradient(90deg, #3b82f6 0%, #2563eb 100%);
  border-radius: 0.125rem;
  transition: width 0.1s linear;
  box-shadow: 0 0 6px rgba(59, 130, 246, 0.5);
}

.volume-scrubber {
  position: absolute;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 0.75rem;
  height: 0.75rem;
  background: #ffffff;
  border: 2px solid #3b82f6;
  border-radius: 50%;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.3);
  opacity: 0;
  transition: opacity 0.2s ease, transform 0.2s ease;
  pointer-events: none;
}

.volume-slider-container:hover .volume-scrubber {
  opacity: 1;
}

.volume-slider-container:active .volume-scrubber {
  transform: translate(-50%, -50%) scale(1.2);
}

.volume-tooltip {
  position: absolute;
  bottom: 100%;
  transform: translateX(-50%);
  margin-bottom: 0.5rem;
  padding: 0.25rem 0.5rem;
  background: rgba(0, 0, 0, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 0.25rem;
  color: #ffffff;
  font-size: 0.625rem;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
  pointer-events: none;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
  z-index: 30;
}

.volume-tooltip::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  border: 0.25rem solid transparent;
  border-top-color: rgba(0, 0, 0, 0.9);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .volume-control {
    gap: 0.5rem;
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

  .volume-slider-container {
    width: 4rem;
  }

  .volume-track {
    height: 0.375rem;
  }

  .volume-slider-container:hover .volume-track {
    height: 0.5rem;
  }

  .volume-scrubber {
    width: 0.875rem;
    height: 0.875rem;
  }
}

/* Hide volume slider on very small screens */
@media (max-width: 480px) {
  .volume-slider-container {
    display: none;
  }
}

/* Enhanced visual feedback during dragging */
.volume-slider-container:active .volume-track {
  height: 0.5rem;
}

.volume-slider-container:active .volume-fill {
  box-shadow: 0 0 10px rgba(59, 130, 246, 0.8);
}

/* Glassmorphism effect enhancement */
@supports (backdrop-filter: blur(4px)) or (-webkit-backdrop-filter: blur(4px)) {
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
