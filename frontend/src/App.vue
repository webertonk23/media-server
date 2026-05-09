<template>
  <div id="app" class="app-container">
    <div v-if="hasGlobalError" class="error-boundary-container">
      <ErrorMessage 
        :error="globalError?.message || 'An unexpected error occurred'" 
        :on-retry="handleRetry"
      />
    </div>
    
    <RouterView v-else v-slot="{ Component, route }">
      <Transition :name="getTransitionName(route)" mode="out-in">
        <Suspense>
          <component :is="Component" :key="route.path" />
          
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

const hasGlobalError = ref(false)
const globalError = ref<Error | null>(null)

onErrorCaptured((err: Error, _instance, info) => {
  console.error('Global error captured:', err, info)
  
  hasGlobalError.value = true
  globalError.value = err
  
  return false
})

const handleRetry = () => {
  hasGlobalError.value = false
  globalError.value = null
  window.location.reload()
}

const route = useRoute()

const getTransitionName = (currentRoute: typeof route) => {
  if (currentRoute.meta?.layout === 'player') {
    return 'fade'
  }
  
  return 'slide'
}

onMounted(() => {
  document.body.classList.add('dark-theme')
  
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
.app-container {
  min-height: 100vh;
  background-color: var(--color-cinema-dark-900);
  color: var(--color-cinema-text-primary);
}

.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background-color: var(--color-cinema-dark-900);
}

.error-boundary-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 2rem;
  background-color: var(--color-cinema-dark-900);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--transition-fast);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

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

