/**
 * Media Service
 * 
 * Provides methods for fetching media items from the backend API.
 * Handles media listing, searching, filtering, and individual item retrieval.
 * 
 * **Validates: Requirements 2.1, 4.1, 8.2, 16.3**
 */
import apiClient from './api';
import type { MediaItem, PaginatedResponse } from '@/types/media';
import type { SearchParams } from '@/types/api';
/**
 * Get paginated media items with optional filtering
 * 
 * Fetches media items from GET /media endpoint with support for:
 * - Pagination (page, limit)
 * - Search by title (search)
 * - Filter by type (type: 'movie', 'series', 'episode')
 * 
 * **Validates: Requirements 2.1, 16.3**
 * 
 * @param params - Search and pagination parameters
 * @returns Promise resolving to paginated media items
 * @throws ApiError if request fails
 * 
 * @example
 * ```typescript
 * 
 * const result = await getMedia({ page: 1, limit: 20 });
 * 
 * 
 * const movies = await getMedia({ search: 'Matrix', type: 'movie' });
 * ```
 */
export async function getMedia(params: SearchParams = {}): Promise<PaginatedResponse<MediaItem>> {
  const response = await apiClient.get<PaginatedResponse<MediaItem>>('/media', {
    params: {
      page: params.page || 1,
      limit: params.limit || 20,
      search: params.search || '',
      type: params.type || '',
    },
  });
  return response.data;
}
/**
 * Get a specific media item by ID
 * 
 * Fetches detailed information for a single media item from GET /media/:id endpoint.
 * 
 * **Validates: Requirements 4.1, 16.3**
 * 
 * @param id - ULID of the media item
 * @returns Promise resolving to the media item
 * @throws ApiError with status 404 if media not found
 * 
 * @example
 * ```typescript
 * const media = await getMediaById('01HQXYZ123ABC456DEF789GHI');
 * console.log(media.title, media.year);
 * ```
 */
export async function getMediaById(id: string): Promise<MediaItem> {
  const response = await apiClient.get<MediaItem>(`/media/${id}`);
  return response.data;
}
/**
 * Get paginated movies with optional filtering
 * 
 * Fetches only movie-type media items from GET /movies endpoint.
 * Convenience method that filters by type='movie'.
 * 
 * **Validates: Requirements 2.1, 16.3**
 * 
 * @param params - Search and pagination parameters
 * @returns Promise resolving to paginated movie items
 * @throws ApiError if request fails
 * 
 * @example
 * ```typescript
 * 
 * const movies = await getMovies({ page: 1, limit: 20 });
 * 
 * 
 * const results = await getMovies({ search: 'Inception' });
 * ```
 */
export async function getMovies(params: SearchParams = {}): Promise<PaginatedResponse<MediaItem>> {
  const response = await apiClient.get<PaginatedResponse<MediaItem>>('/movies', {
    params: {
      page: params.page || 1,
      limit: params.limit || 20,
      search: params.search || '',
    },
  });
  return response.data;
}
/**
 * Search media items by query string
 * 
 * Searches across all media types using the search parameter.
 * This is a convenience method that wraps getMedia() with search parameter.
 * 
 * **Validates: Requirements 8.2, 16.3**
 * 
 * @param query - Search query string to match against titles
 * @param page - Page number for pagination (default: 1)
 * @returns Promise resolving to paginated search results
 * @throws ApiError if request fails
 * 
 * @example
 * ```typescript
 * 
 * const results = await searchMedia('Star Wars', 1);
 * 
 * 
 * const nextPage = await searchMedia('Star Wars', 2);
 * ```
 */
export async function searchMedia(query: string, page: number = 1): Promise<PaginatedResponse<MediaItem>> {
  return getMedia({
    search: query,
    page,
    limit: 20,
  });
}
/**
 * Get seasons for a series
 * 
 * @param seriesId - ULID of the series
 * @returns Promise resolving to an array of seasons
 */
export async function getSeriesSeasons(seriesId: string) {
  const response = await apiClient.get(`/series/${seriesId}/seasons`);
  return response.data;
}
/**
 * Get episodes for a season
 * 
 * @param seasonId - ULID of the season
 * @returns Promise resolving to an array of episodes
 */
export async function getSeasonEpisodes(seasonId: string) {
  const response = await apiClient.get(`/seasons/${seasonId}/episodes`);
  return response.data;
}
