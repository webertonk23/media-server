import { API_BASE_URL } from '@/utils/constants';

export function getStreamUrl(mediaId: string, startTime?: number, audioIndex?: number): string {
  let url = `${API_BASE_URL}/stream/${mediaId}`;
  const params = [];
  if (startTime && startTime > 0) {
    params.push(`start=${Math.floor(startTime)}`);
  }
  if (audioIndex !== undefined) {
    params.push(`audio=${audioIndex}`);
  }
  if (params.length > 0) {
    url += `?${params.join('&')}`;
  }
  return url;
}

export function createRangeHeaders(start: number, end?: number): Record<string, string> {
  const rangeValue = end !== undefined ? `bytes=${start}-${end}` : `bytes=${start}-`;
  return {
    'Range': rangeValue,
  };
}
export function parseContentRange(contentRange: string): { start: number; end: number; total: number } | null {
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
