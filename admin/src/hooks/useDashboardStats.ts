import { useEffect, useState, useCallback } from 'react';
import api from '../services/api';
import type { DashboardState } from '../types/dashboard';

export const useDashboardData = () => {
  const [data, setData] = useState<DashboardState>({
    stats: null,
    revenue: null,
    orderStats: null,

    recentOrders: [],
    topProducts: [],
    lowStockProducts: [],

    orderAnalytics: [],
    userAnalytics: [],

    recentActivity: [],

    isLoading: true,
    error: null,
  });

  const fetchDashboard = useCallback(async () => {
    setData((prev) => ({
      ...prev,
      isLoading: true,
      error: null,
    }));

    try {
      const results = await Promise.allSettled([
        api.get('/dashboard/stats'),
        api.get('/dashboard/revenue'),
        api.get('/dashboard/orders/stats'),
        api.get('/dashboard/orders/recent'),
        api.get('/dashboard/products/top'),
        api.get('/dashboard/products/low-stock'),
        api.get('/dashboard/analytics/orders'),
        api.get('/dashboard/analytics/users'),
        api.get('/dashboard/activity'),
      ]);

      const [
        statsRes,
        revenueRes,
        orderStatsRes,
        recentOrdersRes,
        topProductsRes,
        lowStockRes,
        orderAnalyticsRes,
        userAnalyticsRes,
        activityRes,
      ] = results;

      setData({
        stats:
          statsRes.status === 'fulfilled'
            ? statsRes.value.data?.data ?? null
            : null,

        revenue:
          revenueRes.status === 'fulfilled'
            ? revenueRes.value.data?.data ?? null
            : null,

        orderStats:
          orderStatsRes.status === 'fulfilled'
            ? orderStatsRes.value.data?.data ?? null
            : null,

        recentOrders:
          recentOrdersRes.status === 'fulfilled'
            ? recentOrdersRes.value.data?.data ?? []
            : [],

        topProducts:
          topProductsRes.status === 'fulfilled'
            ? topProductsRes.value.data?.data ?? []
            : [],

        lowStockProducts:
          lowStockRes.status === 'fulfilled'
            ? lowStockRes.value.data?.data ?? []
            : [],

        orderAnalytics:
          orderAnalyticsRes.status === 'fulfilled'
            ? orderAnalyticsRes.value.data?.data ?? []
            : [],

        userAnalytics:
          userAnalyticsRes.status === 'fulfilled'
            ? userAnalyticsRes.value.data?.data ?? []
            : [],

        recentActivity:
          activityRes.status === 'fulfilled'
            ? activityRes.value.data?.data ?? []
            : [],

        isLoading: false,
        error: null,
      });
    } catch (err) {
      console.error('Dashboard fetch error:', err);
      setData((prev) => ({
        ...prev,
        isLoading: false,
        error: 'Failed to load dashboard data',
      }));
    }
  }, []);

  useEffect(() => {
    fetchDashboard();
  }, [fetchDashboard]);

  return {
    ...data,
    refetch: fetchDashboard, // 🔥 optional but useful
  };
};
