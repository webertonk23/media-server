<template>
  <div class="media-details">
    <!-- Back button -->
    <button class="back-btn" @click="$emit('back')">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <polyline points="15 18 9 12 15 6"/>
      </svg>
    </button>

    <!-- Backdrop -->
    <div class="details-backdrop">
      <img
        v-if="media.backdrop"
        :src="media.backdrop"
        :alt="media.title"
        class="backdrop-img"
      />
      <div v-else class="backdrop-fallback"></div>
      <div class="backdrop-gradient"></div>
    </div>

    <!-- Main content -->
    <div class="details-body">
      <!-- Left: Poster + Info -->
      <div class="details-left">
        <!-- Poster -->
        <div class="poster-wrap">
          <img
            v-if="media.poster"
            :src="media.poster"
            :alt="media.title"
            class="poster-img"
          />
          <div v-else class="poster-fallback">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"/>
            </svg>
          </div>
        </div>

        <!-- Info section -->
        <div class="details-info">
          <h1 class="details-title">{{ media.title }}</h1>

          <!-- Meta row -->
          <div class="details-meta">
            <span v-if="media.year" class="meta-year">{{ media.year }}</span>
            <span v-if="media.year" class="meta-sep">•</span>
            <span class="meta-badge" :class="`type-${media.type}`">
              {{ typeLabel(media.type) }}
            </span>
            <span v-if="media.quality" class="meta-sep">•</span>
            <span v-if="media.quality" class="meta-quality">{{ media.quality }}</span>
          </div>

          <!-- Overview -->
          <p v-if="media.overview" class="details-overview">{{ media.overview }}</p>

          <!-- Progress bar if in progress -->
          <div v-if="progress && !progress.finished" class="progress-section">
            <div class="progress-track">
              <div class="progress-fill" :style="{ width: `${progressPct}%` }"></div>
            </div>
            <span class="progress-text">
              {{ formatTime(progress.position) }} / {{ formatTime(progress.duration) }}
            </span>
          </div>

          <!-- Action buttons -->
          <div class="action-btns">
            <button class="btn btn-primary action-play" @click="handleAction">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M8 5v14l11-7z"/>
              </svg>
              {{ actionLabel }}
            </button>

            <button class="btn btn-ghost action-watchlist">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
              </svg>
              Minha Lista
            </button>
          </div>

          <!-- Metadata table -->
          <div v-if="mediaMeta.length > 0" class="meta-table">
            <div v-for="row in mediaMeta" :key="row.label" class="meta-row">
              <span class="meta-label">{{ row.label }}</span>
              <span class="meta-value">{{ row.value }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right: Mini Player area -->
      <div class="details-right">
        <!-- Thumbnail / Player area -->
        <div class="mini-player-area">
          <img
            v-if="media.backdrop"
            :src="media.backdrop"
            :alt="media.title"
            class="mini-player-thumb"
          />
          <div v-else class="mini-player-fallback">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
              <path d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
              <path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
          <!-- Play overlay -->
          <div class="mini-player-overlay" @click="handleAction">
            <div class="mini-play-btn">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M8 5v14l11-7z"/>
              </svg>
            </div>
          </div>
          <!-- Duration badge -->
          <div v-if="progress" class="mini-progress-bar">
            <div class="mini-progress-fill" :style="{ width: `${progressPct}%` }"></div>
          </div>
        </div>

        <!-- Title in right panel -->
        <div class="right-panel-info">
          <h2 class="right-title">{{ media.title }}</h2>
          <div class="right-meta">
            <span v-if="media.year">{{ media.year }}</span>
            <span class="meta-sep">•</span>
            <span class="meta-badge" :class="`type-${media.type}`">{{ typeLabel(media.type) }}</span>
            <span v-if="media.quality" class="meta-sep">•</span>
            <span v-if="media.quality" class="meta-quality">{{ media.quality }}</span>
          </div>
          <p v-if="media.overview" class="right-overview">{{ media.overview }}</p>
        </div>

        <!-- Quick action buttons for right panel -->
        <div class="right-actions">
          <button class="btn btn-primary right-play-btn" @click="handleAction">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M8 5v14l11-7z"/>
            </svg>
            Assistir
          </button>
          <button class="btn btn-ghost right-watchlist-btn">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"/>
            </svg>
            Assistir do Início
          </button>
          <button class="btn btn-ghost right-list-btn">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <line x1="12" y1="5" x2="12" y2="19"/>
              <line x1="5" y1="12" x2="19" y2="12"/>
            </svg>
            + Minha Lista
          </button>
        </div>

        <!-- Details table (right side) -->
        <div class="right-meta-table" v-if="mediaMeta.length > 0">
          <div v-for="row in mediaMeta" :key="row.label" class="right-meta-row">
            <span class="right-meta-label">{{ row.label }}</span>
            <span class="right-meta-value">{{ row.value }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Seasons and Episodes Section -->
    <div v-if="media.type === 'series' && seasons && seasons.length > 0" class="episodes-section">
      <div class="episodes-header">
        <h2 class="episodes-title">Episódios</h2>
        <div class="season-selector">
          <select 
            :value="selectedSeasonId" 
            @change="(e) => $emit('season-select', (e.target as HTMLSelectElement).value)"
            class="season-select"
          >
            <option v-for="season in seasons" :key="season.id" :value="season.id">
              Temporada {{ season.number }} <span v-if="season.name">- {{ season.name }}</span>
            </option>
          </select>
          <div class="select-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="6 9 12 15 18 9"/>
            </svg>
          </div>
        </div>
      </div>

      <div class="episodes-list" v-if="episodes && episodes.length > 0">
        <div 
          v-for="episode in episodes" 
          :key="episode.id" 
          class="episode-card"
          @click="$emit('play', episode)"
        >
          <div class="episode-thumb-wrap">
            <img v-if="episode.still" :src="episode.still" :alt="episode.title" class="episode-thumb">
            <div v-else class="episode-thumb-fallback">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor">
                <path d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"/>
                <path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <div class="episode-play-overlay">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M8 5v14l11-7z"/>
              </svg>
            </div>
          </div>
          <div class="episode-info">
            <div class="episode-meta">
              <span class="episode-number">{{ episode.episode_number }}</span>
              <h3 class="episode-title">{{ episode.title }}</h3>
              <span v-if="episode.runtime" class="episode-runtime">{{ episode.runtime }}m</span>
            </div>
            <p v-if="episode.overview" class="episode-overview">{{ episode.overview }}</p>
          </div>
        </div>
      </div>
      <div v-else class="episodes-empty">
        <p>Nenhum episódio encontrado para esta temporada.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { MediaItem, Season, Episode } from '@/types/media'
import type { ProgressData } from '@/types/player'

interface Props {
  media: MediaItem
  progress?: ProgressData | null
  seasons?: Season[]
  episodes?: Episode[]
  selectedSeasonId?: string | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'play', media: MediaItem | Episode): void
  (e: 'continue', media: MediaItem): void
  (e: 'back'): void
  (e: 'season-select', seasonId: string): void
}>()

const progressPct = computed(() => {
  if (!props.progress?.duration) return 0
  return Math.min((props.progress.position / props.progress.duration) * 100, 100)
})

const actionLabel = computed(() =>
  props.progress && !props.progress.finished ? 'Continuar' : 'Assistir'
)

const typeLabel = (type: string) => {
  const map: Record<string, string> = { movie: 'Filme', series: 'Série', episode: 'Série' }
  return map[type] || type
}

const handleAction = () => {
  if (props.progress && !props.progress.finished) {
    emit('continue', props.media)
  } else {
    emit('play', props.media)
  }
}

const formatTime = (seconds: number) => {
  const h = Math.floor(seconds / 3600)
  const m = Math.floor((seconds % 3600) / 60)
  const s = Math.floor(seconds % 60)
  if (h > 0) return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
  return `${m}:${String(s).padStart(2, '0')}`
}

const mediaMeta = computed(() => {
  const rows: { label: string; value: string }[] = []
  if (props.media.year) rows.push({ label: 'Ano', value: String(props.media.year) })
  rows.push({ label: 'Tipo', value: typeLabel(props.media.type) })
  return rows
})
</script>

<style scoped>
.media-details {
  position: relative;
  min-height: 100vh;
}

/* Back button */
.back-btn {
  position: absolute;
  top: 1.25rem;
  left: 1.5rem;
  z-index: 20;
  width: 2.5rem;
  height: 2.5rem;
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

.back-btn:hover {
  background: var(--color-accent);
  border-color: var(--color-accent);
}

.back-btn svg { width: 1.25rem; height: 1.25rem; }

/* Backdrop */
.details-backdrop {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 55vh;
  overflow: hidden;
  z-index: 0;
}

.backdrop-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  object-position: center 20%;
}

.backdrop-fallback {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #1a1a2e, #0f3460);
}

.backdrop-gradient {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    to bottom,
    rgba(10, 10, 15, 0.3) 0%,
    rgba(10, 10, 15, 0.8) 60%,
    var(--color-cinema-dark-900) 100%
  );
}

