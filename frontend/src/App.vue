<template>
  <div id="app" class="app-container">
    <!-- Error boundary wrapper -->
    <div v-if="hasGlobalError" class="error-boundary-container">
      <ErrorMessage 
        :error="globalError?.message || 'An unexpected error occurred'" 
        :on-retry="handleRetry"
      />
    </div>
    
    <!-- Main application router view -->
    <RouterView v-else v-slot="{ Component, route }">
      <Transition :name="getTransitionName(route)" mode="out-in">
        <Suspense>
          <!-- Main component with error handling -->
          <component :is="Component" :key="route.path" />
          
          <!-- Loading fallback -->
          <template #fallback>
            <div class="loading-container">
              <LoadingSpinner size="large" />
            </div>
          </template>
        </Suspense>
      </Transition>
    </RouterView>
  </div>
</template>

<script setup lang="ts">
import { ref, onErrorCaptured, onMounted } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import LoadingSpinner from './components/common/LoadingSpinner.vue'
import ErrorMessage from './components/common/ErrorMessage.vue'

/**
 * App.vue - Main Application Entry Component
 * 
 * Requirements:
 * - 15.1: Apply global dark theme
 * - 14.5: Handle app-level error boundary
 * 
 * Features:
 * - Global error boundary for uncaught component errors
 * - Route transition animations
 * - Suspense for async component loading
 * - Dark cinematic theme application
 */

// ============================================
// ERROR BOUNDARY STATE
// Requirement: 14.5
// ============================================

const hasGlobalError = ref(false)
const globalError = ref<Error | null>(null)

/**
 * Global error handler
 * Catches errors from child components
 * Requirement: 14.5
 */
onErrorCaptured((err: Error, _instance, info) => {
  console.error('Global error captured:', err, info)
  
  // Set global error state
  hasGlobalError.value = true
  globalError.value = err
  
  // Log to external error tracking service if configured
  // logErrorToService(err, info)
  
  // Prevent error from propagating further
  return false
})

/**
 * Handle retry after error
 * Resets error state and reloads the application
 */
const handleRetry = () => {
  hasGlobalError.value = false
  globalError.value = null
  
  // Reload the current route
  window.location.reload()
}

// ============================================
// ROUTE TRANSITIONS
// Requirement: 15.3
// ============================================

const route = useRoute()

/**
 * Determine transition animation based on route
 * Player routes use fade, others use slide
 */
const getTransitionName = (currentRoute: typeof route) => {
  // No transition for player to avoid visual disruption
  if (currentRoute.meta?.layout === 'player') {
    return 'fade'
  }
  
  // Default slide transition for other routes
  return 'slide'
}

// ============================================
// LIFECYCLE HOOKS
// ============================================

onMounted(() => {
  // Apply dark theme class to body
  // Requirement: 15.1
  document.body.classList.add('dark-theme')
  
  // Set up global error handlers
  window.addEventListener('error', (event) => {
    console.error('Uncaught error:', event.error)
    hasGlobalError.value = true
    globalError.value = event.error
  })
  
  window.addEventListener('unhandledrejection', (event) => {
    console.error('Unhandled promise rejection:', event.reason)
    hasGlobalError.value = true
    globalError.value = new Error(event.reason)
  })
})
</script>

<style scoped>
/* ============================================
   APP CONTAINER
   Requirement: 15.1 - Global dark theme
   ============================================ */

.app-container {
  min-height: 100vh;
  background-color: var(--color-cinema-dark-900);
  color: var(--color-cinema-text-primary);
}

/* ============================================
   LOADING CONTAINER
   Requirement: 11.3
   ============================================ */

.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: var(--color-cinema-dark-900);
}

/* ============================================
   ERROR BOUNDARY CONTAINER
   Requirement: 14.5
   ============================================ */

.error-boundary-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 2rem;
  background-color: var(--color-cinema-dark-900);
}

/* ============================================
   ROUTE TRANSITIONS
   Requirement: 15.3 - Smooth transitions
   ============================================ */

/* Fade transition for player routes */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--transition-fast);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Slide transition for standard routes */
.slide-enter-active,
.slide-leave-active {
  transition: all var(--transition-base);
}

.slide-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.slide-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}
</style>
