// dashboard.types.ts

export interface DashboardStats {
  total_users: number;
  active_users: number;
  inactive_users: number;
  new_users_today: number;
  new_users_this_week: number;
}

export interface RevenueStats {
  total_revenue: number;
  revenue_today: number;
  revenue_this_month: number;
  revenue_growth_percent: number;
}

export interface OrderStats {
  total_orders: number;
  completed_orders: number;
  pending_orders: number;
  cancelled_orders: number;
}

export interface Order {
  id: number;
  order_number: string;
  user_id: number;
  total_amount: number;
  status: 'pending' | 'completed' | 'cancelled';
  created_at: string;
}

export interface ProductSummary {
  id: number;
  name: string;
  total_sold: number;
}

export interface LowStockProduct {
  id: number;
  name: string;
  stock: number;
}

export interface OrderAnalytics {
  date: string;
  total_orders: number;
  total_revenue: number;
}

export interface UserAnalytics {
  date: string;
  user_count: number;
}

export interface ActivityLog {
  id: number;
  action: string;
  resource: string;
  ip_address: string;
  user_agent: string;
  created_at: string;
  user: {
    id: number;
    username: string;
    email: string;
  };
}

export interface DashboardState {
  stats: DashboardStats | null;
  revenue: RevenueStats | null;
  orderStats: OrderStats | null;

  recentOrders: Order[];
  topProducts: ProductSummary[];
  lowStockProducts: LowStockProduct[];

  orderAnalytics: OrderAnalytics[];
  userAnalytics: UserAnalytics[];

  recentActivity: ActivityLog[];

  isLoading: boolean;
  error: string | null;
}
