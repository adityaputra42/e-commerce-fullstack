'use client';

import { useState, useEffect } from 'react';
import api from '@/lib/api';
import type { Product } from '@/types/product';
import ProductCard from '@/components/product/ProductCard';
import { ArrowRight, Sparkles, TrendingUp, ShieldCheck, Zap } from 'lucide-react';
import { motion } from 'framer-motion';

export default function Home() {
  const [products, setProducts] = useState<Product[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    fetchProducts();
  }, []);

  const fetchProducts = async () => {
    try {
      const response = await api.get('/products?page=1&limit=10');
      const data = response.data?.data?.products || response.data?.data || [];
      setProducts(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Error fetching products:', error);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex flex-col min-h-screen overflow-x-hidden">
      {/* Hero Section */}
      <section className="relative min-h-[90vh] flex items-center pt-20 pb-20 overflow-hidden">
        {/* Background Gradients */}
        <div className="absolute top-0 right-0 w-[50%] h-full bg-linear-to-l from-teal-50/50 to-transparent -z-10 blur-3xl" />
        <div className="absolute bottom-0 left-0 w-[40%] h-[80%] bg-linear-to-tr from-rose-50/30 to-transparent -z-10 blur-3xl" />
        
        <div className="container mx-auto px-6 grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <motion.div 
            initial={{ opacity: 0, x: -50 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.8 }}
            className="space-y-8"
          >
            <div className="inline-flex items-center gap-2 px-4 py-2 bg-teal-50 dark:bg-teal-900/30 rounded-full text-primary dark:text-teal-400 text-xs font-black uppercase tracking-widest animate-pulse">
              <Sparkles className="w-4 h-4" />
              New Collection 2026 available now
            </div>
            
            <h1 className="text-6xl md:text-7xl lg:text-8xl font-black text-slate-900 dark:text-white leading-[1.1] tracking-tighter italic">
              ELEVATE <br />
              <span className="text-primary outline-text">EVERY</span> <br />
              MOMENT.
            </h1>
            
            <p className="max-w-md text-lg text-slate-500 font-medium leading-relaxed">
              Experience the fusion of high-performance materials and timeless design. Our latest collection is crafted for those who demand excellence.
            </p>
            
            <div className="flex flex-wrap gap-4 pt-4">
              <button className="h-16 px-10 bg-primary text-white rounded-2xl text-sm font-black uppercase tracking-widest hover:bg-secondary transition-all shadow-2xl shadow-teal-600/30 active:scale-95 group flex items-center gap-3">
                Shop Collection
                <ArrowRight className="w-5 h-5 group-hover:translate-x-2 transition-transform" />
              </button>
              <button className="h-16 px-10 bg-white dark:bg-slate-800 border-2 border-slate-100 dark:border-slate-700 text-slate-900 dark:text-white rounded-2xl text-sm font-black uppercase tracking-widest hover:border-teal-600 transition-all active:scale-95">
                Lookbook
              </button>
            </div>
          </motion.div>

          {/* Hero Image / Design */}
          <motion.div 
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 1 }}
            className="relative"
          >
            <div className="relative aspect-square max-w-125 mx-auto">
               <div className="absolute inset-0 bg-primary rounded-[4rem] rotate-6 -z-10 opacity-10 animate-pulse" />
               <div className="absolute inset-0 bg-rose-500 rounded-[4rem] -rotate-3 -z-10 opacity-10 animate-pulse delay-1000" />
               <img 
                 src="https://images.unsplash.com/photo-1515886657613-9f3515b0c78f?auto=format&fit=crop&q=80" 
                 alt="Hero Model"
                 className="w-full h-full object-cover rounded-[3rem] shadow-2xl shadow-teal-500/20"
               />
               
               {/* Floating Badges */}
               <motion.div 
                 animate={{ y: [0, -10, 0] }}
                 transition={{ duration: 4, repeat: Infinity }}
                 className="absolute -top-10 -right-10 p-6 bg-white dark:bg-slate-800 rounded-3xl shadow-2xl flex items-center gap-4"
                >
                 <div className="w-12 h-12 bg-emerald-50 text-emerald-600 rounded-2xl flex items-center justify-center">
                    <TrendingUp className="w-6 h-6" />
                 </div>
                 <div>
                    <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Trending</p>
                    <p className="text-xl font-black text-slate-900 dark:text-white">+840%</p>
                 </div>
               </motion.div>
            </div>
          </motion.div>
        </div>
      </section>

      {/* Features Bar */}
      <section className="bg-slate-900 text-white py-12">
        <div className="container mx-auto px-6 grid grid-cols-2 md:grid-cols-4 gap-8">
            <div className="flex items-center gap-4">
               <ShieldCheck className="w-8 h-8 text-teal-400" />
               <div>
                  <h4 className="text-sm font-black uppercase tracking-tighter">Premium Quality</h4>
                  <p className="text-[10px] text-slate-400 font-bold uppercase tracking-widest">Certified Materials</p>
               </div>
            </div>
            <div className="flex items-center gap-4">
               <Zap className="w-8 h-8 text-teal-400" />
               <div>
                  <h4 className="text-sm font-black uppercase tracking-tighter">Rapid Delivery</h4>
                  <p className="text-[10px] text-slate-400 font-bold uppercase tracking-widest">Express Worldwide</p>
               </div>
            </div>
            <div className="flex items-center gap-4">
               <Sparkles className="w-8 h-8 text-teal-400" />
               <div>
                  <h4 className="text-sm font-black uppercase tracking-tighter">Exclusive Drops</h4>
                  <p className="text-[10px] text-slate-400 font-bold uppercase tracking-widest">Limited Editions</p>
               </div>
            </div>
            <div className="flex items-center gap-4">
               <ShieldCheck className="w-8 h-8 text-teal-400" />
               <div>
                  <h4 className="text-sm font-black uppercase tracking-tighter">Secure Checkout</h4>
                  <p className="text-[10px] text-slate-400 font-bold uppercase tracking-widest">100% Encrypted</p>
               </div>
            </div>
        </div>
      </section>

      {/* Products Section */}
      <section className="py-32 container mx-auto px-6">
        <div className="flex flex-col md:flex-row items-end justify-between gap-6 mb-16">
            <div className="space-y-4">
               <div className="flex items-center gap-2 text-primary">
                  <div className="w-8 h-1 bg-primary" />
                  <span className="text-xs font-black uppercase tracking-[0.3em]">Our Selection</span>
               </div>
               <h2 className="text-5xl font-black text-slate-900 dark:text-white tracking-tighter italic uppercase underline decoration-teal-600/30 underline-offset-8">Featured Items</h2>
            </div>
            <button className="flex items-center gap-3 text-sm font-black text-slate-900 dark:text-white uppercase tracking-widest hover:text-primary transition-colors group">
               View All Products
               <ArrowRight className="w-4 h-4 group-hover:translate-x-2 transition-transform" />
            </button>
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="aspect-4/6 bg-slate-100 dark:bg-slate-800 rounded-3xl animate-pulse" />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-x-8 gap-y-12">
            {products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
            {products.length === 0 && (
                <div className="col-span-full py-20 text-center bg-slate-50 dark:bg-slate-900 rounded-[3rem] border-2 border-dashed border-slate-200 dark:border-slate-800">
                   <p className="text-slate-400 font-bold uppercase tracking-widest italic">No products available at the moment.</p>
                </div>
            )}
          </div>
        )}
      </section>

      {/* Newsletter / CTA Section */}
      <section className="pb-32 container mx-auto px-6">
          <div className="bg-primary rounded-[4rem] p-12 md:p-24 relative overflow-hidden group">
              <div className="absolute top-0 right-0 w-96 h-96 bg-teal rounded-full translate-x-32 -translate-y-32 group-hover:scale-110 transition-transform duration-1000" />
              <div className="absolute bottom-0 left-0 w-64 h-64 bg-secondary rounded-full -translate-x-16 translate-y-16" />
              
              <div className="relative z-10 max-w-2xl space-y-8">
                  <h2 className="text-4xl md:text-6xl font-black text-white leading-tight italic tracking-tight uppercase underline decoration-white/30 underline-offset-8">
                    GET 15% OFF <br />YOUR FIRST ORDER.
                  </h2>
                  <p className="text-lg text-teal-100/80 font-medium">
                    Be the first to know about new arrivals, fashion shows, and exclusive sales events.
                  </p>
                  <form className="flex flex-col sm:flex-row gap-4">
                     <input 
                      type="email" 
                      placeholder="Enter your email" 
                      className="grow h-16 px-8 bg-white/20 border-2 border-white/20 rounded-2xl text-white placeholder:text-white/50 focus:bg-white/5 focus:text-white outline-none transition-all font-bold"
                     />
                     <button className="h-16 px-12 bg-white text-slate-900 rounded-2xl text-sm font-black uppercase tracking-widest hover:bg-slate-50 transition-colors shadow-2xl active:scale-95 whitespace-nowrap">
                        Join StoreX
                     </button>
                  </form>
              </div>
          </div>
      </section>
    </div>
  );
}
