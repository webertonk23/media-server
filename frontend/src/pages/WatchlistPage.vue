<template>
  <DefaultLayout>
    <div class="catalog-page">
      <div class="catalog-header">
        <h1 class="catalog-title">Minha Lista</h1>
        <p class="catalog-subtitle">Seus favoritos salvos</p>
      </div>
      <div class="catalog-empty">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
        </svg>
        <p>Sua lista está vazia</p>
        <span>Adicione filmes e séries à sua lista clicando em "+ Minha Lista"</span>
        <router-link to="/" class="btn-primary">Explorar</router-link>
      </div>
    </div>
  </DefaultLayout>
</template>

<script setup lang="ts">
import DefaultLayout from '@/layouts/DefaultLayout.vue'
</script>

<style scoped>
.catalog-page { padding: 2rem; min-height: 100vh; }
.catalog-header { margin-bottom: 2rem; }
.catalog-title { font-family: var(--font-display); font-size: 2rem; font-weight: 700; margin-bottom: 0.375rem; }
.catalog-subtitle { font-size: 0.9rem; color: var(--color-text-muted); }
.catalog-empty { display: flex; flex-direction: column; align-items: center; gap: 0.875rem; padding: 6rem 2rem; text-align: center; }
.catalog-empty svg { width: 5rem; height: 5rem; color: rgba(255,255,255,0.12); }
.catalog-empty p { font-size: 1.375rem; font-weight: 600; color: var(--color-text-secondary); }
.catalog-empty span { font-size: 0.9rem; color: var(--color-text-muted); max-width: 320px; margin-bottom: 0.5rem; }
</style>
