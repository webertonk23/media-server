/**
 * Progress Service
 * 
 * Provides methods for fetching and saving video playback progress.
 * Handles progress retrieval and persistence for continue watching functionality.
 * 
 * **Validates: Requirements 6.3, 6.4**
 */

import apiClient from './api';
import type { ProgressData } from '@/types/player';
import type { UpdateProgressRequest } from '@/types/api';

/**
 * Get saved progress for a specific media item
 * 
 * Fetches playback progress from GET /progress/:id endpoint.
 * Returns null if no progress has been saved for this media item.
 * 
 * **Validates: Requirement 6.3**
 * 
 * @param mediaId - ULID of the media item
 * @returns Promise resolving to progress data or null if no progress exists
 * @throws ApiError if request fails (except 404 which returns null)
 * 
 * @example
 * ```typescript
 * const progress = await getProgress('01HQXYZ123ABC456DEF789GHI');
 * if (progress) {
 *   console.log(`Resume from ${progress.position}s`);
 * }
 * ```
 */
export async function getProgress(mediaId: string): Promise<ProgressData | null> {
  try {
    const response = await apiClient.get<ProgressData>(`/progress/${mediaId}`);
    return response.data;
  } catch (error: any) {
    // Return null if no progress exists (404 Not Found)
    if (error.status === 404) {
      return null;
    }
    
    // Re-throw other errors
    throw error;
  }
}

/**
 * Save playback progress for a specific media item
 * 
 * Saves current playback position to POST /progress/:id endpoint.
 * Handles errors gracefully by logging them without throwing.
 * 
 * **Validates: Requirement 6.4**
 * 
 * @param mediaId - ULID of the media item
 * @param data - Progress data to save (position, duration, finished)
 * @returns Promise that resolves when progress is saved
 * 
 * @example
 * ```typescript
 * // Save progress at 120 seconds of a 3600 second video
 * await saveProgress('01HQXYZ123ABC456DEF789GHI', {
 *   position: 120,
 *   duration: 3600,
 *   finished: false
 * });
 * 
 * // Mark video as finished (>= 95% watched)
 * await saveProgress('01HQXYZ123ABC456DEF789GHI', {
 *   position: 3420,
 *   duration: 3600,
 *   finished: true
 * });
 * ```
 */
export async function saveProgress(
  mediaId: string,
  data: UpdateProgressRequest
): Promise<void> {
  try {
    await apiClient.post(`/progress/${mediaId}`, data);
  } catch (error: any) {
    // Log error but don't throw - progress save failures should not disrupt playback
    console.error('[Progress Service] Failed to save progress:', {
      mediaId,
      position: data.position,
      duration: data.duration,
      finished: data.finished,
      error: error.message || error,
    });
    
    // Optionally, could implement retry logic here or queue for later
    // For now, we gracefully handle the error by logging only
  }
}