/* Body */
.details-body {
  position: relative;
  z-index: 10;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 3rem;
  padding: 2rem 3rem 4rem;
  padding-top: 30vh;
  max-width: 1400px;
  margin: 0 auto;
}

/* LEFT panel */
.details-left {
  display: flex;
  gap: 1.75rem;
  align-items: flex-start;
}

.poster-wrap {
  flex-shrink: 0;
  width: 180px;
  border-radius: 0.625rem;
  overflow: hidden;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.6);
}

.poster-img {
  width: 100%;
  height: auto;
  display: block;
}

.poster-fallback {
  width: 100%;
  padding-bottom: 150%;
  background: linear-gradient(135deg, var(--color-cinema-dark-700), var(--color-cinema-dark-600));
  position: relative;
}

.poster-fallback svg {
  position: absolute;
  inset: 0;
  margin: auto;
  width: 3rem;
  height: 3rem;
  color: rgba(255,255,255,0.2);
}

.details-info {
  flex: 1;
  min-width: 0;
}

.details-title {
  font-family: var(--font-display);
  font-size: clamp(1.5rem, 3vw, 2.25rem);
  font-weight: 700;
  color: #fff;
  line-height: 1.2;
  margin-bottom: 0.75rem;
}

.details-meta {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-bottom: 1rem;
}

