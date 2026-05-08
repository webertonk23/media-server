/**
 * UI Store
 * 
 * Manages UI state including responsive breakpoints and control visibility.
 * Tracks screen size categories (mobile, tablet, desktop, ultrawide) and player control visibility.
 * 
 * **Validates: Requirements 10.1, 10.2, 10.3, 10.4, 5.4**
 */

import { defineStore } from 'pinia';
import { ref } from 'vue';

/**
 * Responsive breakpoints in pixels
 * 
 * - Mobile: < 768px
 * - Tablet: 768px - 1023px
 * - Desktop: 1024px - 1919px
 * - Ultrawide: >= 1920px
 */
const BREAKPOINTS = {
  MOBILE: 768,
  TABLET: 1024,
  DESKTOP: 1920,
} as const;

/**
 * UI store for managing responsive breakpoints and UI state
 * 
 * State:
 * - isMobile: True when screen width < 768px
 * - isTablet: True when screen width is 768px - 1023px
 * - isDesktop: True when screen width is 1024px - 1919px
 * - isUltrawide: True when screen width >= 1920px
 * - showControls: True when player controls should be visible
 * 
 * Actions:
 * - updateScreenSize: Update breakpoint state based on window width
 * - toggleControls: Toggle player controls visibility
 */
export const useUiStore = defineStore('ui', () => {
  // State
  const isMobile = ref<boolean>(false);
  const isTablet = ref<boolean>(false);
  const isDesktop = ref<boolean>(false);
  const isUltrawide = ref<boolean>(false);
  const showControls = ref<boolean>(true);

  /**
   * Update screen size breakpoint state based on window width
   * 
   * Determines which breakpoint category the current window width falls into
   * and updates the corresponding state flags. Only one breakpoint flag will be true at a time.
   * 
   * **Validates: Requirements 10.1, 10.2, 10.3, 10.4**
   * 
   * @param width - Optional window width in pixels (defaults to window.innerWidth)
   * 
   * @example
   * ```typescript
   * // Update based on current window size
   * updateScreenSize();
   * 
   * // Update with specific width (useful for testing)
   * updateScreenSize(1920);
   * 
   * // Check current breakpoint
   * if (isMobile.value) {
   *   console.log('Mobile layout');
   * }
   * ```
   */
  function updateScreenSize(width?: number): void {
    const windowWidth = width ?? window.innerWidth;

    // Reset all flags
    isMobile.value = false;
    isTablet.value = false;
    isDesktop.value = false;
    isUltrawide.value = false;

    // Set appropriate flag based on width
    if (windowWidth < BREAKPOINTS.MOBILE) {
      // Mobile: < 768px
      isMobile.value = true;
    } else if (windowWidth < BREAKPOINTS.TABLET) {
      // Tablet: 768px - 1023px
      isTablet.value = true;
    } else if (windowWidth < BREAKPOINTS.DESKTOP) {
      // Desktop: 1024px - 1919px
      isDesktop.value = true;
    } else {
      // Ultrawide: >= 1920px
      isUltrawide.value = true;
    }

    console.debug('[UI Store] Screen size updated:', {
      width: windowWidth,
      isMobile: isMobile.value,
      isTablet: isTablet.value,
      isDesktop: isDesktop.value,
      isUltrawide: isUltrawide.value,
    });
  }

  /**
   * Toggle player controls visibility
   * 
   * Toggles the showControls flag between true and false.
   * Used to show/hide video player controls based on user interaction.
   * 
   * **Validates: Requirement 5.4**
   * 
   * @param visible - Optional explicit visibility state (if not provided, toggles current state)
   * 
   * @example
   * ```typescript
   * // Toggle controls visibility
   * toggleControls();
   * 
   * // Explicitly show controls
   * toggleControls(true);
   * 
   * // Explicitly hide controls
   * toggleControls(false);
   * ```
   */
  function toggleControls(visible?: boolean): void {
    if (visible !== undefined) {
      showControls.value = visible;
    } else {
      showControls.value = !showControls.value;
    }

    console.debug('[UI Store] Controls visibility:', showControls.value);
  }

  /**
   * Initialize UI store
   * 
   * Sets up initial screen size state based on current window dimensions.
   * Should be called when the app initializes.
   * 
   * @example
   * ```typescript
   * // Initialize on app mount
   * initialize();
   * ```
   */
  function initialize(): void {
    updateScreenSize();
  }

  return {
    // State
    isMobile,
    isTablet,
    isDesktop,
    isUltrawide,
    showControls,

    // Actions
    updateScreenSize,
    toggleControls,
    initialize,
  };
});
