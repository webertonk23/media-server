<template>
  <div class="app-layout" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
    <!-- Sidebar -->
    <aside class="sidebar" :class="{ 'sidebar-open': mobileSidebarOpen }">
      <!-- Logo -->
      <div class="sidebar-logo">
        <router-link to="/" class="logo-link">
          <div class="logo-icon">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M8 5v14l11-7z"/>
            </svg>
          </div>
          <div class="logo-text">
            <span class="logo-media">MEDIA</span><span class="logo-server">SERVER</span>
          </div>
        </router-link>
      </div>

      <!-- Navigation -->
      <nav class="sidebar-nav">
        <div class="nav-section">
          <router-link to="/" class="nav-item" exact-active-class="nav-item-active" title="Início">
            <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
            </svg>
            <span class="nav-label">Início</span>
          </router-link>

          <router-link to="/movies" class="nav-item" active-class="nav-item-active" title="Filmes">
            <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z"/>
            </svg>
            <span class="nav-label">Filmes</span>
          </router-link>

          <router-link to="/series" class="nav-item" active-class="nav-item-active" title="Séries">
            <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="2" y="3" width="20" height="14" rx="2" ry="2"/>
              <line x1="8" y1="21" x2="16" y2="21"/>
              <line x1="12" y1="17" x2="12" y2="21"/>
            </svg>
            <span class="nav-label">Séries</span>
          </router-link>

          <router-link to="/continue" class="nav-item" active-class="nav-item-active" title="Continuar">
            <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <polyline points="12 6 12 12 16 14"/>
            </svg>
            <span class="nav-label">Continuar</span>
          </router-link>

          <router-link to="/recent" class="nav-item" active-class="nav-item-active" title="Recém Adicionados">
            <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 4v16m8-8H4"/>
            </svg>
            <span class="nav-label">Recém Adicionados</span>
          </router-link>

          <router-link to="/watchlist" class="nav-item" active-class="nav-item-active" title="Minha Lista">
            <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
            </svg>
            <span class="nav-label">Minha Lista</span>
          </router-link>
        </div>

        <div class="nav-divider"></div>

        <div class="nav-section nav-section-bottom">
          <router-link to="/settings" class="nav-item" active-class="nav-item-active" title="Configurações">
            <svg class="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="3"/>
              <path stroke-linecap="round" stroke-linejoin="round" d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-4 0v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83-2.83l.06-.06A1.65 1.65 0 004.68 15a1.65 1.65 0 00-1.51-1H3a2 2 0 010-4h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 012.83-2.83l.06.06A1.65 1.65 0 009 4.68a1.65 1.65 0 001-1.51V3a2 2 0 014 0v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 2.83l-.06.06A1.65 1.65 0 0019.4 9a1.65 1.65 0 001.51 1H21a2 2 0 010 4h-.09a1.65 1.65 0 00-1.51 1z"/>
            </svg>
            <span class="nav-label">Configurações</span>
          </router-link>
        </div>
      </nav>
    </aside>

    <!-- Mobile Sidebar Overlay -->
    <div
      v-if="mobileSidebarOpen"
      class="sidebar-overlay"
      @click="mobileSidebarOpen = false"
    ></div>

    <!-- Main Content Area -->
    <div class="main-area">
      <!-- Top Header -->
      <header class="top-header glass-dark">
        <!-- Mobile menu button -->
        <button class="btn btn-icon mobile-menu-btn" @click="mobileSidebarOpen = !mobileSidebarOpen">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="3" y1="6" x2="21" y2="6"/>
            <line x1="3" y1="12" x2="21" y2="12"/>
            <line x1="3" y1="18" x2="21" y2="18"/>
          </svg>
        </button>

        <!-- Search bar -->
        <div class="header-search">
          <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8"/>
            <line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Buscar filmes, séries..."
            class="search-input"
            @keyup.enter="goToSearch"
            @input="onSearchInput"
          />
          <button v-if="searchQuery" class="btn btn-icon search-clear" @click="clearSearch">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        <!-- User avatar -->
        <div class="header-user">
          <div class="user-avatar">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 12c2.7 0 4.8-2.1 4.8-4.8S14.7 2.4 12 2.4 7.2 4.5 7.2 7.2 9.3 12 12 12zm0 2.4c-3.2 0-9.6 1.6-9.6 4.8v2.4h19.2v-2.4c0-3.2-6.4-4.8-9.6-4.8z"/>
            </svg>
          </div>
          <span class="user-name">Weberton</span>
          <svg class="chevron-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="6 9 12 15 18 9"/>
          </svg>
        </div>
      </header>

      <!-- Page Content -->
      <main class="page-content">
        <slot />
      </main>
    </div>

    <!-- Mobile Bottom Navigation -->
    <nav class="mobile-bottom-nav glass-dark">
      <router-link to="/" class="mobile-nav-item" exact-active-class="mobile-nav-item-active">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
        </svg>
        <span>Início</span>
      </router-link>
      <router-link to="/movies" class="mobile-nav-item" active-class="mobile-nav-item-active">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4"/>
        </svg>
        <span>Filmes</span>
      </router-link>
      <router-link to="/search" class="mobile-nav-item" active-class="mobile-nav-item-active">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <span>Buscar</span>
      </router-link>
      <router-link to="/watchlist" class="mobile-nav-item" active-class="mobile-nav-item-active">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
        </svg>
        <span>Minha Lista</span>
      </router-link>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const sidebarCollapsed = ref(false)