.meta-year {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
}

.meta-sep {
  color: var(--color-text-muted);
}

.meta-badge {
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.125rem 0.5rem;
  border-radius: 0.25rem;
}

.meta-quality {
  font-size: 0.75rem;
  font-weight: 700;
  color: var(--color-cinema-accent-secondary);
  background: rgba(245, 197, 24, 0.1);
  padding: 0.125rem 0.4rem;
  border-radius: 0.25rem;
  border: 1px solid rgba(245, 197, 24, 0.2);
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

.details-overview {
  font-size: 0.875rem;
  line-height: 1.65;
  color: var(--color-text-secondary);
  margin-bottom: 1.25rem;
  display: -webkit-box;
  -webkit-line-clamp: 5;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Progress */
.progress-section {
  margin-bottom: 1.25rem;
}

.progress-track {
  height: 3px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
  margin-bottom: 0.375rem;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--color-accent);
}

.progress-text {
  font-size: 0.75rem;
  color: var(--color-text-muted);
}

/* Actions */
.action-btns {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
  margin-bottom: 1.5rem;
}

.action-play svg { width: 1.125rem; height: 1.125rem; }
.action-watchlist svg { width: 1rem; height: 1rem; }

/* Meta table */
.meta-table {
  display: flex;
  flex-direction: column;
  gap: 0.375rem;
}

.meta-row {
  display: flex;
  gap: 0.75rem;
  font-size: 0.8125rem;
}

.meta-label {
  color: var(--color-text-muted);
  min-width: 80px;
  flex-shrink: 0;
}

.meta-value {
  color: var(--color-text-secondary);
}

/* RIGHT panel */
.details-right {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.mini-player-area {
  position: relative;
  border-radius: 0.75rem;
  overflow: hidden;
  aspect-ratio: 16/9;
  background: #000;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
  cursor: pointer;
}

.mini-player-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.mini-player-fallback {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #1a1a2e, #0f3460);
  display: flex;
  align-items: center;
  justify-content: center;
}

.mini-player-fallback svg {
  width: 4rem;
  height: 4rem;
  color: rgba(255, 255, 255, 0.3);
}

.mini-player-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.3);
  transition: background var(--transition-fast);
}

