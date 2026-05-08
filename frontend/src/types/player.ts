/**
 * Player-related types for video playback and progress tracking
 * 
 * **Validates: Requirements 20.1, 5.3, 6.3**
 */

/**
 * PlayerState interface tracks the current state of the video player
 * Used to manage playback controls and UI state
 */
export interface PlayerState {
  /** Whether the video is currently playing */
  isPlaying: boolean;
  
  /** Current playback position in seconds */
  currentTime: number;
  
  /** Total video duration in seconds */
  duration: number;
  
  /** Volume level (0.0 to 1.0) */
  volume: number;
  
  /** Whether the player is in fullscreen mode */
  isFullscreen: boolean;
  
  /** Whether the audio is muted */
  isMuted: boolean;
}

/**
 * ProgressData interface represents saved progress information from the backend
 * Matches the backend progress response structure
 */
export interface ProgressData {
  /** Playback position in seconds */
  position: number;
  
  /** Total video duration in seconds */
  duration: number;
  
  /** Whether the video has been fully watched (>= 95% completion) */
  finished: boolean;
}

/**
 * VideoPlayerControls interface defines the interface for controlling video playback
 * Used by components to interact with the video player
 */
export interface VideoPlayerControls {
  /** Start or resume video playback */
  play: () => void;
  
  /** Pause video playback */
  pause: () => void;
  
  /** Seek to a specific time position in seconds */
  seek: (time: number) => void;
  
  /** Set volume level (0.0 to 1.0) */
  setVolume: (volume: number) => void;
  
  /** Toggle fullscreen mode on/off */
  toggleFullscreen: () => void;
}
