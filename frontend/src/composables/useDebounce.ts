import { ref, watch, type Ref } from 'vue'

/**
 * Composable for debouncing a reactive value
 * 
 * @param value - The reactive value to debounce
 * @param delay - Delay in milliseconds (default: 300ms)
 * @returns Object containing the debounced value and a cancel function
 * 
 * @example
 * ```ts
 * const searchQuery = ref('')
 * const { debouncedValue, cancel } = useDebounce(searchQuery, 300)
 * 
 * watch(debouncedValue, (newValue) => {
 *   // This will only trigger 300ms after the user stops typing
 *   performSearch(newValue)
 * })
 * ```
 */
export function useDebounce<T>(value: Ref<T>, delay: number = 300) {
  const debouncedValue = ref<T>(value.value) as Ref<T>
  let timeoutId: ReturnType<typeof setTimeout> | null = null

  /**
   * Cancel the pending debounced update
   */
  const cancel = () => {
    if (timeoutId !== null) {
      clearTimeout(timeoutId)
      timeoutId = null
    }
  }

  // Watch the input value and debounce updates
  watch(
    value,
    (newValue) => {
      cancel()
      timeoutId = setTimeout(() => {
        debouncedValue.value = newValue
        timeoutId = null
      }, delay)
    },
    { immediate: false }
  )

  return {
    debouncedValue,
    cancel
  }
}
