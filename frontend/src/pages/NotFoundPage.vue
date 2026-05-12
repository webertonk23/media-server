<script setup lang="ts">
import { useRouter } from 'vue-router'
const router = useRouter()
const goHome = () => {
  router.push({ name: 'home' })
}
</script>
<template>
  <div class="flex flex-col items-center justify-center min-h-screen bg-gray-950 text-white px-4">
    <h1 class="text-6xl font-bold mb-4">404</h1>
    <p class="text-xl text-gray-400 mb-8">Page not found</p>
    <button
      @click="goHome"
      class="px-6 py-3 bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors duration-200"
    >
      Go Home
    </button>
  </div>
</template>
