<template>
  <header class="app-header">
    <div class="header-content">
      <!-- Logo/Title -->
      <div class="logo-section">
        <router-link to="/" class="logo-link">
          <svg class="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z" />
          </svg>
          <span class="logo-text">MediaStream</span>
        </router-link>
      </div>
      <!-- Desktop Navigation -->
      <nav class="desktop-nav">
        <router-link to="/" class="nav-link" active-class="nav-link-active">
          Home
        </router-link>
        <router-link to="/movies" class="nav-link" active-class="nav-link-active">
          Movies
        </router-link>
        <router-link to="/series" class="nav-link" active-class="nav-link-active">
          Series
        </router-link>
      </nav>
      <!-- Search Bar -->
      <div class="search-section">
        <div class="search-container">
          <svg class="search-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search media..."
            class="search-input"
            @keyup.enter="handleSearch"
          />
          <button
            v-if="searchQuery"
            @click="clearSearch"
            class="clear-button"
            aria-label="Clear search"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
      <!-- Mobile Menu Toggle -->
      <button
        class="mobile-menu-toggle"
        @click="toggleMobileMenu"
        aria-label="Toggle menu"
      >
        <svg v-if="!isMobileMenuOpen" viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
    <!-- Mobile Navigation Menu -->
    <transition name="mobile-menu">
      <nav v-if="isMobileMenuOpen" class="mobile-nav">
        <router-link
          to="/"
          class="mobile-nav-link"
          active-class="mobile-nav-link-active"
          @click="closeMobileMenu"
        >
          Home
        </router-link>
        <router-link
          to="/movies"
          class="mobile-nav-link"
          active-class="mobile-nav-link-active"
          @click="closeMobileMenu"
        >
          Movies
        </router-link>
        <router-link
          to="/series"
          class="mobile-nav-link"
          active-class="mobile-nav-link-active"
          @click="closeMobileMenu"
        >
          Series
        </router-link>
      </nav>
    </transition>
  </header>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
const router = useRouter()
const searchQuery = ref('')
const isMobileMenuOpen = ref(false)
const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({
      path: '/search',
      query: { q: searchQuery.value.trim() }
    })
    closeMobileMenu()
  }
}
const clearSearch = () => {
  searchQuery.value = ''
}
const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}
const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
}
</script>
<style scoped>
.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 50;
  background: rgba(10, 10, 10, 0.85);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s ease;
}
.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 2560px;
  margin: 0 auto;
  padding: 1rem 2rem;
  gap: 2rem;
}
/* Logo Section */
.logo-section {
  flex-shrink: 0;
}
.logo-link {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  text-decoration: none;
  color: #ffffff;
  transition: opacity 0.2s ease;
}
.logo-link:hover {
  opacity: 0.8;
}
.logo-icon {
  width: 2rem;
  height: 2rem;
  color: #ef4444;
}
.logo-text {
  font-size: 1.5rem;
  font-weight: 700;
  letter-spacing: -0.025em;
}
/* Desktop Navigation */
.desktop-nav {
  display: none;
  gap: 2rem;
}
.nav-link {
  color: rgba(255, 255, 255, 0.7);
  text-decoration: none;
  font-size: 1rem;
  font-weight: 500;
  transition: color 0.2s ease;
  position: relative;
}
.nav-link:hover {
  color: #ffffff;
}
.nav-link-active {
  color: #ffffff;
}
.nav-link-active::after {
  content: '';
  position: absolute;
  bottom: -0.5rem;
  left: 0;
  right: 0;
  height: 2px;
  background: #ef4444;
}
/* Search Section */
.search-section {
  flex: 1;
  max-width: 500px;
  display: none;
}
.search-container {
  position: relative;
  display: flex;
  align-items: center;
}
.search-icon {
  position: absolute;
  left: 1rem;
  width: 1.25rem;
  height: 1.25rem;
  color: rgba(255, 255, 255, 0.5);
  pointer-events: none;
}
.search-input {
  width: 100%;
  padding: 0.75rem 3rem 0.75rem 3rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 0.5rem;
  color: #ffffff;
  font-size: 0.875rem;
  transition: all 0.2s ease;
}
.search-input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}
.search-input:focus {
  outline: none;
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(239, 68, 68, 0.5);
}
.clear-button {
  position: absolute;
  right: 1rem;
  width: 1.25rem;
  height: 1.25rem;
  padding: 0;
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.5);
  cursor: pointer;
  transition: color 0.2s ease;
}
.clear-button:hover {
  color: #ffffff;
}
.clear-button svg {
  width: 100%;
  height: 100%;
}
/* Mobile Menu Toggle */
.mobile-menu-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  padding: 0;
  background: none;
  border: none;
  color: #ffffff;
  cursor: pointer;
  transition: opacity 0.2s ease;
}
.mobile-menu-toggle:hover {
  opacity: 0.7;
}
.mobile-menu-toggle svg {
  width: 1.5rem;
  height: 1.5rem;
}
/* Mobile Navigation */
.mobile-nav {
  display: flex;
  flex-direction: column;
  padding: 1rem 2rem;
  background: rgba(15, 15, 15, 0.95);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}
.mobile-nav-link {
  padding: 1rem 0;
  color: rgba(255, 255, 255, 0.7);
  text-decoration: none;
  font-size: 1.125rem;
  font-weight: 500;
  transition: color 0.2s ease;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}
.mobile-nav-link:hover {
  color: #ffffff;
}
.mobile-nav-link-active {
  color: #ef4444;
}
/* Mobile Menu Transition */
.mobile-menu-enter-active,
.mobile-menu-leave-active {
  transition: all 0.3s ease;
}
.mobile-menu-enter-from,
.mobile-menu-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
/* Tablet and up */
@media (min-width: 768px) {
  .desktop-nav {
    display: flex;
  }
  .search-section {
    display: block;
  }
  .mobile-menu-toggle {
    display: none;
  }
  .mobile-nav {
    display: none;
  }
}
/* Desktop */
@media (min-width: 1024px) {
  .header-content {
    padding: 1.25rem 3rem;
  }
  .search-section {
    max-width: 600px;
  }
}
/* Ultrawide */
@media (min-width: 1920px) {
  .header-content {
    padding: 1.5rem 4rem;
  }
}
</style>
