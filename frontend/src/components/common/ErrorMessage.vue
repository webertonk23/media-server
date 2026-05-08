<template>
  <div class="error-message">
    <div class="error-content">
      <!-- Error Icon -->
      <div class="error-icon">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          />
        </svg>
      </div>

      <!-- Error Message -->
      <div class="error-text">
        <h3 class="error-title">{{ errorTitle }}</h3>
        <p class="error-description">{{ errorMessage }}</p>
      </div>

      <!-- Retry Button (optional) -->
      <button
        v-if="onRetry"
        @click="handleRetry"
        class="retry-button"
      >
        <svg class="retry-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
          />
        </svg>
        <span>Try Again</span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  error?: string | Error | null
  title?: string
  onRetry?: (() => void) | null
}

const props = withDefaults(defineProps<Props>(), {
  error: null,
  title: 'Something went wrong',
  onRetry: null
})

const errorTitle = computed(() => props.title)

const errorMessage = computed(() => {
  if (!props.error) {
    return 'An unexpected error occurred. Please try again.'
  }

  if (typeof props.error === 'string') {
    return props.error
  }

  if (props.error instanceof Error) {
    return props.error.message
  }

  return 'An unexpected error occurred. Please try again.'
})

const handleRetry = () => {
  if (props.onRetry) {
    props.onRetry()
  }
}
</script>

<style scoped>
.error-message {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  min-height: 200px;
}

.error-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  max-width: 500px;
  padding: 2rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 1rem;
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
}

.error-icon {
  width: 4rem;
  height: 4rem;
  margin-bottom: 1.5rem;
  color: #ef4444;
  animation: error-pulse 2s ease-in-out infinite;
}

.error-icon svg {
  width: 100%;
  height: 100%;
  filter: drop-shadow(0 0 8px rgba(239, 68, 68, 0.5));
}

.error-text {
  margin-bottom: 1.5rem;
}

.error-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #ffffff;
  margin: 0 0 0.75rem 0;
}

.error-description {
  font-size: 1rem;
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
  margin: 0;
}

.retry-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: rgba(239, 68, 68, 0.2);
  border: 1px solid rgba(239, 68, 68, 0.5);
  border-radius: 0.5rem;
  color: #ffffff;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.retry-button:hover {
  background: rgba(239, 68, 68, 0.3);
  border-color: #ef4444;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

.retry-button:active {
  transform: translateY(0);
}

.retry-icon {
  width: 1.25rem;
  height: 1.25rem;
}

@keyframes error-pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.7;
  }
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .error-message {
    padding: 1rem;
  }

  .error-content {
    padding: 1.5rem;
  }

  .error-icon {
    width: 3rem;
    height: 3rem;
    margin-bottom: 1rem;
  }

  .error-title {
    font-size: 1.25rem;
  }

  .error-description {
    font-size: 0.875rem;
  }

  .retry-button {
    padding: 0.625rem 1.25rem;
    font-size: 0.875rem;
  }
}
</style>
