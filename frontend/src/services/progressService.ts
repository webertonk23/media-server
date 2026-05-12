import apiClient from './api';
import type { ProgressData } from '@/types/player';
import type { UpdateProgressRequest, ContinueWatchingResponse } from '@/types/api';
export async function getProgress(mediaId: string): Promise<ProgressData | null> {
  try {
    const response = await apiClient.get<ProgressData>(`/progress/${mediaId}`);
    return response.data;
  } catch (error: any) {
    if (error.status === 404) {
      return null;
    }
    throw error;
  }
}
export async function getContinueWatching(): Promise<ContinueWatchingResponse[]> {
  const response = await apiClient.get<ContinueWatchingResponse[]>('/media/continue-watching');
  return response.data;
}
export async function saveProgress(
  mediaId: string,
  data: UpdateProgressRequest
): Promise<void> {
  try {
    await apiClient.post(`/progress/${mediaId}`, data);
  } catch (error: any) {
    console.error('[Progress Service] Failed to save progress:', {
      mediaId,
      position: data.position,
      duration: data.duration,
      finished: data.finished,
      error: error.message || error,
    });
  }
}
