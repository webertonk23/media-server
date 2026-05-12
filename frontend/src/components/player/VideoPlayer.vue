<template>
  <div class="video-player-container" @mousemove="handleMouseMove" @mouseleave="handleMouseLeave">
    <button v-show="showControls" class="back-button" @click="goBack" aria-label="Go Back">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M19 12H5M12 19l-7-7 7-7" />
      </svg>
      <span>Voltar</span>
    </button>
    <video
      ref="videoRef"
      class="video-element"
      :src="streamUrl"
      :style="{ fontFamily: playerState.subtitleFont }"
      @click="togglePlayPause"
    >
      <track
        v-for="(track, index) in playerState.subtitleTracks"
        :key="track.id"
        kind="subtitles"
        :src="getSubtitleUrl(track)"
        :srclang="track.language"
        :label="track.label"
        :default="playerState.selectedSubtitleIndex === index"
      />
    </video>
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
      :audio-tracks="playerState.audioTracks"
      :subtitle-tracks="playerState.subtitleTracks"
      :selected-audio="playerState.selectedAudioIndex"
      :selected-subtitle="playerState.selectedSubtitleIndex"
      :current-font="playerState.subtitleFont"
      @play="controls.play"
      @pause="controls.pause"
      @seek="handleSeek"
      @volume="controls.setVolume"
      @toggle-mute="toggleMute"
      @toggle-fullscreen="controls.toggleFullscreen"
      @select-audio="handleAudioSelect"
      @select-subtitle="handleSubtitleSelect"
      @select-font="controls.setSubtitleFont"
    />
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { usePlayerStore } from '@/stores/playerStore';
import { useVideoPlayer } from '@/composables/useVideoPlayer';
import { useKeyboardShortcuts } from '@/composables/useKeyboardShortcuts';
import { getStreamUrl } from '@/services/streamService';
import type { MediaFile } from '@/types/media';
import type { Track } from '@/types/player';
import LoadingSpinner from '@/components/common/LoadingSpinner.vue';
import ErrorMessage from '@/components/common/ErrorMessage.vue';
import PlayerControls from './PlayerControls.vue';
interface Props {
  mediaId: string;
}
const props = defineProps<Props>();
const router = useRouter();
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
const handleMouseMove = () => {
  showControls.value = true;
  if (hideControlsTimer) clearTimeout(hideControlsTimer);
  hideControlsTimer = setTimeout(() => {
    if (playerState.value.isPlaying) showControls.value = false;
  }, 3000);
};
const handleMouseLeave = () => {
  if (playerState.value.isPlaying) showControls.value = false;
};
const togglePlayPause = () => {
  if (playerState.value.isPlaying) controls.pause();
  else controls.play();
};
const handleRetry = () => {
  if (videoRef.value) videoRef.value.load();
};
const goBack = () => {
  router.back();
};
const getSubtitleUrl = (track: Track) => {
  return `/api/stream/subtitle/${props.mediaId}/${track.index}`;
};
const handleAudioSelect = (index: number) => {
  controls.setAudioTrack(index);
  streamUrl.value = getStreamUrl(props.mediaId, playerState.value.currentTime, index);
  setTimeout(() => {
    if (videoRef.value) videoRef.value.play().catch(() => {});
  }, 100);
};

const handleSubtitleSelect = (index: number) => {
  controls.setSubtitleTrack(index);
  if (videoRef.value && videoRef.value.textTracks) {
    Array.from(videoRef.value.textTracks).forEach((track, i) => {
      track.mode = i === index ? 'showing' : 'hidden';
    });
  }
};
onMounted(async () => {
  try {
    await playerStore.initializePlayer(props.mediaId);
    if (needsTranscoding.value && playerStore.savedProgress) {
      const pos = playerStore.savedProgress.position;
      if (pos > 0) streamUrl.value = getStreamUrl(props.mediaId, pos);
    }
  } catch (err) {}
});
const handleSeek = (time: number) => {
  if (needsTranscoding.value) {
    streamUrl.value = getStreamUrl(props.mediaId, time, playerState.value.selectedAudioIndex);
    setTimeout(() => {
      if (videoRef.value) videoRef.value.play().catch(() => {});
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
}
.video-element {
  width: 100%;
  height: 100%;
  object-fit: contain;
  cursor: pointer;
}
.video-element::cue {
  background: rgba(0, 0, 0, 0.7);
  color: #ffffff;
  font-size: 1.2rem;
}
.back-button {
  position: absolute;
  top: 2rem;
  left: 2rem;
  z-index: 30;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.25rem;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 0.75rem;
  color: #ffffff;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}
.back-button:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: translateX(-4px);
}
.back-button svg {
  width: 1.25rem;
  height: 1.25rem;
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
</style>
