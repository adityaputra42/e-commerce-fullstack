"use client";

import { useCart } from '@/context/CartContext';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Trash2, Plus, Minus, ArrowRight, ShoppingBag } from 'lucide-react';
import Link from 'next/link';
import { useAuth } from '@/context/AuthContext';
import { useRouter } from 'next/navigation';

export default function CartPage() {
  const { items, removeFromCart, updateQuantity, cartTotal, clearCart } = useCart();
  const { isAuthenticated } = useAuth();
  const router = useRouter();

  const handleCheckout = () => {
    if (!isAuthenticated) {
        router.push('/login?redirect=/checkout');
    } else {
        router.push('/checkout');
    }
  };

  if (items.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] space-y-4">
        <div className="w-24 h-24 bg-slate-100 dark:bg-slate-800 rounded-full flex items-center justify-center text-slate-400 mb-4">
            <ShoppingBag className="w-10 h-10" />
        </div>
        <h1 className="text-2xl font-bold">Your cart is empty</h1>
        <p className="text-slate-500">Looks like you haven't added anything yet.</p>
        <Link href="/">
           <Button className="mt-4 gap-2">
             Start Shopping <ArrowRight className="w-4 h-4" />
           </Button>
        </Link>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-6 py-12">
      <h1 className="text-3xl font-black mb-8 italic">SHOPPING CART</h1>
      
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-12">
        {/* Cart Items */}
        <div className="lg:col-span-2 space-y-6">
          {items.map((item) => (
            <Card key={item.id} className="premium-card overflow-hidden">
              <CardContent className="p-0 flex flex-col sm:flex-row gap-6">
                <div className="w-full sm:w-32 aspect-square bg-slate-100 shrink-0">
                    <img src={item.images || "https://images.unsplash.com/photo-1523275335684-37898b6baf30"} alt={item.name} className="w-full h-full object-cover" />
                </div>
                
                <div className="flex-1 p-6 pl-0 flex flex-col justify-between">
                   <div className="flex justify-between items-start">
                      <div>
                        <h3 className="font-bold text-lg mb-1">{item.name}</h3>
                        <p className="text-sm text-slate-500">{item.category?.name || "General"}</p>
                      </div>
                      <p className="font-bold text-primary">
                         {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(item.price)}
                      </p>
                   </div>

                   <div className="flex justify-between items-end mt-4">
                      <div className="flex items-center gap-3 bg-slate-50 dark:bg-slate-900 rounded-lg p-1">
                          <button 
                            onClick={() => updateQuantity(item.id, item.quantity - 1)}
                            className="w-8 h-8 flex items-center justify-center rounded-md hover:bg-white dark:hover:bg-slate-800 shadow-sm transition-colors"
                          >
                            <Minus className="w-4 h-4" />
                          </button>
                          <span className="w-8 text-center font-bold text-sm">{item.quantity}</span>
                          <button 
                            onClick={() => updateQuantity(item.id, item.quantity + 1)}
                            className="w-8 h-8 flex items-center justify-center rounded-md hover:bg-white dark:hover:bg-slate-800 shadow-sm transition-colors"
                          >
                            <Plus className="w-4 h-4" />
                          </button>
                      </div>
                      
                      <button 
                        onClick={() => removeFromCart(item.id)}
                        className="text-slate-400 hover:text-red-500 transition-colors p-2"
                      >
                         <Trash2 className="w-5 h-5" />
                      </button>
                   </div>
                </div>
              </CardContent>
            </Card>
          ))}
          
          <Button variant="ghost" onClick={clearCart} className="text-red-500 hover:text-red-600 hover:bg-red-50">
            Clear Cart
          </Button>
        </div>

        {/* Order Summary */}
        <div className="lg:col-span-1">
           <Card className="premium-card sticky top-24">
              <CardContent className="p-8 space-y-6">
                 <h2 className="text-xl font-bold">Order Summary</h2>
                 
                 <div className="space-y-4 text-sm">
                    <div className="flex justify-between text-slate-500">
                       <span>Subtotal</span>
                       <span>{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(cartTotal)}</span>
                    </div>
                    <div className="flex justify-between text-slate-500">
                       <span>Shipping</span>
                       <span>Calculated at checkout</span>
                    </div>
                 </div>
                 
                 <div className="border-t pt-4 flex justify-between font-black text-lg">
                    <span>Total</span>
                    <span className="text-primary">{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(cartTotal)}</span>
                 </div>

                 <Button className="w-full h-14 rounded-xl font-bold text-lg shadow-xl shadow-primary/20" onClick={handleCheckout}>
                    Checkout <ArrowRight className="w-5 h-5 ml-2" />
                 </Button>
              </CardContent>
           </Card>
        </div>
      </div>
    </div>
  );
}
