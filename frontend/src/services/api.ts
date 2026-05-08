/**
 * Axios HTTP Client Configuration
 * 
 * Configures the Axios instance with base URL, interceptors, and retry logic
 * for communicating with the Go backend API.
 * 
 * **Validates: Requirements 16.1, 16.2, 16.4, 16.5, 14.1, 14.2, 14.3**
 */

import axios from 'axios';
import type { AxiosError, AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import type { ApiError } from '@/types/api';

/**
 * Base URL for the backend API
 * In development: Use full URL to backend server
 * In production: Use relative URL (served from same origin)
 */
const BASE_URL = import.meta.env.DEV ? 'http://localhost:9000/api' : '/api';

/**
 * Maximum number of retry attempts for failed requests
 */
const MAX_RETRIES = 3;

/**
 * Initial delay for exponential backoff (in milliseconds)
 */
const INITIAL_RETRY_DELAY = 1000;

/**
 * Request timeout in milliseconds
 */
const REQUEST_TIMEOUT = 30000;

/**
 * Create and configure the Axios instance
 */
const apiClient: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  timeout: REQUEST_TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
    'Accept': 'application/json',
  },
});

/**
 * Request interceptor to add custom headers and configuration
 * 
 * **Validates: Requirement 16.2**
 */
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Add timestamp to prevent caching issues
    config.headers['X-Request-Time'] = new Date().toISOString();
    
    // Initialize retry count if not present
    if (!config.headers['X-Retry-Count']) {
      config.headers['X-Retry-Count'] = '0';
    }
    
    return config;
  },
  (error: AxiosError) => {
    console.error('[API Request Error]', error);
    return Promise.reject(error);
  }
);

/**
 * Response interceptor for error handling and retry logic
 * 
 * Handles:
 * - 401 Unauthorized
 * - 404 Not Found
 * - 500 Internal Server Error
 * - Network errors
 * - Timeout errors with exponential backoff retry
 * 
 * **Validates: Requirements 16.4, 16.5, 14.1, 14.2, 14.3**
 */
apiClient.interceptors.response.use(
  (response: AxiosResponse) => {
    // Successful response, return as-is
    return response;
  },
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retryCount?: number };
    
    // Handle different error scenarios
    if (error.response) {
      // Server responded with error status
      const status = error.response.status;
      const apiError: ApiError = {
        message: getErrorMessage(status, error),
        status: status,
        code: error.code,
      };
      
      // Log error for debugging (Requirement 14.5)
      console.error(`[API Error ${status}]`, {
        url: originalRequest?.url,
        method: originalRequest?.method,
        status,
        message: apiError.message,
        data: error.response.data,
      });
      
      // Handle specific status codes
      switch (status) {
        case 401:
          // Unauthorized - could redirect to login if auth is implemented
          console.warn('[API] Unauthorized access - authentication required');
          break;
        case 404:
          // Not Found - handled by components
          break;
        case 500:
          // Internal Server Error
          console.error('[API] Server error occurred');
          break;
      }
      
      return Promise.reject(apiError);
    } else if (error.request) {
      // Request was made but no response received (network error)
      const isTimeout = error.code === 'ECONNABORTED' || error.message.includes('timeout');
      
      // Implement retry logic with exponential backoff for timeouts
      if (isTimeout && originalRequest) {
        const retryCount = originalRequest._retryCount || 0;
        
        if (retryCount < MAX_RETRIES) {
          originalRequest._retryCount = retryCount + 1;
          
          // Calculate exponential backoff delay
          const delay = INITIAL_RETRY_DELAY * Math.pow(2, retryCount);
          
          console.warn(`[API] Timeout - Retrying request (${retryCount + 1}/${MAX_RETRIES}) after ${delay}ms`, {
            url: originalRequest.url,
            method: originalRequest.method,
          });
          
          // Update retry count header
          originalRequest.headers['X-Retry-Count'] = String(retryCount + 1);
          
          // Wait for the backoff delay
          await new Promise(resolve => setTimeout(resolve, delay));
          
          // Retry the request
          return apiClient(originalRequest);
        } else {
          console.error('[API] Max retries reached for timeout');
        }
      }
      
      // Network error or max retries exceeded
      const apiError: ApiError = {
        message: isTimeout 
          ? 'Tempo limite excedido. Verifique sua conexão e tente novamente.'
          : 'Sem conexão com o servidor. Verifique sua conexão de rede.',
        status: 0,
        code: error.code,
      };
      
      console.error('[API Network Error]', {
        message: apiError.message,
        code: error.code,
        url: originalRequest?.url,
      });
      
      return Promise.reject(apiError);
    } else {
      // Something else happened
      const apiError: ApiError = {
        message: 'Erro inesperado ao processar a requisição',
        status: 0,
        code: error.code,
      };
      
      console.error('[API Unexpected Error]', error.message);
      
      return Promise.reject(apiError);
    }
  }
);

/**
 * Get user-friendly error message based on status code
 * 
 * **Validates: Requirements 14.1, 14.2, 14.3**
 */
function getErrorMessage(status: number, error: AxiosError): string {
  switch (status) {
    case 401:
      return 'Não autorizado. Autenticação necessária.';
    case 404:
      return 'Mídia não encontrada';
    case 500:
      return 'Erro no servidor, tente novamente';
    case 503:
      return 'Serviço temporariamente indisponível';
    default:
      // Try to extract error message from response
      const responseData = error.response?.data as any;
      if (responseData?.error) {
        return responseData.error;
      }
      return `Erro na requisição (${status})`;
  }
}

/**
 * Export the configured Axios instance
 */
export default apiClient;

/**
 * Export utility function to check if error is an ApiError
 */
export function isApiError(error: any): error is ApiError {
  return error && typeof error.message === 'string' && typeof error.status === 'number';
}
