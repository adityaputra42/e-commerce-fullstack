import axios, { AxiosInstance, AxiosResponse } from 'axios';

// Types
export interface User {
  id: number;
  email: string;
  name: string;
  role: string;
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
    color_varian?: any[]; // Simplified for now, can expand later
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
    // Assuming there's an endpoint to get profile, or we decode token. 
    // For now we might just rely on stored user data or implement a /me endpoint if backend supports it.
    // If not, we might skipped this or use what we stored.
    // Let's assume we store user in localStorage on login.
    return JSON.parse(localStorage.getItem('user') || 'null');
  },
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
  }
};

// ... (interfaces)

export interface Address {
  id: number;
  name: string;
  phone: string;
  street: string;
  city: string;
  state: string;
  zip_code: string;
  is_primary: boolean;
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
    const response = await api.get('/address');
    return response.data.data || [];
  },
  create: async (data: any): Promise<Address> => {
    const response = await api.post('/address', data);
    return response.data.data;
  },
  delete: async (id: number): Promise<void> => {
    await api.delete(`/address/${id}`);
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
