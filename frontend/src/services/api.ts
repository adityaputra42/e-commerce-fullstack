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
  token: string;
  user: User;
}

export interface RegisterResponse {
  user: User;
  token: string;
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

const api: AxiosInstance = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add a request interceptor to attach the token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export const authService = {
  login: async (email: string, password: string): Promise<LoginResponse> => {
    const response = await api.post('/auth/login', { email, password });
    return response.data;
  },
  register: async (userData: any): Promise<RegisterResponse> => {
     const response = await api.post('/auth/register', userData);
     return response.data;
  },
  getCurrentUser: async (): Promise<User> => {
    const response = await api.get('/users/me');
    return response.data.data;
  },
  updateProfile: async (userData: Partial<User>): Promise<User> => {
    // Ideally use /users/me but backend might expect ID. 
    // Let's safe check user ID availability or use /users/me if supported for updates
    // Based on handler, UpdateUser needs ID.
    const currentUser = JSON.parse(localStorage.getItem('user') || '{}');
    if (!currentUser.id) throw new Error("User ID not found");
    
    const response = await api.put(`/users/${currentUser.id}`, userData);
    return response.data.data;
  },
  changePassword: async (data: any): Promise<void> => {
    await api.put('/users/me/password', data);
  },
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
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

export default api;
