"use client";

import React from 'react';
import { AuthProvider } from '@/context/AuthContext';
import { Toaster } from 'sonner';
import { CartProvider } from '@/context/CartContext';

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider>
      <CartProvider>
        {children}
        <Toaster position="top-right" richColors />
      </CartProvider>
    </AuthProvider>
  );
}
