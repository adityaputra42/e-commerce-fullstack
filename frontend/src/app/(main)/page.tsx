'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { toast } from 'sonner';
import api, { categoryService, type Category } from '@/services/api';
import type { Product } from '@/types/product';
import ProductCard from '@/components/product/ProductCard';
import { ArrowRight, ShieldCheck, Truck, RotateCcw, Headphones, Layers } from 'lucide-react';
import { motion } from 'framer-motion';

const FALLBACK_GRADIENTS = [
  'from-teal-500 to-emerald-600',
  'from-rose-500 to-orange-500',
  'from-indigo-500 to-purple-600',
  'from-amber-500 to-pink-500',
];

export default function Home() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [newsletterEmail, setNewsletterEmail] = useState('');

  useEffect(() => {
    fetchProducts();
    categoryService.getAll().then(setCategories).catch((err) => console.error('Error fetching categories:', err));
  }, []);

  const fetchProducts = async () => {
    try {
      const response = await api.get('/products?page=1&limit=8');
      const data = response.data?.data?.products || response.data?.data || [];
      setProducts(Array.isArray(data) ? data : []);
    } catch (error) {
      console.error('Error fetching products:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // No newsletter endpoint exists on the backend yet — this deliberately
  // does not pretend to succeed. A fake "Subscribed!" toast would be a UI
  // lie the same way a fabricated discount percentage would be.
  const handleNewsletterSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    toast.info("Newsletter signup isn't connected yet — coming soon.");
    setNewsletterEmail('');
  };

  return (
    <div className="flex flex-col min-h-screen overflow-x-hidden">
      {/* Hero Section — calm, boxy: flat off-white background, a plain
          bordered image box instead of rotated color blobs, no fabricated
          stat badge. */}
      <section className="py-16 md:py-24">
        <div className="container mx-auto px-6 grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <motion.div
            initial={{ opacity: 0, y: 16 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="space-y-6"
          >
            <div className="inline-flex items-center gap-2 px-3 py-1.5 bg-primary/10 rounded-md text-primary text-xs font-bold uppercase tracking-widest">
              Trending Now
            </div>

            <h1 className="text-5xl md:text-6xl font-black text-neutral-900 leading-[1.1] tracking-tight">
              Discover Products You&apos;ll Love
            </h1>

            <p className="max-w-md text-base text-neutral-500 font-medium leading-relaxed">
              Shop the latest trending products, curated for a modern lifestyle. Quality
              materials, honest pricing, no gimmicks.
            </p>

            <div className="flex flex-wrap gap-3 pt-2">
              <Link
                href="/shop"
                className="h-13 px-8 bg-primary text-white rounded-md text-sm font-bold hover:bg-primary-dark active:scale-[0.98] transition-all flex items-center gap-2"
              >
                Shop Now
                <ArrowRight className="w-4 h-4" />
              </Link>
              <Link
                href="/collections"
                className="h-13 px-8 bg-white border border-neutral-300 text-neutral-900 rounded-md text-sm font-bold hover:border-neutral-900 transition-all active:scale-[0.98] flex items-center"
              >
                Explore Collections
              </Link>
            </div>
          </motion.div>

          <motion.div
            initial={{ opacity: 0, scale: 0.96 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.5, delay: 0.1 }}
            className="relative aspect-square max-w-125 mx-auto w-full rounded-lg overflow-hidden border border-neutral-200"
          >
            <img
              src="https://images.unsplash.com/photo-1515886657613-9f3515b0c78f?auto=format&fit=crop&q=80"
              alt="Featured product model"
              className="w-full h-full object-cover"
            />
          </motion.div>
        </div>
      </section>

      {/* Features Bar — plain, light, no dark full-bleed section. */}
      <section className="border-y border-neutral-200 py-10">
        <div className="container mx-auto px-6 grid grid-cols-2 md:grid-cols-4 gap-8">
          {[
            { icon: Truck, title: 'Free Shipping', subtitle: 'On orders over $50' },
            { icon: ShieldCheck, title: 'Secure Payments', subtitle: '100% secure checkout' },
            { icon: RotateCcw, title: 'Easy Returns', subtitle: '30-day return policy' },
            { icon: Headphones, title: '24/7 Support', subtitle: 'Always here to help' },
          ].map((f) => (
            <div key={f.title} className="flex items-center gap-3">
              <f.icon className="w-6 h-6 text-primary shrink-0" />
              <div>
                <h4 className="text-sm font-bold text-neutral-900">{f.title}</h4>
                <p className="text-xs text-neutral-500">{f.subtitle}</p>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* Shop by Categories — real data from GET /categories, same source
          the Collections page uses, so this isn't a second hardcoded list
          that can drift out of sync with the real category set. */}
      {categories.length > 0 && (
        <section className="py-20 container mx-auto px-6">
          <div className="flex items-end justify-between gap-6 mb-8">
            <h2 className="text-2xl font-black text-neutral-900">Shop by Categories</h2>
            <Link href="/collections" className="flex items-center gap-2 text-sm font-bold text-neutral-600 hover:text-primary transition-colors">
              View All Categories
              <ArrowRight className="w-4 h-4" />
            </Link>
          </div>

          <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-4">
            {categories.slice(0, 6).map((category, i) => (
              <Link
                key={category.id}
                href={`/shop?category_id=${category.id}`}
                className="group relative aspect-square rounded-lg overflow-hidden border border-neutral-200"
              >
                {category.icon ? (
                  <img
                    src={category.icon}
                    alt={category.name}
                    className="absolute inset-0 w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                  />
                ) : (
                  <div className={`absolute inset-0 bg-linear-to-br ${FALLBACK_GRADIENTS[i % FALLBACK_GRADIENTS.length]} flex items-center justify-center`}>
                    <Layers className="w-8 h-8 text-white/30" />
                  </div>
                )}
                <div className="absolute inset-0 bg-linear-to-t from-black/60 to-transparent" />
                <span className="absolute bottom-3 left-3 text-sm font-bold text-white">{category.name}</span>
              </Link>
            ))}
          </div>
        </section>
      )}

      {/* Products Section */}
      <section className="py-20 container mx-auto px-6">
        <div className="flex flex-col md:flex-row items-end justify-between gap-6 mb-10">
          <div className="space-y-2">
            <div className="flex items-center gap-2 text-primary">
              <div className="w-6 h-1 bg-primary" />
              <span className="text-xs font-bold uppercase tracking-[0.2em]">Our Selection</span>
            </div>
            <h2 className="text-3xl font-black text-neutral-900">Featured Items</h2>
          </div>
          <Link href="/shop" className="flex items-center gap-2 text-sm font-bold text-neutral-900 hover:text-primary transition-colors group">
            View All Products
            <ArrowRight className="w-4 h-4 group-hover:translate-x-1 transition-transform" />
          </Link>
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {[...Array(4)].map((_, i) => (
              <div key={i} className="aspect-4/6 bg-neutral-100 rounded-lg animate-pulse" />
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-x-6 gap-y-10">
            {products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
            {products.length === 0 && (
              <div className="col-span-full py-20 text-center bg-neutral-50 rounded-lg border border-dashed border-neutral-200">
                <p className="text-neutral-400 font-bold uppercase tracking-widest">No products available at the moment.</p>
              </div>
            )}
          </div>
        )}
      </section>

      {/* Newsletter Section — no fabricated discount promise. */}
      <section className="pb-24 container mx-auto px-6">
        <div className="bg-neutral-900 rounded-lg p-10 md:p-16">
          <div className="max-w-xl space-y-6">
            <h2 className="text-3xl md:text-4xl font-black text-white leading-tight">
              Stay in the loop.
            </h2>
            <p className="text-neutral-400 font-medium">
              Get notified about new arrivals and product updates. No spam, unsubscribe anytime.
            </p>
            <form onSubmit={handleNewsletterSubmit} className="flex flex-col sm:flex-row gap-3">
              <input
                type="email"
                required
                value={newsletterEmail}
                onChange={(e) => setNewsletterEmail(e.target.value)}
                placeholder="Enter your email"
                className="grow h-13 px-5 bg-white/10 border border-white/20 rounded-md text-white placeholder:text-neutral-400 focus:bg-white/5 outline-none transition-all font-medium"
              />
              <button
                type="submit"
                className="h-13 px-8 bg-white text-neutral-900 rounded-md text-sm font-bold hover:bg-neutral-100 transition-colors active:scale-[0.98] whitespace-nowrap"
              >
                Subscribe
              </button>
            </form>
          </div>
        </div>
      </section>
    </div>
  );
}