.mini-player-area:hover .mini-player-overlay {
  background: rgba(0, 0, 0, 0.15);
}

.mini-play-btn {
  width: 4rem;
  height: 4rem;
  border-radius: 50%;
  background: rgba(229, 9, 20, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 24px rgba(229, 9, 20, 0.6);
  transition: transform var(--transition-fast);
}

.mini-player-area:hover .mini-play-btn {
  transform: scale(1.1);
}

.mini-play-btn svg {
  width: 1.75rem;
  height: 1.75rem;
  color: #fff;
  margin-left: 3px;
}

.mini-progress-bar {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: rgba(255,255,255,0.15);
}

.mini-progress-fill {
  height: 100%;
  background: var(--color-accent);
}

/* Right panel info */
.right-panel-info {}

.right-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #fff;
  margin-bottom: 0.5rem;
}

.right-meta {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
}

.right-overview {
  font-size: 0.8125rem;
  line-height: 1.6;
  color: var(--color-text-secondary);
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Right actions */
.right-actions {
  display: flex;
  gap: 0.625rem;
  flex-wrap: wrap;
}

.right-play-btn {
  padding: 0.6rem 1.25rem;
  font-size: 0.875rem;
}

.right-play-btn svg { width: 1rem; height: 1rem; }

.right-watchlist-btn,
.right-list-btn {
  padding: 0.6rem 1rem;
  font-size: 0.8125rem;
}

.right-watchlist-btn svg,
.right-list-btn svg { width: 0.875rem; height: 0.875rem; }

/* Right meta table */
.right-meta-table {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  padding: 0.75rem 0;
  border-top: 1px solid rgba(255,255,255,0.06);
}

.right-meta-row {
  display: flex;
  gap: 1rem;
  font-size: 0.8rem;
}

.right-meta-label {
  color: var(--color-text-muted);
  min-width: 70px;
}

.right-meta-value {
  color: var(--color-text-secondary);
}

/* Features grid */
.features-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0.875rem;
  padding-top: 0.75rem;
  border-top: 1px solid rgba(255,255,255,0.06);
}

.feature-item {
  display: flex;
  align-items: flex-start;
  gap: 0.625rem;
}

.feature-item svg {
  width: 1.375rem;
  height: 1.375rem;
  flex-shrink: 0;
  color: var(--color-accent);
  margin-top: 1px;
}

