import type { MediaItem } from './media';

export interface ApiError {
  message: string;
  status: number;
  code?: string;
}
export interface UpdateProgressRequest {
  position: number;
  duration: number;
  finished: boolean;
}
export interface ContinueWatchingResponse {
  media: MediaItem;
  position: number;
  duration: number;
  finished: boolean;
}
export interface SearchParams {
  page?: number;
  limit?: number;
  search?: string;
  type?: string;
}
