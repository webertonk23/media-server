import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './assets/styles/main.css'
const app = createApp(App)
const pinia = createPinia()
if (import.meta.env.DEV) {
  console.log('🍍 Pinia devtools enabled')
}
app.use(pinia)
app.use(router)
app.config.errorHandler = (err, instance, info) => {
  console.error('Global error handler:', err)
  console.error('Component:', instance)
  console.error('Error info:', info)
}
if (import.meta.env.DEV) {
  app.config.warnHandler = (msg, instance, trace) => {
    console.warn('Vue warning:', msg)
    console.warn('Component:', instance)
    console.warn('Trace:', trace)
  }
}
if (import.meta.env.DEV) {
  app.config.performance = true
}
app.mount('#app')
if (import.meta.env.DEV) {
  console.log('🎬 Media SPA initialized')
  console.log('📦 Vue version:', app.version)
  console.log('🌐 Base URL:', import.meta.env.BASE_URL)
  console.log('🔧 Mode:', import.meta.env.MODE)
}
if (import.meta.hot) {
  import.meta.hot.accept()
}
export { app, pinia, router }