.feature-item div {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.feature-item strong {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.feature-item span {
  font-size: 0.7rem;
  color: var(--color-text-muted);
  line-height: 1.4;
}

/* Responsive */
@media (max-width: 1024px) {
  .details-body {
    grid-template-columns: 1fr;
    padding: 2rem 1.5rem 4rem;
    padding-top: 25vh;
    gap: 2rem;
  }

  .details-right {
    order: -1;
  }

  .mini-player-area {
    max-width: 480px;
  }
}

@media (max-width: 640px) {
  .details-left {
    flex-direction: column;
  }

  .poster-wrap {
    width: 140px;
  }

  .details-body {
    padding: 1rem 1rem 4rem;
    padding-top: 20vh;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }
}

/* Episodes Section */
.episodes-section {
  position: relative;
  z-index: 10;
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 3rem 4rem;
}

.episodes-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding-bottom: 1rem;
}

.episodes-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #fff;
}

.season-selector {
  position: relative;
  min-width: 200px;
}

.season-select {
  appearance: none;
  width: 100%;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #fff;
  padding: 0.625rem 2.5rem 0.625rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.9rem;
  font-weight: 500;
  cursor: pointer;
  outline: none;
  transition: border-color var(--transition-fast);
}

.season-select:hover, .season-select:focus {
  border-color: rgba(255, 255, 255, 0.3);
}

.season-select option {
  background: var(--color-cinema-dark-900);
  color: #fff;
}

.select-icon {
  position: absolute;
  right: 1rem;
  top: 50%;
  transform: translateY(-50%);
  pointer-events: none;
  color: rgba(255, 255, 255, 0.5);
}

.select-icon svg {
  width: 1rem;
  height: 1rem;
}

.episodes-list {
  display: flex;
  gap: 1.25rem;
  overflow-x: auto;
  padding-bottom: 1.5rem;
  scroll-snap-type: x mandatory;
  scrollbar-width: thin;
  scrollbar-color: var(--color-cinema-dark-500) transparent;
}

.episodes-list::-webkit-scrollbar {
  height: 6px;
}

.episodes-list::-webkit-scrollbar-thumb {
  background: var(--color-cinema-dark-500);
  border-radius: 10px;
}

.episode-card {
  flex: 0 0 280px;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  padding: 0.75rem;
  border-radius: 0.75rem;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.05);
  cursor: pointer;
  transition: all var(--transition-fast);
  scroll-snap-align: start;
}

.episode-card:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
  transform: translateX(4px);
}

.episode-thumb-wrap {
  position: relative;
  width: 100%;
  flex-shrink: 0;
  aspect-ratio: 16/9;
  border-radius: 0.5rem;
  overflow: hidden;
  background: #000;
}

.episode-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-fast);
}

.episode-card:hover .episode-thumb {
  transform: scale(1.05);
}

.episode-thumb-fallback {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(255,255,255,0.05), rgba(255,255,255,0.1));
}

.episode-thumb-fallback svg {
  width: 2rem;
  height: 2rem;
  color: rgba(255, 255, 255, 0.2);
}

.episode-play-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.episode-card:hover .episode-play-overlay {
  opacity: 1;
}

.episode-play-overlay svg {
  width: 2.5rem;
  height: 2.5rem;
  color: #fff;
  filter: drop-shadow(0 2px 4px rgba(0,0,0,0.5));
}

.episode-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.episode-meta {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 0.5rem;
}

.episode-number {
  font-size: 1.125rem;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.5);
  min-width: 1.5rem;
}

.episode-title {
  font-size: 1rem;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.episode-runtime {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.4);
  margin-left: auto;
}

.episode-overview {
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin: 0;
}

.episodes-empty {
  padding: 3rem;
  text-align: center;
  color: var(--color-text-muted);
  background: rgba(255, 255, 255, 0.02);
  border-radius: 0.75rem;
}

@media (max-width: 1024px) {
  .episodes-section {
    padding: 0 1.5rem 4rem;
  }
}

@media (max-width: 640px) {
  .episodes-section {
    padding: 0 1rem 4rem;
  }
  
  .episodes-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
  
  .season-selector {
    width: 100%;
  }

  .episode-card {
    flex-direction: column;
    gap: 1rem;
  }

  .episode-thumb-wrap {
    width: 100%;
  }
}
</style>
