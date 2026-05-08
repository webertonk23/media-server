import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import { visualizer } from 'rollup-plugin-visualizer'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    // Bundle analyzer - generates stats.html after build
    visualizer({
      filename: './dist/stats.html',
      open: false,
      gzipSize: true,
      brotliSize: true,
    }) as any,
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:9000',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    },
  },
  build: {
    target: 'es2015',
    minify: 'terser',
    terserOptions: {
      compress: {
        drop_console: true, // Remove console.log in production
        drop_debugger: true,
        pure_funcs: ['console.log', 'console.info'], // Remove specific console methods
      },
    },
    rollupOptions: {
      output: {
        manualChunks(id) {
          // Split vendor chunks for better caching
          if (id.includes('node_modules')) {
            // Vue core libraries (must check before general vue check)
            if (id.includes('@vue/') || (id.includes('/vue/') && !id.includes('vue-router'))) {
              return 'vue-vendor';
            }
            // Vue Router
            if (id.includes('vue-router')) {
              return 'vue-router-vendor';
            }
            // Pinia state management
            if (id.includes('pinia')) {
              return 'pinia-vendor';
            }
            // Axios HTTP client
            if (id.includes('axios')) {
              return 'axios-vendor';
            }
            // Other node_modules
            return 'vendor';
          }
          
          // Split route components into separate chunks
          if (id.includes('/src/pages/')) {
            const pageName = id.split('/src/pages/')[1].split('.vue')[0];
            return `page-${pageName.toLowerCase()}`;
          }
          
          // Split large component groups
          if (id.includes('/src/components/player/')) {
            return 'player-components';
          }
          if (id.includes('/src/components/media/')) {
            return 'media-components';
          }
        },
      },
      // Enable tree-shaking
      treeshake: {
        moduleSideEffects: false,
        propertyReadSideEffects: false,
      },
    },
    chunkSizeWarningLimit: 500, // Warn if chunks exceed 500KB
    // Optimize for Raspberry Pi
    cssCodeSplit: true,
    sourcemap: false,
    // Additional optimizations
    reportCompressedSize: true,
    cssMinify: true,
  },
  optimizeDeps: {
    include: ['vue', 'vue-router', 'pinia', 'axios'],
  },
})