const mobileSidebarOpen = ref(false)
const searchQuery = ref('')

// Close mobile sidebar on route change
watch(() => route.path, () => {
  mobileSidebarOpen.value = false
})

const goToSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({ path: '/search', query: { q: searchQuery.value.trim() } })
  }
}

const onSearchInput = () => {
  // Could implement instant search dropdown here
}

const clearSearch = () => {
  searchQuery.value = ''
}
</script>

<style scoped>
/* ==========================================
   APP LAYOUT
   ========================================== */

.app-layout {
  display: flex;
  min-height: 100vh;
  background: var(--color-cinema-dark-900);
  position: relative;
}

/* ==========================================
   SIDEBAR
   ========================================== */

.sidebar {
  width: var(--sidebar-width);
  min-height: 100vh;
  background: var(--color-cinema-dark-800);
  border-right: 1px solid rgba(255, 255, 255, 0.06);
  display: flex;
  flex-direction: column;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 40;
  transition: transform var(--transition-base);
  overflow-y: auto;
  overflow-x: hidden;
  scrollbar-width: none;
}

.sidebar::-webkit-scrollbar { display: none; }

.sidebar-logo {
  padding: 1.5rem 1.25rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  flex-shrink: 0;
}

.logo-link {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  text-decoration: none;
}

.logo-icon {
  width: 2rem;
  height: 2rem;
  background: var(--color-accent);
  border-radius: 0.375rem;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: 0 2px 12px var(--color-accent-glow);
}

.logo-icon svg {
  width: 1.125rem;
  height: 1.125rem;
  color: #fff;
}

.logo-text {
  font-family: var(--font-display);
  font-size: 1rem;
  font-weight: 700;
  letter-spacing: 0.05em;
  white-space: nowrap;
}

.logo-media {
  color: #fff;
}

.logo-server {
  color: var(--color-accent);
}

/* Navigation */
.sidebar-nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 1rem 0;
}

.nav-section {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
  padding: 0 0.75rem;
}

.nav-section-bottom {
  margin-top: auto;
  padding-top: 1rem;
}

.nav-divider {
  margin: 1rem 1.25rem;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.875rem;
  padding: 0.75rem 0.875rem;
  border-radius: 0.5rem;
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all var(--transition-fast);
  white-space: nowrap;
  overflow: hidden;
  position: relative;
}

.nav-item:hover {
  color: var(--color-text-primary);
  background: rgba(255, 255, 255, 0.06);
}

.nav-item-active {
  color: var(--color-text-primary) !important;
  background: rgba(229, 9, 20, 0.12) !important;
}

.nav-item-active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 60%;
  background: var(--color-accent);
  border-radius: 0 2px 2px 0;
}

