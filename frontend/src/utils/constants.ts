export const API_BASE_URL = '/api';export const API_ENDPOINTS = {
  MEDIA: '/media',
  MEDIA_BY_ID: (id: string) => `/media/${id}`,
  PROGRESS: (id: string) => `/progress/${id}`,
  STREAM: (id: string) => `/stream/${id}`,
  SEARCH: '/media',
} as const;
export const PROGRESS_SAVE_INTERVAL = 10000;
export const COMPLETION_THRESHOLD = 0.95;
export const SEARCH_DEBOUNCE_DELAY = 300;
export const DEFAULT_PAGE_SIZE = 20;
export const DEFAULT_PAGE = 1;
export const SEEK_STEP = 10;
export const VOLUME_STEP = 0.1;
export const CONTROLS_HIDE_DELAY = 3000;
export const IMAGE_LAZY_LOAD_THRESHOLD = 100;
export const BREAKPOINTS = {
  MOBILE: 768,
  TABLET: 1024,
  DESKTOP: 1920,
} as const;
