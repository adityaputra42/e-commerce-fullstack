import { Users, UserPlus, Shield } from 'lucide-react';
import {
  PieChart,
  Pie,
  Cell,
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  CartesianGrid,
} from 'recharts';

import { useDashboardData } from '../hooks/useDashboardStats';
import { useAuthStore } from '../hooks/useAuth';

/* ======================
 * Constants & Helpers
 * ====================== */

const COLORS = [
  '#0088FE',
  '#00C49F',
  '#FFBB28',
  '#FF8042',
  '#8884d8',
  '#82ca9d',
];

const renderRoleLabel = ({ role_name, percent }: any) =>
  `${role_name} (${(percent * 100).toFixed(0)}%)`;

const getDisplayName = (user: any) =>
  user?.first_name ||
  user?.username ||
  user?.email ||
  'Unknown';

/* ======================
 * Component
 * ====================== */

const DashboardPage = () => {
  /**
   * 🔐 AUTH STATE (WAJIB)
   */
  const { isInitialized, isAuthenticated } = useAuthStore();

  /**
   * 📊 DASHBOARD DATA
   */
  const {
    stats,
    roleDistribution,
    recentActivity,
    userAnalytics,
    isLoading,
    error,
  } = useDashboardData();

  /* ======================
   * AUTH GATE (FIX BLANK)
   * ====================== */

  // ⛔ Tunggu Zustand hydrate selesai
  if (!isInitialized) {
    return (
      <div className="h-screen flex items-center justify-center text-gray-500">
        Initializing dashboard...
      </div>
    );
  }

  // ⛔ Safety net (biasanya sudah di ProtectedRoute)
  if (!isAuthenticated) {
    return null;
  }

  /* ======================
   * Loading & Error
   * ====================== */

  if (isLoading) {
    return (
      <div className="p-6 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {Array.from({ length: 4 }).map((_, i) => (
          <div
            key={i}
            className="h-28 bg-gray-200 animate-pulse rounded-lg"
          />
        ))}
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-6 text-center text-red-600">
        Failed to load dashboard data. Please try again.
      </div>
    );
  }

  /* ======================
   * Derived Values
   * ====================== */

  const totalUsers = stats?.total_users ?? 0;
  const activeUsers = stats?.active_users ?? 0;
  const newUsersToday = stats?.new_users_today ?? 0;
  const totalRoles = stats?.total_roles ?? 0;

  /* ======================
   * Render
   * ====================== */

  return (
    <div className="p-6 bg-gray-100 min-h-screen">
      <h1 className="text-3xl font-bold text-gray-800 mb-6">
        Dashboard Overview
      </h1>

      {/* ======================
       * Stats Cards
       * ====================== */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <StatsCard
          title="Total Users"
          value={totalUsers}
          icon={<Users className="w-6 h-6" />}
          color="blue"
        />
        <StatsCard
          title="Active Users"
          value={activeUsers}
          icon={<UserPlus className="w-6 h-6" />}
          color="green"
        />
        <StatsCard
          title="New Registrations"
          value={newUsersToday}
          icon={<UserPlus className="w-6 h-6" />}
          color="yellow"
        />
        <StatsCard
          title="Total Roles"
          value={totalRoles}
          icon={<Shield className="w-6 h-6" />}
          color="red"
        />
      </div>

      {/* ======================
       * Charts
       * ====================== */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        {/* Role Distribution */}
        <Card title="User Role Distribution">
          {roleDistribution?.length ? (
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={roleDistribution}
                  dataKey="user_count"
                  nameKey="role_name"
                  cx="50%"
                  cy="50%"
                  outerRadius={100}
                  label={renderRoleLabel}
                  labelLine={false}
                >
                  {roleDistribution.map((_, index) => (
                    <Cell
                      key={index}
                      fill={COLORS[index % COLORS.length]}
                    />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <EmptyState message="No role distribution data available." />
          )}
        </Card>

        {/* User Analytics */}
        <Card title="New User Analytics">
          {userAnalytics?.length ? (
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={userAnalytics}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar
                  dataKey="user_count"
                  name="New Users"
                  fill={COLORS[0]}
                />
              </BarChart>
            </ResponsiveContainer>
          ) : (
            <EmptyState message="No user analytics data available." />
          )}
        </Card>
      </div>

      {/* ======================
       * Recent Activity
       * ====================== */}
      <Card title="Recent Activity">
        {recentActivity?.length ? (
          <ul className="divide-y divide-gray-200">
            {recentActivity.map((activity: any) => (
              <li
                key={activity.id}
                className="py-3 flex flex-col md:flex-row md:justify-between gap-1"
              >
                <span className="text-gray-700 break-words">
                  User{' '}
                  <span className="font-medium">
                    {getDisplayName(activity.user)}
                  </span>{' '}
                  {activity.action} {activity.resource} from{' '}
                  {activity.ip_address}
                </span>
                <span className="text-sm text-gray-500">
                  {new Date(activity.created_at).toLocaleString()}
                </span>
              </li>
            ))}
          </ul>
        ) : (
          <EmptyState message="No recent activity." />
        )}
      </Card>
    </div>
  );
};

export default DashboardPage;

/* ======================
 * Reusable Components
 * ====================== */

const Card = ({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) => (
  <div className="bg-white p-6 rounded-lg shadow-md">
    <h2 className="text-xl font-semibold text-gray-800 mb-4">
      {title}
    </h2>
    {children}
  </div>
);

const StatsCard = ({
  title,
  value,
  icon,
  color,
}: {
  title: string;
  value: number;
  icon: React.ReactNode;
  color: 'blue' | 'green' | 'yellow' | 'red';
}) => {
  const colorMap: any = {
    blue: 'bg-blue-100 text-blue-600',
    green: 'bg-green-100 text-green-600',
    yellow: 'bg-yellow-100 text-yellow-600',
    red: 'bg-red-100 text-red-600',
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-md flex justify-between items-center">
      <div>
        <p className="text-sm font-medium text-gray-500">{title}</p>
        <p className="text-3xl font-bold text-gray-900">{value}</p>
      </div>
      <div className={`p-3 rounded-full ${colorMap[color]}`}>
        {icon}
      </div>
    </div>
  );
};

const EmptyState = ({ message }: { message: string }) => (
  <p className="text-gray-500">{message}</p>
);
