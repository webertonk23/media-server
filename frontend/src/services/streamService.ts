/**
 * Stream Service
 * 
 * Provides methods for constructing streaming URLs and handling HTTP Range requests
 * for video playback with seeking support.
 * 
 * **Validates: Requirements 9.1, 9.2**
 */

/**
 * Base URL for the backend API
 * In development: Use full URL to backend server
 * In production: Use relative URL (served from same origin)
 */
const BACKEND_BASE_URL = import.meta.env.DEV ? 'http://localhost:9000/api' : '/api';

/**
 * Get the streaming URL for a specific media item
 * 
 * Constructs the URL for the backend streaming endpoint GET /stream/:id.
 * This URL is used by the HTML5 video player to fetch video content.
 * 
 * **Validates: Requirement 9.1**
 * 
 * @param mediaId - ULID of the media item to stream
 * @returns Complete streaming URL pointing to the backend endpoint
 * 
 * @example
 * ```typescript
 * const streamUrl = getStreamUrl('01HQXYZ123ABC456DEF789GHI');
 * // Returns: 'http://localhost:9000/api/stream/01HQXYZ123ABC456DEF789GHI' (in dev)
 * ```
 * // Use with HTML5 video element
 * videoElement.src = getStreamUrl(mediaId);
 * ```
 */
export function getStreamUrl(mediaId: string): string {
  return `${BACKEND_BASE_URL}/stream/${mediaId}`;
}

/**
 * Create HTTP Range request headers for video seeking
 * 
 * Generates the Range header for HTTP Range requests, enabling video seeking
 * and partial content delivery. The backend must support Range requests (206 Partial Content).
 * 
 * **Validates: Requirement 9.2**
 * 
 * @param start - Starting byte position (inclusive)
 * @param end - Ending byte position (inclusive, optional). If omitted, requests from start to end of file
 * @returns Headers object with Range header set
 * 
 * @example
 * ```typescript
 * // Request bytes 0-1023 (first 1KB)
 * const headers = createRangeHeaders(0, 1023);
 * // Returns: { 'Range': 'bytes=0-1023' }
 * 
 * // Request from byte 1024 to end of file
 * const headers = createRangeHeaders(1024);
 * // Returns: { 'Range': 'bytes=1024-' }
 * 
 * // Use with fetch API for seeking
 * fetch(streamUrl, { headers: createRangeHeaders(startByte, endByte) });
 * ```
 */
export function createRangeHeaders(start: number, end?: number): Record<string, string> {
  const rangeValue = end !== undefined ? `bytes=${start}-${end}` : `bytes=${start}-`;
  
  return {
    'Range': rangeValue,
  };
}

/**
 * Parse Content-Range header from response
 * 
 * Extracts range information from the Content-Range response header.
 * Useful for understanding what portion of the file was returned by the server.
 * 
 * Format: "bytes start-end/total"
 * Example: "bytes 0-1023/5242880" means bytes 0-1023 of a 5MB file
 * 
 * @param contentRange - Content-Range header value from response
 * @returns Parsed range information or null if invalid format
 * 
 * @example
 * ```typescript
 * const range = parseContentRange('bytes 0-1023/5242880');
 * // Returns: { start: 0, end: 1023, total: 5242880 }
 * 
 * const invalid = parseContentRange('invalid');
 * // Returns: null
 * ```
 */
export function parseContentRange(contentRange: string): { start: number; end: number; total: number } | null {
  // Match pattern: bytes start-end/total
  const match = contentRange.match(/bytes (\d+)-(\d+)\/(\d+)/);
  
  if (!match || !match[1] || !match[2] || !match[3]) {
    return null;
  }
  
  return {
    start: parseInt(match[1], 10),
    end: parseInt(match[2], 10),
    total: parseInt(match[3], 10),
  };
}
