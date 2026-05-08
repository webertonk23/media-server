<template>
  <div
    class="progress-bar-container"
    @mousedown="handleMouseDown"
    @mousemove="handleMouseMove"
    @mouseleave="handleMouseLeave"
    ref="progressBarRef"
  >
    <!-- Buffered Progress (background) -->
    <div class="progress-track">
      <!-- Buffered ranges -->
      <div
        v-for="(range, index) in bufferedRanges"
        :key="index"
        class="buffered-range"
        :style="{ left: `${range.start}%`, width: `${range.width}%` }"
      />

      <!-- Played Progress -->
      <div class="progress-fill" :style="{ width: `${progressPercentage}%` }" />

      <!-- Scrubber -->
      <div
        class="progress-scrubber"
        :style="{ left: `${progressPercentage}%` }"
      />
    </div>

    <!-- Time Tooltip -->
    <div
      v-if="showTooltip"
      class="time-tooltip"
      :style="{ left: `${tooltipPosition}%` }"
    >
      {{ tooltipTime }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { formatTime } from '@/utils/formatters';

/**
 * ProgressBar Component
 * 
 * Displays video playback progress with seekable scrubber,
 * buffered ranges visualization, and time tooltip on hover.
 * 
 * **Validates: Requirements 5.3, 6.1**
 */

interface Props {
  /** Current playback position in seconds */
  currentTime: number;
  /** Total video duration in seconds */
  duration: number;
}

interface Emits {
  (e: 'seek', time: number): void;
}

const props = defineProps<Props>();
const emit = defineEmits<Emits>();

// Refs
const progressBarRef = ref<HTMLElement | null>(null);

// State
const isDragging = ref(false);
const showTooltip = ref(false);
const tooltipPosition = ref(0);
const tooltipTime = ref('00:00');
const bufferedRanges = ref<Array<{ start: number; width: number }>>([]);

/**
 * Calculate progress percentage
 * 
 * **Validates: Requirement 5.3**
 */
const progressPercentage = computed(() => {
  if (!props.duration || props.duration === 0) {
    return 0;
  }
  return (props.currentTime / props.duration) * 100;
});

/**
 * Calculate time from mouse position
 * 
 * @param clientX - Mouse X position
 * @returns Time in seconds
 */
const calculateTimeFromPosition = (clientX: number): number => {
  if (!progressBarRef.value) {
    return 0;
  }

  const rect = progressBarRef.value.getBoundingClientRect();
  const position = Math.max(0, Math.min(clientX - rect.left, rect.width));
  const percentage = position / rect.width;
  return percentage * props.duration;
};

/**
 * Calculate percentage from mouse position
 * 
 * @param clientX - Mouse X position
 * @returns Percentage (0-100)
 */
const calculatePercentageFromPosition = (clientX: number): number => {
  if (!progressBarRef.value) {
    return 0;
  }

  const rect = progressBarRef.value.getBoundingClientRect();
  const position = Math.max(0, Math.min(clientX - rect.left, rect.width));
  return (position / rect.width) * 100;
};

/**
 * Handle mouse down to start dragging
 * 
 * **Validates: Requirement 6.1**
 */
const handleMouseDown = (event: MouseEvent): void => {
  isDragging.value = true;
  handleSeek(event);
};

/**
 * Handle mouse move for tooltip and dragging
 * 
 * **Validates: Requirements 5.3, 6.1**
 */
const handleMouseMove = (event: MouseEvent): void => {
  // Update tooltip
  const time = calculateTimeFromPosition(event.clientX);
  const percentage = calculatePercentageFromPosition(event.clientX);
  
  tooltipPosition.value = percentage;
  tooltipTime.value = formatTime(time);
  showTooltip.value = true;

  // Handle dragging
  if (isDragging.value) {
    handleSeek(event);
  }
};

/**
 * Handle mouse leave to hide tooltip
 */
const handleMouseLeave = (): void => {
  showTooltip.value = false;
};

/**
 * Handle seek operation
 * 
 * **Validates: Requirement 6.1**
 */
const handleSeek = (event: MouseEvent): void => {
  const time = calculateTimeFromPosition(event.clientX);
  emit('seek', time);
};

/**
 * Handle global mouse up to stop dragging
 */
const handleGlobalMouseUp = (): void => {
  isDragging.value = false;
};

/**
 * Handle global mouse move for dragging outside the progress bar
 */
const handleGlobalMouseMove = (event: MouseEvent): void => {
  if (isDragging.value && progressBarRef.value) {
    const time = calculateTimeFromPosition(event.clientX);
    emit('seek', time);
  }
};

/**
 * Update buffered ranges visualization
 * 
 * Note: This is a placeholder implementation. In a real scenario,
 * you would get buffered ranges from the video element's buffered property.
 * 
 * **Validates: Requirement 5.3**
 */
const updateBufferedRanges = (): void => {
  // This would typically be called with the video element's buffered property
  // For now, we'll leave it empty as it requires access to the video element
  // which should be passed from the parent component if needed
  bufferedRanges.value = [];
};

/**
 * Attach global event listeners for dragging
 */
onMounted(() => {
  document.addEventListener('mouseup', handleGlobalMouseUp);
  document.addEventListener('mousemove', handleGlobalMouseMove);
  updateBufferedRanges();
});

/**
 * Remove global event listeners
 */
onUnmounted(() => {
  document.removeEventListener('mouseup', handleGlobalMouseUp);
  document.removeEventListener('mousemove', handleGlobalMouseMove);
});
</script>

<style scoped>
.progress-bar-container {
  position: relative;
  width: 100%;
  height: 1.5rem;
  display: flex;
  align-items: center;
  cursor: pointer;
  user-select: none;
}

.progress-track {
  position: relative;
  width: 100%;
  height: 0.375rem;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 0.25rem;
  overflow: visible;
  transition: height 0.2s ease;
}

.progress-bar-container:hover .progress-track {
  height: 0.5rem;
}

.buffered-range {
  position: absolute;
  top: 0;
  height: 100%;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 0.25rem;
}

.progress-fill {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: linear-gradient(90deg, #3b82f6 0%, #2563eb 100%);
  border-radius: 0.25rem;
  transition: width 0.1s linear;
  box-shadow: 0 0 8px rgba(59, 130, 246, 0.5);
}

.progress-scrubber {
  position: absolute;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 1rem;
  height: 1rem;
  background: #ffffff;
  border: 2px solid #3b82f6;
  border-radius: 50%;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  opacity: 0;
  transition: opacity 0.2s ease, transform 0.2s ease;
  pointer-events: none;
}

.progress-bar-container:hover .progress-scrubber {
  opacity: 1;
}

.progress-bar-container:active .progress-scrubber {
  transform: translate(-50%, -50%) scale(1.2);
}

.time-tooltip {
  position: absolute;
  bottom: 100%;
  transform: translateX(-50%);
  margin-bottom: 0.5rem;
  padding: 0.375rem 0.75rem;
  background: rgba(0, 0, 0, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 0.375rem;
  color: #ffffff;
  font-size: 0.75rem;
  font-weight: 600;
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
  pointer-events: none;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.5);
  z-index: 30;
}

.time-tooltip::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  border: 0.375rem solid transparent;
  border-top-color: rgba(0, 0, 0, 0.9);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .progress-bar-container {
    height: 2rem;
  }

  .progress-track {
    height: 0.5rem;
  }

  .progress-bar-container:hover .progress-track {
    height: 0.625rem;
  }

  .progress-scrubber {
    width: 1.25rem;
    height: 1.25rem;
  }

  .time-tooltip {
    font-size: 0.625rem;
    padding: 0.25rem 0.5rem;
  }
}

/* Enhanced visual feedback during dragging */
.progress-bar-container:active .progress-track {
  height: 0.625rem;
}

.progress-bar-container:active .progress-fill {
  box-shadow: 0 0 12px rgba(59, 130, 246, 0.8);
}
</style>
