export const MediaType = {
  Movie: 'movie',
  Series: 'series',
  Episode: 'episode'
} as const;

export type MediaType = typeof MediaType[keyof typeof MediaType];

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

export interface MediaItem {
  id: string;
  type: MediaType;
  title: string;
  year?: number;
  overview?: string;
  poster?: string;
  backdrop?: string;
  stream_url?: string;
  quality?: string;
}

export interface PaginatedResponse<T> {
  page: number;
  limit: number;
  total: number;
  items: T[];
}

export interface Season {
  id: string;
  number: number;
  name?: string;
  overview?: string;
  poster?: string;
  episode_count?: number;
}

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
  quality?: string;
}
