<template>
  <div class="loading-spinner" :class="sizeClass">
    <div class="spinner-ring">
      <div class="spinner-segment"></div>
      <div class="spinner-segment"></div>
      <div class="spinner-segment"></div>
      <div class="spinner-segment"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  size?: 'small' | 'medium' | 'large'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'medium'
})

const sizeClass = computed(() => `spinner-${props.size}`)
</script>

<style scoped>
.loading-spinner {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.spinner-ring {
  position: relative;
  display: inline-block;
}

.spinner-segment {
  position: absolute;
  border-radius: 50%;
  border-style: solid;
  border-color: transparent;
  animation: spinner-rotate 1.2s cubic-bezier(0.5, 0, 0.5, 1) infinite;
}

.spinner-segment:nth-child(1) {
  border-top-color: #ef4444;
  animation-delay: -0.45s;
}

.spinner-segment:nth-child(2) {
  border-top-color: #f87171;
  animation-delay: -0.3s;
}

.spinner-segment:nth-child(3) {
  border-top-color: #fca5a5;
  animation-delay: -0.15s;
}

.spinner-segment:nth-child(4) {
  border-top-color: rgba(239, 68, 68, 0.3);
}

/* Small size */
.spinner-small .spinner-ring {
  width: 1.5rem;
  height: 1.5rem;
}

.spinner-small .spinner-segment {
  width: 1.5rem;
  height: 1.5rem;
  border-width: 2px;
}

/* Medium size */
.spinner-medium .spinner-ring {
  width: 3rem;
  height: 3rem;
}

.spinner-medium .spinner-segment {
  width: 3rem;
  height: 3rem;
  border-width: 3px;
}

/* Large size */
.spinner-large .spinner-ring {
  width: 4.5rem;
  height: 4.5rem;
}

.spinner-large .spinner-segment {
  width: 4.5rem;
  height: 4.5rem;
  border-width: 4px;
}

@keyframes spinner-rotate {
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
}

/* Cinematic glow effect */
.spinner-segment {
  filter: drop-shadow(0 0 4px rgba(239, 68, 68, 0.5));
}

/* Smooth animation */
@media (prefers-reduced-motion: reduce) {
  .spinner-segment {
    animation-duration: 1.5s;
  }
}
</style>
