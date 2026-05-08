/**
 * Keyboard Shortcuts Composable
 * 
 * Provides keyboard event handlers for video player controls.
 * Maps keyboard shortcuts to player actions for improved user experience.
 * 
 * **Validates: Requirements 17.1, 17.2, 17.3, 17.4, 17.5, 5.6**
 */

import { onMounted, onUnmounted } from 'vue';
import type { VideoPlayerControls } from '@/types/player';
import { usePlayerStore } from '@/stores/playerStore';

/**
 * Seek interval in seconds for arrow key navigation
 */
const SEEK_INTERVAL = 10;

/**
 * Volume adjustment step for arrow up/down keys
 */
const VOLUME_STEP = 0.1;

/**
 * Composable for managing keyboard shortcuts in the video player
 * 
 * Provides keyboard event handlers that map to player actions:
 * - Space: toggle play/pause
 * - F: toggle fullscreen
 * - Arrow Right: seek forward 10 seconds
 * - Arrow Left: seek backward 10 seconds
 * - Arrow Up/Down: adjust volume
 * 
 * Automatically attaches event listeners on mount and removes them on unmount.
 * 
 * **Validates: Requirements 17.1, 17.2, 17.3, 17.4, 17.5, 5.6**
 * 
 * @param controls - Video player controls interface
 * 
 * @example
 * ```vue
 * <script setup lang="ts">
 * import { ref } from 'vue';
 * import { useVideoPlayer } from '@/composables/useVideoPlayer';
 * import { useKeyboardShortcuts } from '@/composables/useKeyboardShortcuts';
 * 
 * const videoRef = ref<HTMLVideoElement | null>(null);
 * const { controls } = useVideoPlayer(videoRef);
 * 
 * // Enable keyboard shortcuts
 * useKeyboardShortcuts(controls);
 * </script>
 * 
 * <template>
 *   <video ref="videoRef" />
 * </template>
 * ```
 */
export function useKeyboardShortcuts(controls: VideoPlayerControls) {
  const playerStore = usePlayerStore();

  /**
   * Handle keyboard events and map to player actions
   * 
   * Processes keyboard input and triggers appropriate player controls.
   * Prevents default browser behavior for handled keys.
   * 
   * **Validates: Requirements 17.1, 17.2, 17.3, 17.4, 17.5**
   * 
   * @param event - Keyboard event
   */
  const handleKeyDown = (event: KeyboardEvent): void => {
    // Ignore keyboard shortcuts if user is typing in an input field
    const target = event.target as HTMLElement;
    if (
      target.tagName === 'INPUT' ||
      target.tagName === 'TEXTAREA' ||
      target.isContentEditable
    ) {
      return;
    }

    const { key } = event;
    const { playerState } = playerStore;

    switch (key) {
      case ' ':
      case 'Spacebar': // For older browsers
        // Toggle play/pause
        // **Validates: Requirement 17.1**
        event.preventDefault();
        if (playerState.isPlaying) {
          controls.pause();
        } else {
          controls.play();
        }
        console.debug('[useKeyboardShortcuts] Space: toggled play/pause');
        break;

      case 'f':
      case 'F':
        // Toggle fullscreen
        // **Validates: Requirement 17.2**
        event.preventDefault();
        controls.toggleFullscreen();
        console.debug('[useKeyboardShortcuts] F: toggled fullscreen');
        break;

      case 'ArrowRight':
        // Seek forward 10 seconds
        // **Validates: Requirement 17.3**
        event.preventDefault();
        const forwardTime = Math.min(
          playerState.currentTime + SEEK_INTERVAL,
          playerState.duration
        );
        controls.seek(forwardTime);
        console.debug('[useKeyboardShortcuts] Arrow Right: seek forward to', forwardTime);
        break;

      case 'ArrowLeft':
        // Seek backward 10 seconds
        // **Validates: Requirement 17.4**
        event.preventDefault();
        const backwardTime = Math.max(playerState.currentTime - SEEK_INTERVAL, 0);
        controls.seek(backwardTime);
        console.debug('[useKeyboardShortcuts] Arrow Left: seek backward to', backwardTime);
        break;

      case 'ArrowUp':
        // Increase volume
        // **Validates: Requirement 17.5**
        event.preventDefault();
        const increasedVolume = Math.min(playerState.volume + VOLUME_STEP, 1.0);
        controls.setVolume(increasedVolume);
        console.debug('[useKeyboardShortcuts] Arrow Up: volume increased to', increasedVolume);
        break;

      case 'ArrowDown':
        // Decrease volume
        // **Validates: Requirement 17.5**
        event.preventDefault();
        const decreasedVolume = Math.max(playerState.volume - VOLUME_STEP, 0);
        controls.setVolume(decreasedVolume);
        console.debug('[useKeyboardShortcuts] Arrow Down: volume decreased to', decreasedVolume);
        break;

      default:
        // Ignore other keys
        break;
    }
  };

  /**
   * Attach keyboard event listener
   * 
   * Adds the keydown event listener to the document.
   * Called automatically on component mount.
   */
  const attachEventListener = (): void => {
    document.addEventListener('keydown', handleKeyDown);
    console.debug('[useKeyboardShortcuts] Event listener attached');
  };

  /**
   * Detach keyboard event listener
   * 
   * Removes the keydown event listener from the document.
   * Called automatically on component unmount.
   */
  const detachEventListener = (): void => {
    document.removeEventListener('keydown', handleKeyDown);
    console.debug('[useKeyboardShortcuts] Event listener detached');
  };

  /**
   * Initialize keyboard shortcuts
   * 
   * Attaches event listener on component mount.
   */
  onMounted(() => {
    attachEventListener();
  });

  /**
   * Cleanup keyboard shortcuts
   * 
   * Removes event listener on component unmount.
   */
  onUnmounted(() => {
    detachEventListener();
  });

  return {
    // Expose for testing or manual control if needed
    attachEventListener,
    detachEventListener,
  };
}
