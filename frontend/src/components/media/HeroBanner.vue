<template>
  <div class="hero-banner">
    <!-- Slides -->
    <div class="hero-slides" :style="{ transform: `translateX(-${currentSlide * 100}%)` }">
      <div
        v-for="item in featuredItems"
        :key="item.id"
        class="hero-slide"
      >
        <!-- Backdrop -->
        <div class="hero-backdrop">
          <img
            v-if="item.backdrop"
            :src="item.backdrop"
            :alt="`${item.title} backdrop`"
            class="backdrop-image"
          />
          <div v-else class="backdrop-fallback">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" class="fallback-icon">
              <path d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z"/>
            </svg>
          </div>
          <!-- Gradient overlays -->
          <div class="gradient-left"></div>
          <div class="gradient-bottom"></div>
        </div>
        <!-- Content -->
        <div class="hero-content">
          <div class="hero-info">
            <!-- Title -->
            <h1 class="hero-title">{{ item.title }}</h1>
            <!-- Metadata row -->
            <div class="hero-meta">
              <span v-if="item.year" class="meta-item">{{ item.year }}</span>
              <span v-if="item.duration" class="meta-dot">•</span>
              <span v-if="item.duration" class="meta-item">{{ item.duration }}</span>
              <span class="meta-dot">•</span>
              <span class="meta-item type-badge" :class="`type-${item.type}`">
                {{ typeLabel(item.type) }}
              </span>
              <span v-if="item.genres?.length" class="meta-dot">•</span>
              <span v-if="item.genres?.length" class="meta-item meta-genres">
                {{ item.genres.join(', ') }}
              </span>
            </div>
            <!-- Overview -->
            <p v-if="item.overview" class="hero-overview">{{ item.overview }}</p>
            <!-- Actions -->
            <div class="hero-actions">
              <button class="btn btn-primary hero-btn-play" @click="$emit('play', item)">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M8 5v14l11-7z"/>
                </svg>
                Assistir
              </button>
              <button class="btn btn-ghost hero-btn-info" @click="$emit('more-info', item)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <line x1="12" y1="8" x2="12" y2="12"/>
                  <line x1="12" y1="16" x2="12.01" y2="16"/>
                </svg>
                Mais informações
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- Navigation Arrows -->
    <button v-if="featuredItems.length > 1" class="hero-arrow hero-arrow-left" @click="prevSlide">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="15 18 9 12 15 6"/>
      </svg>
    </button>
    <button v-if="featuredItems.length > 1" class="hero-arrow hero-arrow-right" @click="nextSlide">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="9 18 15 12 9 6"/>
      </svg>
    </button>
    <!-- Dots navigation -->
    <div v-if="featuredItems.length > 1" class="hero-dots">
      <button
        v-for="(_, i) in featuredItems"
        :key="i"
        class="hero-dot"
        :class="{ 'hero-dot-active': i === currentSlide }"
        @click="goToSlide(i)"
      ></button>
    </div>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import type { MediaItem } from '@/types/media'
interface Props {
  media: MediaItem
  allMedia?: MediaItem[]
}
const props = defineProps<Props>()
defineEmits<{
  (e: 'play', media: MediaItem): void
  (e: 'more-info', media: MediaItem): void
}>()
const currentSlide = ref(0)
let autoplayTimer: ReturnType<typeof setInterval> | null = null
const featuredItems = computed<(MediaItem & { duration?: string; genres?: string[] })[]>(() => {
  const items = props.allMedia?.slice(0, 5) || [props.media]
  return items.map(item => ({
    ...item,
    duration: undefined,
    genres: undefined,
  }))
})
const typeLabel = (type: string) => {
  const map: Record<string, string> = {
    movie: 'Filme',
    series: 'Série',
    episode: 'Episódio',
  }
  return map[type] || type
}
const nextSlide = () => {
  currentSlide.value = (currentSlide.value + 1) % featuredItems.value.length
}
const prevSlide = () => {
  currentSlide.value = (currentSlide.value - 1 + featuredItems.value.length) % featuredItems.value.length
}
const goToSlide = (i: number) => {
  currentSlide.value = i
  resetAutoplay()
}
const startAutoplay = () => {
  if (featuredItems.value.length > 1) {
    autoplayTimer = setInterval(nextSlide, 6000)
  }
}
const resetAutoplay = () => {
  if (autoplayTimer) clearInterval(autoplayTimer)
  startAutoplay()
}
onMounted(startAutoplay)
onUnmounted(() => {
  if (autoplayTimer) clearInterval(autoplayTimer)
})
</script>
<style scoped>
.hero-banner {
  position: relative;
  width: 100%;
  height: 72vh;
  min-height: 480px;
  max-height: 760px;
  overflow: hidden;
}
/* Slides Container */
.hero-slides {
  display: flex;
  width: 100%;
  height: 100%;
  transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}
