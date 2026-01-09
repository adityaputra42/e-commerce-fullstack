"use client";

import React from 'react';
import { AuthProvider } from '@/context/AuthContext';
import { Toaster } from '@/components/ui/sonner';
import { ThemeProvider } from "next-themes";
import { CartProvider } from '@/context/CartContext';

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <AuthProvider>
      <ThemeProvider
        attribute="class"
        defaultTheme="dark"
        enableSystem={false}
        disableTransitionOnChange
      >
        <CartProvider>
          {children}
          <Toaster />
        </CartProvider>
      </ThemeProvider>
    </AuthProvider>
  );
}
