/**
 * Application-wide constants
 * 
 * **Validates: Requirements 6.1, 6.5, 8.1**
 */

/**
 * API Configuration
 */
export const API_BASE_URL = 'http://localhost:9000';

/**
 * API Endpoints
 */
export const API_ENDPOINTS = {
  MEDIA: '/media',
  MEDIA_BY_ID: (id: string) => `/media/${id}`,
  PROGRESS: (id: string) => `/progress/${id}`,
  STREAM: (id: string) => `/stream/${id}`,
  SEARCH: '/media',
} as const;

/**
 * Progress Tracking Configuration
 */
export const PROGRESS_SAVE_INTERVAL = 10000; // 10 seconds in milliseconds
export const COMPLETION_THRESHOLD = 0.95; // 95% completion marks video as finished

/**
 * Search Configuration
 */
export const SEARCH_DEBOUNCE_DELAY = 300; // 300ms debounce for search input

/**
 * Pagination Configuration
 */
export const DEFAULT_PAGE_SIZE = 20;
export const DEFAULT_PAGE = 1;

/**
 * Player Configuration
 */
export const SEEK_STEP = 10; // Seconds to skip forward/backward with arrow keys
export const VOLUME_STEP = 0.1; // Volume increment/decrement step (0-1 range)

/**
 * UI Configuration
 */
export const CONTROLS_HIDE_DELAY = 3000; // 3 seconds before hiding player controls
export const IMAGE_LAZY_LOAD_THRESHOLD = 100; // Pixels before viewport to start loading images

/**
 * Breakpoints for responsive design
 */
export const BREAKPOINTS = {
  MOBILE: 768,
  TABLET: 1024,
  DESKTOP: 1920,
} as const;