.hero-slide {
  min-width: 100%;
  height: 100%;
  position: relative;
}
/* Backdrop */
.hero-backdrop {
  position: absolute;
  inset: 0;
}
.backdrop-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center 20%;
}
.backdrop-fallback {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}
.fallback-icon {
  width: 5rem;
  height: 5rem;
  color: rgba(255, 255, 255, 0.15);
}
.gradient-left {
  position: absolute;
  top: 0;
  left: 0;
  width: 60%;
  height: 100%;
  background: linear-gradient(
    to right,
    rgba(10, 10, 15, 0.95) 0%,
    rgba(10, 10, 15, 0.7) 50%,
    transparent 100%
  );
}
.gradient-bottom {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 50%;
  background: linear-gradient(
    to top,
    var(--color-cinema-dark-900) 0%,
    transparent 100%
  );
}
/* Hero Content */
.hero-content {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 3rem 3rem 4rem;
  z-index: 5;
}
.hero-info {
  max-width: 580px;
}
.hero-title {
  font-family: var(--font-display);
  font-size: clamp(2rem, 5vw, 3.5rem);
  font-weight: 700;
  color: #fff;
  line-height: 1.1;
  margin-bottom: 0.875rem;
  text-shadow: 0 2px 20px rgba(0, 0, 0, 0.5);
}
.hero-meta {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.375rem;
  margin-bottom: 1rem;
}
.meta-item {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  font-weight: 400;
}
.meta-dot {
  color: var(--color-text-muted);
  font-size: 0.75rem;
}
.meta-genres {
  color: var(--color-text-secondary);
}
.type-badge {
  padding: 0.125rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-weight: 600;
}
.type-movie {
  background: rgba(229, 9, 20, 0.2);
  color: #ff6b6b;
  border: 1px solid rgba(229, 9, 20, 0.3);
}
.type-series {
  background: rgba(99, 102, 241, 0.2);
  color: #a5b4fc;
  border: 1px solid rgba(99, 102, 241, 0.3);
}
.hero-overview {
  font-size: 0.9375rem;
  line-height: 1.65;
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 1.75rem;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  max-width: 520px;
}
.hero-actions {
  display: flex;
  gap: 0.875rem;
  flex-wrap: wrap;
}
.hero-btn-play {
  min-width: 140px;
}
.hero-btn-info {
  min-width: 180px;
}
/* Navigation arrows */
.hero-arrow {
  position: absolute;
  top: 50%;
  transform: translateY(-60%);
  z-index: 10;
  width: 2.75rem;
  height: 2.75rem;
  border-radius: 50%;
  background: rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.15);
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
  backdrop-filter: blur(8px);
}
.hero-arrow:hover {
  background: rgba(229, 9, 20, 0.7);
  border-color: var(--color-accent);
  transform: translateY(-60%) scale(1.1);
}
.hero-arrow svg { width: 1.25rem; height: 1.25rem; }
.hero-arrow-left { left: 1.5rem; }
.hero-arrow-right { right: 1.5rem; }
/* Dots */
.hero-dots {
  position: absolute;
  bottom: 1.5rem;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 0.5rem;
  z-index: 10;
}
.hero-dot {
  width: 2rem;
  height: 3px;
  border-radius: 2px;
  background: rgba(255, 255, 255, 0.3);
  transition: all var(--transition-fast);
}
.hero-dot-active {
  background: var(--color-accent);
  width: 2.5rem;
}
/* Responsive */
@media (max-width: 768px) {
  .hero-banner {
    height: 60vh;
    min-height: 380px;
  }
  .hero-content {
    padding: 1.5rem 1.25rem 3rem;
  }
  .hero-title {
    font-size: 1.75rem;
  }
  .hero-overview {
    font-size: 0.875rem;
    -webkit-line-clamp: 2;
  }
  .hero-btn-play,
  .hero-btn-info {
    padding: 0.625rem 1.25rem;
    font-size: 0.875rem;
  }
  .gradient-left {
    width: 100%;
    background: linear-gradient(
      to top,
      rgba(10, 10, 15, 0.98) 0%,
      rgba(10, 10, 15, 0.5) 60%,
      transparent 100%
    );
  }
}
</style>
