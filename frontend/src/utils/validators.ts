/**
 * Utility functions for validating data
 * 
 * **Validates: Requirements 6.3, 16.3**
 */
/**
 * Validate if a string is a valid ULID format
 * ULIDs are 26-character strings using Crockford's base32 alphabet
 * 
 * @param id - String to validate as ULID
 * @returns true if valid ULID format, false otherwise
 * 
 * @example
 * isValidMediaId("01ARZ3NDEKTSV4RRFFQ69G5FAV") 
 * isValidMediaId("invalid") 
 * isValidMediaId("") 
 */
export function isValidMediaId(id: string): boolean {
  if (!id || typeof id !== 'string') {
    return false;
  }
  const ulidRegex = /^[0123456789ABCDEFGHJKMNPQRSTVWXYZ]{26}$/;
  return ulidRegex.test(id);
}
/**
 * Validate if progress position and duration are valid numbers
 * Position must be non-negative and not exceed duration
 * Duration must be positive
 * 
 * @param position - Current playback position in seconds
 * @param duration - Total video duration in seconds
 * @returns true if both values are valid, false otherwise
 * 
 * @example
 * isValidProgress(100, 200) 
 * isValidProgress(0, 100) 
 * isValidProgress(-10, 100) 
 * isValidProgress(150, 100) 
 * isValidProgress(50, 0) 
 * isValidProgress(50, -100) 
 */
export function isValidProgress(position: number, duration: number): boolean {
  if (!isFinite(position) || !isFinite(duration)) {
    return false;
  }
  if (position < 0) {
    return false;
  }
  if (duration <= 0) {
    return false;
  }
  if (position > duration) {
    return false;
  }
  return true;
}
