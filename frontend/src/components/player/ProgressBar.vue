<template>
  <div class="progress-bar-container" @mousedown="handleMouseDown" @mousemove="handleMouseMove"
    @mouseleave="handleMouseLeave" ref="progressBarRef">
    <div class="progress-track">
      <div v-for="(range, index) in bufferedRanges" :key="index" class="buffered-range"
        :style="{ left: `${range.start}%`, width: `${range.width}%` }" />
      <div class="progress-fill" :style="{ width: `${progressPercentage}%` }" />
      <div class="progress-scrubber" :style="{ left: `${progressPercentage}%` }" />
    </div>
    <div v-if="showTooltip" class="time-tooltip" :style="{ left: `${tooltipPosition}%` }">
      {{ tooltipTime }}
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { formatTime } from '@/utils/formatters';
interface Props {
  currentTime: number;
  duration: number;
}
interface Emits {
  (e: 'seek', time: number): void;
}
const props = defineProps<Props>();
const emit = defineEmits<Emits>();
const progressBarRef = ref<HTMLElement | null>(null);
const isDragging = ref(false);
const showTooltip = ref(false);
const tooltipPosition = ref(0);
const tooltipTime = ref('00:00');
const bufferedRanges = ref<Array<{ start: number; width: number }>>([]);
const progressPercentage = computed(() => {
  if (!props.duration || props.duration === 0) return 0;
  return (props.currentTime / props.duration) * 100;
});
const calculateTimeFromPosition = (clientX: number): number => {
  if (!progressBarRef.value) return 0;
  const rect = progressBarRef.value.getBoundingClientRect();
  const position = Math.max(0, Math.min(clientX - rect.left, rect.width));
  const percentage = position / rect.width;
  return percentage * props.duration;
};
const calculatePercentageFromPosition = (clientX: number): number => {
  if (!progressBarRef.value) return 0;
  const rect = progressBarRef.value.getBoundingClientRect();
  const position = Math.max(0, Math.min(clientX - rect.left, rect.width));
  return (position / rect.width) * 100;
};
const handleMouseDown = (event: MouseEvent): void => {
  isDragging.value = true;
  handleSeek(event);
};
const handleMouseMove = (event: MouseEvent): void => {
  const time = calculateTimeFromPosition(event.clientX);
  const percentage = calculatePercentageFromPosition(event.clientX);
  tooltipPosition.value = percentage;
  tooltipTime.value = formatTime(time);
  showTooltip.value = true;
  if (isDragging.value) {
    handleSeek(event);
  }
};
const handleMouseLeave = (): void => {
  showTooltip.value = false;
};
const handleSeek = (event: MouseEvent): void => {
  const time = calculateTimeFromPosition(event.clientX);
  emit('seek', time);
};
const handleGlobalMouseUp = (): void => {
  isDragging.value = false;
};
const handleGlobalMouseMove = (event: MouseEvent): void => {
  if (isDragging.value && progressBarRef.value) {
    const time = calculateTimeFromPosition(event.clientX);
    emit('seek', time);
  }
};
onMounted(() => {
  document.addEventListener('mouseup', handleGlobalMouseUp);
  document.addEventListener('mousemove', handleGlobalMouseMove);
});
onUnmounted(() => {
  document.removeEventListener('mouseup', handleGlobalMouseUp);
  document.removeEventListener('mousemove', handleGlobalMouseMove);
});
</script>
<style scoped>
.progress-bar-container {
  position: relative;
  width: 100%;
  height: 20px;
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 0 4px;
}
.progress-track {
  position: relative;
  width: 100%;
  height: 4px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
  transition: height 0.2s ease;
}
.progress-bar-container:hover .progress-track {
  height: 6px;
}
.buffered-range {
  position: absolute;
  top: 0;
  height: 100%;
  background: rgba(255, 255, 255, 0.15);
  border-radius: 2px;
}
.progress-fill {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: #e50914;
  border-radius: 2px;
  z-index: 2;
}
.progress-scrubber {
  position: absolute;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 14px;
  height: 14px;
  background: #e50914;
  border-radius: 50%;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
  opacity: 0;
  transition: opacity 0.2s ease, transform 0.1s ease;
  z-index: 3;
}
.progress-bar-container:hover .progress-scrubber {
  opacity: 1;
}
.progress-bar-container:active .progress-scrubber {
  transform: translate(-50%, -50%) scale(1.3);
}
.time-tooltip {
  position: absolute;
  bottom: 25px;
  transform: translateX(-50%);
  padding: 4px 8px;
  background: rgba(20, 20, 20, 0.95);
  border-radius: 4px;
  color: white;
  font-size: 14px;
  font-weight: 500;
  pointer-events: none;
  white-space: nowrap;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.5);
  z-index: 100;
}
.time-tooltip::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  border: 6px solid transparent;
  border-top-color: rgba(20, 20, 20, 0.95);
}
</style>
