/**
 * API-related types for HTTP requests and responses
 * 
 * **Validates: Requirements 20.1, 16.3**
 */

/**
 * ApiError interface represents error responses from the backend API
 * Matches the backend error response structure: { "error": "message" }
 */
export interface ApiError {
  /** Error message describing what went wrong */
  message: string;
  
  /** HTTP status code (e.g., 404, 500) */
  status: number;
  
  /** Optional error code for specific error types */
  code?: string;
}

/**
 * UpdateProgressRequest interface represents the payload for saving playback progress
 * Matches the backend UpdateProgressRequest DTO
 */
export interface UpdateProgressRequest {
  /** Current playback position in seconds */
  position: number;
  
  /** Total video duration in seconds */
  duration: number;
  
  /** Whether the video has been fully watched (>= 95% completion) */
  finished: boolean;
}

/**
 * SearchParams interface defines query parameters for media search and filtering
 * Used with GET /media endpoint
 */
export interface SearchParams {
  /** Page number for pagination (default: 1) */
  page?: number;
  
  /** Number of items per page (default: 20) */
  limit?: number;
  
  /** Search query string to filter by title */
  search?: string;
  
  /** Media type filter: 'movie', 'series', or 'episode' */
  type?: string;
}