.nav-icon {
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
  stroke-width: 1.75;
}

.nav-label {
  font-size: 0.875rem;
}

/* Sidebar Overlay for mobile */
.sidebar-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 39;
  backdrop-filter: blur(4px);
}

/* ==========================================
   MAIN AREA
   ========================================== */

.main-area {
  flex: 1;
  min-width: 0;
  margin-left: var(--sidebar-width);
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* ==========================================
   TOP HEADER
   ========================================== */

.top-header {
  position: sticky;
  top: 0;
  z-index: 30;
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 0.875rem 2rem;
  height: 64px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.mobile-menu-btn {
  display: none;
  align-items: center;
  justify-content: center;
  width: 2.25rem;
  height: 2.25rem;
  border-radius: 0.375rem;
  color: var(--color-text-secondary);
  transition: all var(--transition-fast);
  flex-shrink: 0;
}

.mobile-menu-btn:hover {
  color: var(--color-text-primary);
  background: rgba(255, 255, 255, 0.08);
}

.mobile-menu-btn svg {
  width: 1.25rem;
  height: 1.25rem;
}

.header-search {
  flex: 1;
  max-width: 480px;
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 0.875rem;
  width: 1rem;
  height: 1rem;
  color: var(--color-text-muted);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 0.6rem 2.5rem 0.6rem 2.5rem;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.5rem;
  color: var(--color-text-primary);
  font-size: 0.875rem;
  transition: all var(--transition-fast);
}

.search-input::placeholder {
  color: var(--color-text-muted);
}

.search-input:focus {
  outline: none;
  background: rgba(255, 255, 255, 0.09);
  border-color: rgba(229, 9, 20, 0.4);
}

.search-clear {
  position: absolute;
  right: 0.875rem;
  width: 1rem;
  height: 1rem;
  color: var(--color-text-muted);
  transition: color var(--transition-fast);
}

.search-clear:hover { color: var(--color-text-primary); }
.search-clear svg { width: 100%; height: 100%; }

.header-user {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  margin-left: auto;
  cursor: pointer;
  padding: 0.375rem 0.75rem;
  border-radius: 0.5rem;
  transition: background var(--transition-fast);
}

.header-user:hover {
  background: rgba(255, 255, 255, 0.06);
}

.user-avatar {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--color-accent) 0%, #b00 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.user-avatar svg {
  width: 1.25rem;
  height: 1.25rem;
  color: #fff;
}

.user-name {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.chevron-icon {
  width: 1rem;
  height: 1rem;
  color: var(--color-text-muted);
}

/* ==========================================
   PAGE CONTENT
   ========================================== */

.page-content {
  flex: 1;
  overflow-x: hidden;
}

/* ==========================================
   MOBILE BOTTOM NAV
   ========================================== */

.mobile-bottom-nav {
  display: none;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  z-index: 50;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  padding: 0.5rem 1rem;
  padding-bottom: max(0.5rem, env(safe-area-inset-bottom));
}

.mobile-nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.25rem;
  flex: 1;
  padding: 0.5rem;
  color: var(--color-text-muted);
  font-size: 0.7rem;
  font-weight: 500;
  transition: color var(--transition-fast);
  border-radius: 0.375rem;
}

.mobile-nav-item svg {
  width: 1.375rem;
  height: 1.375rem;
}

.mobile-nav-item:hover,
.mobile-nav-item-active {
  color: var(--color-accent);
}

/* ==========================================
   RESPONSIVE
   ========================================== */

@media (max-width: 1024px) {
  .sidebar {
    transform: translateX(-100%);
  }

  .sidebar.sidebar-open {
    transform: translateX(0);
  }

  .main-area {
    margin-left: 0;
  }

  .mobile-menu-btn {
    display: flex;
  }
}

@media (max-width: 768px) {
  .mobile-bottom-nav {
    display: flex;
  }

  .page-content {
    padding-bottom: 5rem;
  }

  .top-header {
    padding: 0.75rem 1rem;
  }

  .header-search {
    max-width: none;
  }

  .user-name {
    display: none;
  }

  .chevron-icon {
    display: none;
  }
}
</style>
