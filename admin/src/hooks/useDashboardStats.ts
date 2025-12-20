import { useEffect, useState } from 'react';
import api from '../services/api';
import { useAuthStore } from './useAuth';

export const useDashboardData = () => {
  const { isAuthenticated, isInitialized } = useAuthStore();

  const [stats, setStats] = useState<any>(null);
  const [roleDistribution, setRoleDistribution] = useState<any[]>([]);
  const [recentActivity, setRecentActivity] = useState<any[]>([]);
  const [userAnalytics, setUserAnalytics] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isInitialized || !isAuthenticated) {
      setIsLoading(false);
      return;
    }

    const fetchDashboard = async () => {
      try {
        setIsLoading(true);
        setError(null);

     
        const [usersRes, rolesRes] = await Promise.all([
          api.get('/users').catch(() => null),
          api.get('/roles').catch(() => null),
        ]);

        setStats({
          total_users: usersRes?.data?.data?.length ?? 0,
          active_users: usersRes?.data?.data?.filter((u: any) => u.is_active)?.length ?? 0,
          inactive_users: usersRes?.data?.data?.filter((u: any) => !u.is_active)?.length ?? 0,
          new_users_today: 0,
          new_users_this_week: 0,
          total_roles: rolesRes?.data?.data?.length ?? 0,
        });

        setRoleDistribution(
          rolesRes?.data?.data?.map((r: any) => ({
            role_name: r.name,
            user_count: r.users?.length ?? 0,
          })) ?? []
        );

        setRecentActivity([]); // backend belum ada
        setUserAnalytics([]);  // backend belum ada
      } catch (err) {
        console.error('Dashboard fetch error:', err);
        setError('Failed to load dashboard data');
      } finally {
        setIsLoading(false);
      }
    };

    fetchDashboard();
  }, [isAuthenticated, isInitialized]);

  return {
    stats,
    roleDistribution,
    recentActivity,
    userAnalytics,
    isLoading,
    error,
  };
};
