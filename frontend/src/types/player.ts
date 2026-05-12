/**
 * Player-related types for video playback and progress tracking
 */
export interface Track {
  index: number;
  id: string;
  language: string;
  label: string;
  codec: string;
}
/**
 * PlayerState interface tracks the current state of the video player
 */
export interface PlayerState {
  isPlaying: boolean;
  currentTime: number;
  duration: number;
  volume: number;
  isFullscreen: boolean;
  isMuted: boolean;
  audioTracks: Track[];
  subtitleTracks: Track[];
  selectedAudioIndex: number;
  selectedSubtitleIndex: number;
  subtitleFont: string;
}
/**
 * ProgressData interface represents saved progress information from the backend
 */
export interface ProgressData {
  position: number;
  duration: number;
  finished: boolean;
}
/**
 * VideoPlayerControls interface defines the interface for controlling video playback
 */
export interface VideoPlayerControls {
  play: () => void;
  pause: () => void;
  seek: (time: number) => void;
  setVolume: (volume: number) => void;
  toggleFullscreen: () => void;
  setAudioTrack: (index: number) => void;
  setSubtitleTrack: (index: number) => void;
  setSubtitleFont: (font: string) => void;
}
