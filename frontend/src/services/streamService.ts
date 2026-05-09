const BACKEND_BASE_URL = import.meta.env.DEV ? 'http://localhost:9000/api' : '/api';


export function getStreamUrl(mediaId: string, startTime?: number): string {
  let url = `${BACKEND_BASE_URL}/stream/${mediaId}`;
  if (startTime && startTime > 0) {
    url += `?start=${Math.floor(startTime)}`;
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
