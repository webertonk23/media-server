import { ref, onMounted, onUnmounted } from 'vue'
/**
 * Breakpoint definitions matching the design requirements
 * - mobile: < 768px
 * - tablet: 768px - 1023px
 * - desktop: 1024px - 1919px
 * - ultrawide: >= 1920px
 */
const BREAKPOINTS = {
  mobile: 768,
  tablet: 1024,
  desktop: 1920,
} as const
/**
 * Composable for responsive breakpoint detection
 * 
 * Detects the current screen size and provides reactive boolean refs
 * for each breakpoint. Updates automatically on window resize with
 * debouncing to avoid excessive updates.
 * 
 * @returns Object containing reactive boolean refs for each breakpoint
 * 
 * @example
 * ```ts
 * const { isMobile, isTablet, isDesktop, isUltrawide } = useMediaQuery()
 * 
 * 
 * const columns = computed(() => {
 *   if (isMobile.value) return 2
 *   if (isTablet.value) return 3
 *   if (isDesktop.value) return 4
 *   return 6 
 * })
 * ```
 */
export function useMediaQuery() {
  const isMobile = ref(false)
  const isTablet = ref(false)
  const isDesktop = ref(false)
  const isUltrawide = ref(false)
  let debounceTimeout: ReturnType<typeof setTimeout> | null = null
  /**
   * Update breakpoint states based on current window width
   */
  const updateBreakpoints = () => {
    const width = window.innerWidth
    isMobile.value = width < BREAKPOINTS.mobile
    isTablet.value = width >= BREAKPOINTS.mobile && width < BREAKPOINTS.tablet
    isDesktop.value = width >= BREAKPOINTS.tablet && width < BREAKPOINTS.desktop
    isUltrawide.value = width >= BREAKPOINTS.desktop
  }
  /**
   * Debounced resize handler to avoid excessive updates
   * Uses 150ms debounce delay for smooth performance
   */
  const handleResize = () => {
    if (debounceTimeout !== null) {
      clearTimeout(debounceTimeout)
    }
    debounceTimeout = setTimeout(() => {
      updateBreakpoints()
      debounceTimeout = null
    }, 150)
  }
  onMounted(() => {
    updateBreakpoints()
    window.addEventListener('resize', handleResize)
  })
  onUnmounted(() => {
    window.removeEventListener('resize', handleResize)
    if (debounceTimeout !== null) {
      clearTimeout(debounceTimeout)
    }
  })
  return {
    isMobile,
    isTablet,
    isDesktop,
    isUltrawide,
  }
}
