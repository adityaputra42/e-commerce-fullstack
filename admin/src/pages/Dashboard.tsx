import { useDashboardData } from '../hooks/useDashboardStats';
import {
  Users,
  DollarSign,
  ShoppingCart,
  Activity,
} from 'lucide-react';
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
} from 'recharts';
import { StatCard } from '../components/common/StatCard';

const DashboardPage = () => {
  const {
    stats,
    revenue,
    orderStats,
    userAnalytics,
    recentActivity,
    isLoading,
    error,
  } = useDashboardData();

  if (isLoading) {
    return <div className="p-6 text-center">Loading dashboard data...</div>;
  }

  if (error) {
    return <div className="p-6 text-center text-red-500">{error}</div>;
  }

  return (
    <div className="p-6 bg-gray-100 min-h-screen">
      <h1 className="text-3xl font-bold text-gray-800 mb-6">
        Dashboard Overview
      </h1>

      {/* ===== STATS CARDS ===== */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatCard
          title="Total Users"
          value={stats?.total_users ?? 0}
          icon={<Users />}
          color="blue"
        />
        <StatCard
          title="Active Users"
          value={stats?.active_users ?? 0}
          icon={<Users />}
          color="green"
        />
        <StatCard
          title="Total Revenue"
          value={`$${revenue?.total_revenue ?? 0}`}
          icon={<DollarSign />}
          color="yellow"
        />
        <StatCard
          title="Total Orders"
          value={orderStats?.total_orders ?? 0}
          icon={<ShoppingCart />}
          color="red"
        />
      </div>

      {/* ===== ANALYTICS ===== */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-4">
            User Growth Analytics
          </h2>

          {userAnalytics.length > 0 ? (
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={userAnalytics}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis />
                <Tooltip />
                <Bar
                  dataKey="user_count"
                  name="New Users"
                  fill="#3b82f6"
                />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <p className="text-gray-500">No analytics data.</p>
          )}
        </div>

        {/* ===== RECENT ACTIVITY ===== */}
        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-xl font-semibold mb-4">
            Recent Activity
          </h2>

          {recentActivity.length > 0 ? (
            <ul className="divide-y">
              {recentActivity.map((activity) => (
                <li
                  key={activity.id}
                  className="py-3 flex justify-between"
                >
                  <span className="text-gray-700">
                    <span className="font-medium">
                      {activity.user?.username ??
                        activity.user?.email ??
                        'Unknown'}
                    </span>{' '}
                    {activity.action} {activity.resource}
                  </span>
                  <span className="text-sm text-gray-500">
                    {new Date(
                      activity.created_at
                    ).toLocaleString()}
                  </span>
                </li>
              ))}
            </ul>
          ) : (
            <p className="text-gray-500">No recent activity.</p>
          )}
        </div>
      </div>
    </div>
  );
};

export default DashboardPage;
