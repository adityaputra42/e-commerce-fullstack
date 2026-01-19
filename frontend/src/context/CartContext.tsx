"use client";

import React, { createContext, useContext, useState, useEffect } from 'react';
import { Product, ColorVariant, SizeVariant } from '@/types/product';
import { toast } from 'sonner';

export interface CartItem {
  product: Product;
  quantity: number;
  selectedColorVariant?: ColorVariant;
  selectedSizeVariant?: SizeVariant;
}

interface CartContextType {
  items: CartItem[];
  addToCart: (product: Product, colorVariant?: ColorVariant, sizeVariant?: SizeVariant) => void;
  removeFromCart: (productId: number, colorVariantId?: number, sizeVariantId?: number) => void;
  updateQuantity: (productId: number, quantity: number, colorVariantId?: number, sizeVariantId?: number) => void;
  clearCart: () => void;
  cartTotal: number;
  cartCount: number;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export function CartProvider({ children }: { children: React.ReactNode }) {
  const [items, setItems] = useState<CartItem[]>([]);

  // Load cart from local storage on mount
  useEffect(() => {
    const storedCart = localStorage.getItem('cart');
    if (storedCart) {
      try {
        const parsedCart = JSON.parse(storedCart) as CartItem[];
        // Filter out invalid items (items without product data)
        const validItems = parsedCart.filter(item => 
          item && item.product && typeof item.product.price === 'number'
        );
        setItems(validItems);
        
        // If we filtered out invalid items, update localStorage
        if (validItems.length !== parsedCart.length) {
          console.warn(`Removed ${parsedCart.length - validItems.length} invalid items from cart`);
        }
      } catch (e) {
        console.error("Failed to parse cart", e);
        // Clear corrupted cart data
        localStorage.removeItem('cart');
      }
    }
  }, []);

  // Save cart to local storage whenever it changes
  useEffect(() => {
    localStorage.setItem('cart', JSON.stringify(items));
  }, [items]);

  const addToCart = (product: Product, colorVariant?: ColorVariant, sizeVariant?: SizeVariant) => {
    setItems((prev) => {
      // Find existing item with same product and variant combination
      const existing = prev.find((item) => 
        item.product?.id === product.id &&
        item.selectedColorVariant?.id === colorVariant?.id &&
        item.selectedSizeVariant?.id === sizeVariant?.id
      );
      
      if (existing) {
        toast.info("Updated quantity in cart");
        return prev.map((item) =>
          item.product?.id === product.id &&
          item.selectedColorVariant?.id === colorVariant?.id &&
          item.selectedSizeVariant?.id === sizeVariant?.id
            ? { ...item, quantity: item.quantity + 1 }
            : item
        );
      }
      
      toast.success("Added to cart");
      return [...prev, { 
        product, 
        quantity: 1,
        selectedColorVariant: colorVariant,
        selectedSizeVariant: sizeVariant
      }];
    });
  };

  const removeFromCart = (productId: number, colorVariantId?: number, sizeVariantId?: number) => {
    setItems((prev) => prev.filter((item) => 
      !(item.product?.id === productId &&
        item.selectedColorVariant?.id === colorVariantId &&
        item.selectedSizeVariant?.id === sizeVariantId)
    ));
    toast.success("Removed from cart");
  };

  const updateQuantity = (productId: number, quantity: number, colorVariantId?: number, sizeVariantId?: number) => {
    if (quantity < 1) {
      removeFromCart(productId, colorVariantId, sizeVariantId);
      return;
    }
    setItems((prev) =>
      prev.map((item) => 
        item.product?.id === productId &&
        item.selectedColorVariant?.id === colorVariantId &&
        item.selectedSizeVariant?.id === sizeVariantId
          ? { ...item, quantity }
          : item
      )
    );
  };

  const clearCart = () => {
    setItems([]);
    localStorage.removeItem('cart');
  };

  const cartTotal = items.reduce((total, item) => {
    // Safety check: ensure product and price exist
    if (item?.product && typeof item.product.price === 'number') {
      return total + item.product.price * item.quantity;
    }
    return total;
  }, 0);
  
  const cartCount = items.reduce((count, item) => {
    // Safety check: ensure item exists
    if (item && typeof item.quantity === 'number') {
      return count + item.quantity;
    }
    return count;
  }, 0);

  return (
    <CartContext.Provider
      value={{
        items,
        addToCart,
        removeFromCart,
        updateQuantity,
        clearCart,
        cartTotal,
        cartCount,
      }}
    >
      {children}
    </CartContext.Provider>
  );
}

export const useCart = () => {
  const context = useContext(CartContext);
  if (context === undefined) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
};
