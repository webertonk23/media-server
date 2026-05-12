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
          <svg v-if="!isPlaying" viewBox="0 0 24 24" fill="currentColor">
            <path d="M8 5v14l11-7z" />
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="currentColor">
            <path d="M6 4h4v16H6V4zm8 0h4v16h-4V4z" />
          </svg>
        </button>
        <div class="time-display">
          <span class="current-time">{{ formattedCurrentTime }}</span>
          <span class="time-separator">/</span>
          <span class="total-time">{{ formattedDuration }}</span>
        </div>
      </div>
      <!-- Center Section: Settings (Audio, Subtitles, Fonts) -->
      <div class="controls-section controls-center">
        <!-- Audio Selection -->
        <div class="dropdown">
          <button class="control-button" title="Audio">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 1a3 3 0 0 0-3 3v8a3 3 0 0 0 6 0V4a3 3 0 0 0-3-3z"/>
              <path d="M19 10v2a7 7 0 0 1-14 0v-2M12 19v4M8 23h8"/>
            </svg>
          </button>
          <div class="dropdown-content">
            <div 
              v-for="(track, index) in audioTracks" 
              :key="track.id" 
              class="dropdown-item"
              :class="{ active: selectedAudio === index }"
              @click="emit('select-audio', index)"
            >
              {{ track.label || track.language || `Audio ${index + 1}` }}
            </div>
          </div>
        </div>
        <!-- Subtitle Selection -->
        <div class="dropdown">
          <button class="control-button" title="Subtitles">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
          </button>
          <div class="dropdown-content">
            <div 
              class="dropdown-item"
              :class="{ active: selectedSubtitle === -1 }"
              @click="emit('select-subtitle', -1)"
            >
              Off
            </div>
            <div 
              v-for="(track, index) in subtitleTracks" 
              :key="track.id" 
              class="dropdown-item"
              :class="{ active: selectedSubtitle === index }"
              @click="emit('select-subtitle', index)"
            >
              {{ track.label || track.language || `Subtitle ${index + 1}` }}
            </div>
          </div>
        </div>
        <!-- Font Selection -->
        <div class="dropdown">
          <button class="control-button font-button" title="Font">
            <span>Ag</span>
          </button>
          <div class="dropdown-content">
            <div 
              v-for="font in availableFonts" 
              :key="font" 
              class="dropdown-item"
              :class="{ active: currentFont === font }"
              :style="{ fontFamily: font }"
              @click="emit('select-font', font)"
            >
              {{ font }}
            </div>
          </div>
        </div>
      </div>
      <!-- Right Section: Volume and Fullscreen -->
      <div class="controls-section controls-right">
        <VolumeControl
          :volume="volume"
          :is-muted="isMuted"
          @volume="handleVolume"
          @toggle-mute="handleToggleMute"
        />
        <button
          class="control-button fullscreen-button"
          @click="handleToggleFullscreen"
          :aria-label="isFullscreen ? 'Exit Fullscreen' : 'Enter Fullscreen'"
        >
          <svg v-if="!isFullscreen" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
          </svg>
          <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M9 9V4.5M9 9H4.5M9 9L3.75 3.75M15 9h4.5M15 9V4.5M15 9l5.25-5.25M9 15v4.5M9 15H4.5M9 15l-5.25 5.25M15 15h4.5M15 15v4.5m0-4.5l5.25 5.25" />
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed } from 'vue';
import { formatTime } from '@/utils/formatters';
import type { Track } from '@/types/player';
import ProgressBar from './ProgressBar.vue';
import VolumeControl from './VolumeControl.vue';
interface Props {
  isPlaying: boolean;
  currentTime: number;
  duration: number;
  volume: number;
  isMuted: boolean;
  isFullscreen: boolean;
  audioTracks: Track[];
  subtitleTracks: Track[];
  selectedAudio: number;
  selectedSubtitle: number;
  currentFont: string;
}
interface Emits {
  (e: 'play'): void;
  (e: 'pause'): void;
  (e: 'seek', time: number): void;
  (e: 'volume', volume: number): void;
  (e: 'toggle-mute'): void;
  (e: 'toggle-fullscreen'): void;
  (e: 'select-audio', index: number): void;
  (e: 'select-subtitle', index: number): void;
  (e: 'select-font', font: string): void;
}
const props = defineProps<Props>();
const emit = defineEmits<Emits>();
const availableFonts = ['Inter', 'Roboto', 'Arial', 'Verdana', 'Georgia', 'Times New Roman'];
const formattedCurrentTime = computed(() => formatTime(props.currentTime));
const formattedDuration = computed(() => formatTime(props.duration));
const togglePlayPause = (): void => {
  if (props.isPlaying) {
    emit('pause');
  } else {
    emit('play');
  }
};
const handleSeek = (time: number): void => {
  emit('seek', time);
};
const handleVolume = (volume: number): void => {
  emit('volume', volume);
};
const handleToggleMute = (): void => {
  emit('toggle-mute');
};
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
  padding: 1.5rem;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.9) 0%, rgba(0, 0, 0, 0.7) 50%, transparent 100%);
  backdrop-filter: blur(8px);
  z-index: 20;
}
.controls-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
  margin-top: 1rem;
}
.controls-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}
.controls-center {
  flex: 1;
  justify-content: center;
}
.control-button {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 0.5rem;
  color: #ffffff;
  cursor: pointer;
  transition: all 0.2s ease;
}
.control-button:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.3);
}
.control-button svg {
  width: 1.25rem;
  height: 1.25rem;
}
.font-button span {
  font-size: 1rem;
  font-weight: 600;
}
.dropdown {
  position: relative;
  display: inline-block;
}
.dropdown-content {
  display: none;
  position: absolute;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
  margin-bottom: 0.5rem;
  min-width: 160px;
  background: rgba(20, 20, 20, 0.95);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.75rem;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
  padding: 0.5rem;
  z-index: 50;
}
.dropdown:hover .dropdown-content {
  display: block;
}
.dropdown-item {
  padding: 0.75rem 1rem;
  color: rgba(255, 255, 255, 0.8);
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}
.dropdown-item:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
}
.dropdown-item.active {
  background: var(--primary-color, #e50914);
  color: #ffffff;
}
.time-display {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #ffffff;
}
</style>
