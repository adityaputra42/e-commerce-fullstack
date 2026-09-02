import axios, { AxiosInstance, AxiosResponse } from 'axios';

// Types
export interface User {
  id: number;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  role_id: number;
  is_active: boolean;
  role: {
    id: number;
    name: string;
    permissions: any[];
  };
}

export interface LoginResponse {
  access_token: string;
  expires_at: string;
  user: User;
}

export interface RegisterResponse {
  access_token: string;
  expires_at: string;
  user: User;
}

export interface Product {
    id: number;
    name: string;
    description: string;
    price: number;
    images: string;
    rating?: number;
    category: {
        id: number;
        name: string;
    };
    color_varian?: any[];
}

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
let inMemoryAccessToken: string | null = null;

export function setAccessToken(token: string | null) {
  inMemoryAccessToken = token;
}

export function getAccessToken(): string | null {
  return inMemoryAccessToken;
}

const api: AxiosInstance = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  // WAJIB true supaya browser ikut mengirim cookie httpOnly refresh_token
  // ke backend. Backend juga harus set Access-Control-Allow-Credentials
  // (sudah dikonfigurasi di CORS backend) dan origin spesifik (bukan "*").
  withCredentials: true,
});

// Request interceptor: tempel access token dari memory (bukan localStorage).
api.interceptors.request.use(
  (config) => {
    if (inMemoryAccessToken) {
      config.headers.Authorization = `Bearer ${inMemoryAccessToken}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor: kalau access token expired (401), coba refresh
// sekali lewat cookie httpOnly, lalu ulangi request yang gagal. Kalau
// refresh juga gagal (cookie tidak ada/expired), biarkan 401 apa adanya —
// AuthContext yang akan menangani logout.
let refreshPromise: Promise<string | null> | null = null;

async function refreshAccessToken(): Promise<string | null> {
  if (!refreshPromise) {
    refreshPromise = axios
      .post(`${API_URL}/auth/refresh`, {}, { withCredentials: true })
      .then((res) => {
        const newToken = res.data?.data?.access_token ?? null;
        setAccessToken(newToken);
        return newToken;
      })
      .catch(() => {
        setAccessToken(null);
        return null;
      })
      .finally(() => {
        refreshPromise = null;
      });
  }
  return refreshPromise;
}

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    if (
      error.response?.status === 401 &&
      !originalRequest._retry &&
      !originalRequest.url?.includes('/auth/')
    ) {
      originalRequest._retry = true;
      const newToken = await refreshAccessToken();
      if (newToken) {
        originalRequest.headers.Authorization = `Bearer ${newToken}`;
        return api(originalRequest);
      }
    }
    return Promise.reject(error);
  }
);

export const authService = {
  login: async (email: string, password: string): Promise<LoginResponse> => {
    const response = await api.post('/auth/login', { email, password });
    return response.data.data;
  },
  register: async (userData: any): Promise<RegisterResponse> => {
     const response = await api.post('/auth/register', userData);
     return response.data.data;
  },
  // refresh menukar cookie httpOnly refresh_token jadi access token baru.
  // Dipanggil AuthContext saat aplikasi dimuat (page load/refresh) karena
  // access token di memory hilang setiap reload.
  refresh: async (): Promise<{ access_token: string } | null> => {
    try {
      const response = await api.post('/auth/refresh');
      return response.data.data;
    } catch {
      return null;
    }
  },
  getCurrentUser: async (): Promise<User> => {
    const response = await api.get('/users/me');
    return response.data.data;
  },
  updateProfile: async (userId: number, userData: Partial<User>): Promise<User> => {
    const response = await api.put(`/users/${userId}`, userData);
    return response.data.data;
  },
  changePassword: async (data: any): Promise<void> => {
    await api.put('/users/me/password', data);
  },
  logout: async (): Promise<void> => {
    try {
      await api.post('/auth/logout');
    } finally {
      setAccessToken(null);
    }
  }
};

// ... (interfaces)

export interface Address {
  id: number;
  recipient_name: string;
  recipient_phone_number: string;
  full_address: string;
  city: string;
  province: string;
  district: string;
  village: string;
  postal_code: string;
}

export interface ShippingMethod {
  id: number;
  name: string;
  cost: number;
  estimated_delivery: string;
}

export interface PaymentMethod {
  id: number;
  name: string;
  type: string;
}

export interface Order {
    id: string;
    product: Product;
    quantity: number;
    unit_price: number;
    subtotal: number;
    status: string;
}

export interface Transaction {
    tx_id: string;
    address: Address;
    shipping: ShippingMethod;
    payment_method: PaymentMethod;
    shipping_price: number;
    total_price: number;
    status: string;
    orders: Order[];
    created_at: string;
}

// ... (existing services)

export const addressService = {
  getAll: async (): Promise<Address[]> => {
    const response = await api.get('/addresses');
    return response.data.data || [];
  },
  create: async (data: any): Promise<Address> => {
    const response = await api.post('/addresses', data);
    return response.data.data;
  },
  update: async (id: number, data: any): Promise<Address> => {
    const response = await api.put(`/addresses/${id}`, data);
    return response.data.data;
  },
  delete: async (id: number): Promise<void> => {
    await api.delete(`/addresses/${id}`);
  }
};

export const shippingService = {
  getAll: async (): Promise<ShippingMethod[]> => {
    const response = await api.get('/shipping');
    return response.data.data || [];
  }
};

export const paymentService = {
  getAllMethods: async (): Promise<PaymentMethod[]> => {
    const response = await api.get('/payment-methods');
    return response.data.data || [];
  }
};

export const transactionService = {
  createTransaction: async (data: any) => {
    const response = await api.post('/transactions', data);
    return response.data;
  },
  getAll: async () => {
      const response = await api.get('/transactions');
      return response.data.data;
  },
  getById: async (id: string) => {
      const response = await api.get(`/transactions/${id}`);
      return response.data.data;
  }
};

export const productService = {
    getAll: async (): Promise<Product[]> => {
        const response = await api.get('/products');
        return response.data.data || response.data; // Adjust based on actual API response structure
    },

    getById: async (id: string): Promise<Product> => {
        const response = await api.get(`/products/${id}`);
        return response.data.data || response.data;
    }
}

export interface Category {
  id: number;
  name: string;
  icon: string;
}

// categoryService mirrors productService's shape on purpose — one
// convention for "call the backend, unwrap `.data.data`, fall back to
// `.data`" instead of every page inventing its own unwrapping logic.
export const categoryService = {
  getAll: async (): Promise<Category[]> => {
    const response = await api.get('/categories');
    return response.data.data || response.data || [];
  },
};

export default api;
