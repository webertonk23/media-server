import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { MediaItem, PaginatedResponse } from '@/types/media';
import type { SearchParams } from '@/types/api';
import * as mediaService from '@/services/mediaService';

interface PaginationState {
  page: number;
  limit: number;
  total: number;
  hasMore: boolean;
}

export const useMediaStore = defineStore('media', () => {
  const mediaItems = ref<MediaItem[]>([]);
  const currentMedia = ref<MediaItem | null>(null);
  const loading = ref<boolean>(false);
  const error = ref<string | null>(null);
  const pagination = ref<PaginationState>({
    page: 1,
    limit: 20,
    total: 0,
    hasMore: false,
  });

  const getMediaById = computed(() => {
    return (id: string): MediaItem | undefined => {
      return mediaItems.value.find(item => item.id === id);
    };
  });

  async function fetchMedia(params: SearchParams = {}, append: boolean = false): Promise<PaginatedResponse<MediaItem>> {
    loading.value = true;
    error.value = null;

    try {
      const response = await mediaService.getMedia(params);

      if (append) {
        mediaItems.value = [...mediaItems.value, ...response.items];
      } else {
        mediaItems.value = response.items;
      }

      pagination.value = {
        page: response.page,
        limit: response.limit,
        total: response.total,
        hasMore: response.page * response.limit < response.total,
      };

      return response;
    } catch (err: any) {
      error.value = err.message || 'Erro ao carregar mídia';
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function fetchMediaById(id: string): Promise<MediaItem> {
    loading.value = true;
    error.value = null;

    try {
      const cached = getMediaById.value(id);
      if (cached) {
        currentMedia.value = cached;
        loading.value = false;
        return cached;
      }

      const media = await mediaService.getMediaById(id);
      
      currentMedia.value = media;

      const existingIndex = mediaItems.value.findIndex(item => item.id === id);
      if (existingIndex === -1) {
        mediaItems.value.push(media);
      } else {
        mediaItems.value[existingIndex] = media;
      }

      return media;
    } catch (err: any) {
      error.value = err.message || 'Erro ao carregar detalhes da mídia';
      currentMedia.value = null;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function searchMedia(query: string, page: number = 1): Promise<PaginatedResponse<MediaItem>> {
    loading.value = true;
    error.value = null;

    try {
      const response = await mediaService.searchMedia(query, page);

      if (page === 1) {
        mediaItems.value = response.items;
      } else {
        mediaItems.value = [...mediaItems.value, ...response.items];
      }

      pagination.value = {
        page: response.page,
        limit: response.limit,
        total: response.total,
        hasMore: response.page * response.limit < response.total,
      };

      return response;
    } catch (err: any) {
      error.value = err.message || 'Erro ao buscar mídia';
      throw err;
    } finally {
      loading.value = false;
    }
  }

  function clearCurrentMedia(): void {
    currentMedia.value = null;
  }

  function clearAll(): void {
    mediaItems.value = [];
    currentMedia.value = null;
    error.value = null;
    pagination.value = {
      page: 1,
      limit: 20,
      total: 0,
      hasMore: false,
    };
  }

  return {
    mediaItems,
    currentMedia,
    loading,
    error,
    pagination,
    getMediaById,
    fetchMedia,
    fetchMediaById,
    searchMedia,
    clearCurrentMedia,
    clearAll,
  };
});

