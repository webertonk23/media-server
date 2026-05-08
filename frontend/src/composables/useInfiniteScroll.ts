/**
 * Infinite Scroll Composable
 * 
 * Provides intersection observer-based infinite scrolling functionality
 * for lazy loading content when the user scrolls near the bottom of the page.
 * 
 * **Validates: Requirements 2.5, 11.2**
 */

import { ref, onMounted, onUnmounted } from 'vue';

/**
 * Options for configuring the infinite scroll behavior
 */
export interface UseInfiniteScrollOptions {
  /**
   * Root margin for the intersection observer
   * Positive values trigger the callback before the sentinel enters viewport
   * @default '200px' - triggers 200px before sentinel is visible
   */
  rootMargin?: string;

  /**
   * Intersection threshold (0.0 to 1.0)
   * @default 0.0 - triggers as soon as any part of sentinel is visible
   */
  threshold?: number;

  /**
   * Whether infinite scroll is enabled
   * @default true
   */
  enabled?: boolean;
}

/**
 * Composable for implementing infinite scroll with Intersection Observer
 * 
 * Monitors a sentinel element and triggers a callback when it enters the viewport,
 * enabling lazy loading of content as the user scrolls. Useful for implementing
 * infinite scroll in media grids and lists.
 * 
 * **Validates: Requirements 2.5, 11.2**
 * 
 * @param callback - Function to call when sentinel enters viewport (e.g., load more items)
 * @param options - Configuration options for intersection observer
 * @returns Object containing loading state and sentinel element ref
 * 
 * @example
 * ```vue
 * <script setup lang="ts">
 * import { useInfiniteScroll } from '@/composables/useInfiniteScroll';
 * import { useMediaStore } from '@/stores/mediaStore';
 * 
 * const mediaStore = useMediaStore();
 * 
 * const loadMore = async () => {
 *   await mediaStore.fetchMoreMedia();
 * };
 * 
 * const { sentinelRef, isLoading } = useInfiniteScroll(loadMore, {
 *   rootMargin: '200px', // Trigger 200px before reaching bottom
 * });
 * </script>
 * 
 * <template>
 *   <div class="media-grid">
 *     <MediaCard v-for="item in mediaStore.items" :key="item.id" :item="item" />
 *     <!-- Sentinel element for infinite scroll -->
 *     <div ref="sentinelRef" class="sentinel" />
 *   </div>
 * </template>
 * ```
 */
export function useInfiniteScroll(
  callback: () => void | Promise<void>,
  options: UseInfiniteScrollOptions = {}
) {
  const {
    rootMargin = '200px',
    threshold = 0.0,
    enabled = true,
  } = options;

  // Ref to the sentinel element that triggers loading
  const sentinelRef = ref<HTMLElement | null>(null);

  // Loading state to prevent multiple simultaneous loads
  const isLoading = ref(false);

  // Intersection observer instance
  let observer: IntersectionObserver | null = null;

  /**
   * Handle intersection observer callback
   * 
   * Triggers the callback when the sentinel element enters the viewport
   * and loading is not already in progress.
   * 
   * **Validates: Requirement 2.5**
   * 
   * @param entries - Intersection observer entries
   */
  const handleIntersection = async (entries: IntersectionObserverEntry[]) => {
    const [entry] = entries;

    // Only trigger if sentinel is intersecting, not already loading, and enabled
    if (entry && entry.isIntersecting && !isLoading.value && enabled) {
      isLoading.value = true;

      try {
        await callback();
      } catch (error) {
        console.error('[useInfiniteScroll] Error in callback:', error);
      } finally {
        isLoading.value = false;
      }
    }
  };

  /**
   * Initialize the intersection observer
   * 
   * Creates and configures the intersection observer to monitor the sentinel element.
   * 
   * **Validates: Requirement 11.2**
   */
  const initObserver = () => {
    if (!sentinelRef.value) {
      console.warn('[useInfiniteScroll] Sentinel element not found');
      return;
    }

    // Create intersection observer with configured options
    observer = new IntersectionObserver(handleIntersection, {
      root: null, // Use viewport as root
      rootMargin,
      threshold,
    });

    // Start observing the sentinel element
    observer.observe(sentinelRef.value);

    console.debug('[useInfiniteScroll] Observer initialized', {
      rootMargin,
      threshold,
      enabled,
    });
  };

  /**
   * Disconnect the intersection observer
   * 
   * Stops observing the sentinel element and cleans up the observer.
   */
  const disconnectObserver = () => {
    if (observer) {
      observer.disconnect();
      observer = null;
      console.debug('[useInfiniteScroll] Observer disconnected');
    }
  };

  /**
   * Reset the loading state
   * 
   * Manually reset the loading state if needed (e.g., after an error).
   */
  const reset = () => {
    isLoading.value = false;
  };

  /**
   * Initialize observer on mount
   * 
   * Sets up the intersection observer when the component mounts.
   */
  onMounted(() => {
    // Use nextTick to ensure sentinel element is rendered
    setTimeout(() => {
      if (sentinelRef.value) {
        initObserver();
      }
    }, 0);
  });

  /**
   * Cleanup on unmount
   * 
   * Disconnects the observer and cleans up resources.
   */
  onUnmounted(() => {
    disconnectObserver();
  });

  return {
    /**
     * Ref to attach to the sentinel element in the template
     * This element should be placed at the bottom of the scrollable content
     */
    sentinelRef,

    /**
     * Loading state indicating whether content is currently being loaded
     * Use this to show loading indicators or prevent duplicate requests
     */
    isLoading,

    /**
     * Manually reset the loading state
     */
    reset,

    /**
     * Manually disconnect the observer (useful for conditional infinite scroll)
     */
    disconnect: disconnectObserver,

    /**
     * Manually reinitialize the observer (useful after disconnect)
     */
    reconnect: initObserver,
  };
}
