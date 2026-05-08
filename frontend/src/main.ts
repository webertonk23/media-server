/**
 * main.ts - Application Entry Point
 * 
 * Requirements:
 * - 12.1: Configure Pinia for state management
 * - 18.2: Initialize Vue app with proper configuration
 * 
 * This file initializes the Vue 3 application with:
 * - Pinia for centralized state management
 * - Vue Router for SPA navigation
 * - Global styles and theme
 * - Error handling configuration
 */

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// Import global styles
// Requirement: 15.1, 15.2, 15.3, 15.4, 15.5
import './assets/styles/main.css'

/**
 * Create Vue application instance
 * Requirement: 18.2
 */
const app = createApp(App)

/**
 * Create and configure Pinia store
 * Requirement: 12.1
 * 
 * Pinia provides centralized state management with:
 * - Type-safe stores with TypeScript
 * - Devtools integration for debugging
 * - Hot module replacement support
 * - Modular store architecture
 */
const pinia = createPinia()

// Enable Pinia devtools in development
if (import.meta.env.DEV) {
  // Pinia devtools are automatically enabled in development
  console.log('🍍 Pinia devtools enabled')
}

/**
 * Register plugins
 * Order matters: Pinia before Router to allow router guards to access stores
 */
app.use(pinia)  // Requirement: 12.1
app.use(router) // Requirement: 13.1

/**
 * Global error handler
 * Requirement: 14.5
 * 
 * Catches errors that escape component error boundaries
 */
app.config.errorHandler = (err, instance, info) => {
  console.error('Global error handler:', err)
  console.error('Component:', instance)
  console.error('Error info:', info)
  
  // Log to external error tracking service if configured
  // Example: Sentry, LogRocket, etc.
  // if (import.meta.env.PROD) {
  //   logErrorToService(err, { component: instance, info })
  // }
}

/**
 * Global warning handler (development only)
 * Helps catch potential issues during development
 */
if (import.meta.env.DEV) {
  app.config.warnHandler = (msg, instance, trace) => {
    console.warn('Vue warning:', msg)
    console.warn('Component:', instance)
    console.warn('Trace:', trace)
  }
}

/**
 * Performance monitoring (development only)
 * Requirement: 13.1, 13.2
 */
if (import.meta.env.DEV) {
  app.config.performance = true
  console.log('📊 Performance monitoring enabled')
}

/**
 * Global properties (if needed)
 * Example: app.config.globalProperties.$api = apiClient
 */
// app.config.globalProperties.$apiBaseUrl = import.meta.env.VITE_API_BASE_URL || 'http://localhost:9000'

/**
 * Mount application to DOM
 * Requirement: 18.2
 * 
 * The app is mounted to #app element in index.html
 * This should be the last operation in this file
 */
app.mount('#app')

/**
 * Log application info in development
 */
if (import.meta.env.DEV) {
  console.log('🎬 Media SPA initialized')
  console.log('📦 Vue version:', app.version)
  console.log('🌐 Base URL:', import.meta.env.BASE_URL)
  console.log('🔧 Mode:', import.meta.env.MODE)
}

/**
 * Hot Module Replacement (HMR) for Pinia stores
 * Preserves store state during development
 */
if (import.meta.hot) {
  import.meta.hot.accept()
  console.log('🔥 HMR enabled')
}

/**
 * Service Worker registration (if needed for PWA)
 * Uncomment to enable offline support
 */
// if ('serviceWorker' in navigator && import.meta.env.PROD) {
//   window.addEventListener('load', () => {
//     navigator.serviceWorker.register('/sw.js').then(
//       (registration) => {
//         console.log('SW registered:', registration)
//       },
//       (error) => {
//         console.error('SW registration failed:', error)
//       }
//     )
//   })
// }

/**
 * Export app instance for testing purposes
 */
export { app, pinia, router }
