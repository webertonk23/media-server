/**
 * Core media types matching backend Go DTOs
 * 
 * **Validates: Requirements 20.1, 16.3**
 */

/**
 * Media type constants matching backend MediaType
 */
export const MediaType = {
  Movie: 'movie',
  Series: 'series',
  Episode: 'episode'
} as const;

/**
 * Media type union type
 */
export type MediaType = typeof MediaType[keyof typeof MediaType];

/**
 * MediaFile interface matching backend MediaFile model
 */
export interface MediaFile {
  id: number;
  ulid: string;
  media_item_id: number;
  path: string;
  size: number;
  fingerprint: string;
  quality: string;
  status: string;
  original_path?: string;
  error_message?: string;
  retry_count: number;
  created_at: string;
  updated_at: string;
}

/**
 * MediaItem interface matching backend MediaItemResponse DTO
 * Represents a media item (movie, series, or episode) with associated metadata
 */
export interface MediaItem {
  /** ULID public identifier */
  id: string;
  
  /** Media type: 'movie', 'series', or 'episode' */
  type: MediaType;
  
  /** Media title */
  title: string;
  
  /** Release year (optional) */
  year?: number;
  
  /** Media overview/description (optional) */
  overview?: string;
  
  /** Poster image URL (optional) */
  poster?: string;
  
  /** Backdrop image URL (optional) */
  backdrop?: string;
  
  /** Streaming URL (optional) */
  stream_url?: string;
}

/**
 * Generic paginated response matching backend PaginatedResponse DTO
 * Used for paginated API responses
 */
export interface PaginatedResponse<T> {
  /** Current page number */
  page: number;
  
  /** Items per page limit */
  limit: number;
  
  /** Total number of items */
  total: number;
  
  /** Array of items for current page */
  items: T[];
}

/**
 * Season interface matching backend SeasonResponse DTO
 */
export interface Season {
  id: string;
  number: number;
  name?: string;
  overview?: string;
  poster?: string;
  episode_count?: number;
}

/**
 * Episode interface matching backend EpisodeResponse DTO
 */
export interface Episode {
  id: string;
  type: string;
  title: string;
  overview?: string;
  season_number: number;
  episode_number: number;
  still?: string;
  runtime?: number;
  stream_url?: string;
}
