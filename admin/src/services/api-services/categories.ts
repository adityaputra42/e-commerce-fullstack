import api from '../api-client';
import type { ApiResponse } from '../../types/api';

export interface Category {
  id: number;
  name: string;
}

export interface CategoryListResponse {
  categories: Category[];
  total: number;
  page: number;
  limit: number;
}

/**
 * Categories API Service
 * Handles all category-related API calls
 */
export const categoriesApi = {
  /**
   * Get all categories
   * GET /categories
   */
  async getCategories(limit: number = 100): Promise<Category[]> {
    const response = await api.get<ApiResponse<any>>('/categories', {
      params: { limit },
    });
    // Handle both { data: { categories: [] } } and { data: [] }
    const data = response.data?.data?.categories || response.data?.data || [];
    return Array.isArray(data) ? data : [];
  },
};
