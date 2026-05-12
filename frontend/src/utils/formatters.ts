/**
 * Utility functions for formatting data for display
 * 
 * **Validates: Requirements 4.2, 5.3**
 */
/**
 * Format duration in seconds to human-readable format (e.g., "1h 23m", "45m", "2h 15m")
 * Used for displaying video length in media cards and details
 * 
 * @param seconds - Duration in seconds
 * @returns Formatted duration string (e.g., "1h 23m")
 * 
 * @example
 * formatDuration(5000) 
 * formatDuration(2700) 
 * formatDuration(7800) 
 * formatDuration(45) 
 */
export function formatDuration(seconds: number): string {
  if (seconds < 0 || !isFinite(seconds)) {
    return '0m';
  }
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  if (hours > 0) {
    return `${hours}h ${minutes}m`;
  }
  return `${minutes}m`;
}
/**
 * Format time in seconds to HH:MM:SS or MM:SS format
 * Used for displaying current playback position in video player
 * 
 * @param seconds - Time in seconds
 * @returns Formatted time string (e.g., "01:23:45" or "23:45")
 * 
 * @example
 * formatTime(5025) 
 * formatTime(1425) 
 * formatTime(45) 
 */
export function formatTime(seconds: number): string {
  if (seconds < 0 || !isFinite(seconds)) {
    return '00:00';
  }
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = Math.floor(seconds % 60);
  const paddedMinutes = String(minutes).padStart(2, '0');
  const paddedSeconds = String(secs).padStart(2, '0');
  if (hours > 0) {
    const paddedHours = String(hours).padStart(2, '0');
    return `${paddedHours}:${paddedMinutes}:${paddedSeconds}`;
  }
  return `${paddedMinutes}:${paddedSeconds}`;
}
/**
 * Extract and format year from date string
 * Used for displaying release year in media cards and details
 * 
 * @param date - Date string in ISO format or year string
 * @returns Formatted year string (e.g., "2022")
 * 
 * @example
 * formatYear("2022-05-06") 
 * formatYear("2022") 
 * formatYear("invalid") 
 */
export function formatYear(date: string): string {
  if (!date || typeof date !== 'string') {
    return '';
  }
  const yearMatch = date.match(/^\d{4}/);
  if (yearMatch) {
    return yearMatch[0];
  }
  const dateObj = new Date(date);
  if (!isNaN(dateObj.getTime())) {
    return String(dateObj.getFullYear());
  }
  return '';
}
