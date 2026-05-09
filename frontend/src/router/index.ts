import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'home',
    component: () => import('../pages/HomePage.vue'),
    meta: { title: 'Início' }
  },
  {
    path: '/movies',
    name: 'movies',
    component: () => import('../pages/CatalogPage.vue'),
    meta: { title: 'Filmes' },
    props: { type: 'movie' }
  },
  {
    path: '/series',
    name: 'series',
    component: () => import('../pages/CatalogPage.vue'),
    meta: { title: 'Séries' },
    props: { type: 'series' }
  },
  {
    path: '/continue',
    name: 'continue',
    component: () => import('../pages/ContinuePage.vue'),
    meta: { title: 'Continuar Assistindo' }
  },
  {
    path: '/recent',
    name: 'recent',
    component: () => import('../pages/RecentPage.vue'),
    meta: { title: 'Recém Adicionados' }
  },
  {
    path: '/watchlist',
    name: 'watchlist',
    component: () => import('../pages/WatchlistPage.vue'),
    meta: { title: 'Minha Lista' }
  },
  {
    path: '/search',
    name: 'search',
    component: () => import('../pages/SearchPage.vue'),
    meta: { title: 'Buscar' }
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('../pages/SettingsPage.vue'),
    meta: { title: 'Configurações' }
  },
  {
    path: '/transcode-monitor',
    name: 'transcode-monitor',
    component: () => import('../pages/TranscodeMonitorPage.vue'),
    meta: { title: 'Monitor de Transcode' }
  },
  {
    path: '/logs',
    name: 'logs',
    component: () => import('../pages/LogsPage.vue'),
    meta: { title: 'Logs do Sistema' }
  },
  {
    path: '/media/:id',
    name: 'media-detail',
    component: () => import('../pages/MediaDetailPage.vue'),
    meta: { title: 'Detalhes' },
    props: true
  },
  {
    path: '/player/:id',
    name: 'player',
    component: () => import('../pages/PlayerPage.vue'),
    meta: { title: 'Player', layout: 'player' },
    props: true
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('../pages/NotFoundPage.vue'),
    meta: { title: 'Não encontrado' }
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  scrollBehavior(_to, _from, savedPosition) {
    if (savedPosition) return savedPosition
    return { top: 0 }
  }
})

router.beforeEach((to, _from, next) => {
  if (to.meta.title) {
    document.title = `${to.meta.title} — MediaServer`
  }
  next()
})

export default router
