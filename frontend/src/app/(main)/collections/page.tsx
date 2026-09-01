'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { ArrowRight, Layers } from 'lucide-react';
import { motion } from 'framer-motion';
import { categoryService, type Category } from '@/services/api';

// A small, fixed palette so category cards don't all look identical when a
// category has no icon image set. Cycled by index, not random, so the page
// looks the same on every reload instead of shuffling on the user.
const FALLBACK_GRADIENTS = [
  'from-teal-500 to-emerald-600',
  'from-rose-500 to-orange-500',
  'from-indigo-500 to-purple-600',
  'from-amber-500 to-pink-500',
  'from-sky-500 to-blue-600',
  'from-slate-700 to-slate-900',
];

export default function CollectionsPage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(false);

  useEffect(() => {
    categoryService
      .getAll()
      .then(setCategories)
      .catch((err) => {
        console.error('Error fetching categories:', err);
        setError(true);
      })
      .finally(() => setIsLoading(false));
  }, []);

  return (
    <div className="min-h-screen pt-32 pb-24">
      <div className="container mx-auto px-6">
        <div className="mb-16 space-y-4">
          <div className="flex items-center gap-2 text-primary">
            <div className="w-8 h-1 bg-primary" />
            <span className="text-xs font-black uppercase tracking-[0.3em]">Browse By</span>
          </div>
          <h1 className="text-5xl md:text-6xl font-black text-slate-900 dark:text-white tracking-tighter italic uppercase">
            Collections
          </h1>
          <p className="max-w-md text-slate-500 font-medium">
            Every category we carry, in one place. Pick one to jump straight into a filtered shop view.
          </p>
        </div>

        {isLoading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
            {[...Array(6)].map((_, i) => (
              <div key={i} className="aspect-video bg-slate-100 dark:bg-slate-800 rounded-[2.5rem] animate-pulse" />
            ))}
          </div>
        ) : error ? (
          <div className="py-24 text-center bg-slate-50 dark:bg-slate-900 rounded-[3rem] border-2 border-dashed border-slate-200 dark:border-slate-800">
            <p className="text-slate-400 font-bold uppercase tracking-widest italic">
              Couldn't load collections right now.
            </p>
          </div>
        ) : categories.length === 0 ? (
          <div className="py-24 text-center bg-slate-50 dark:bg-slate-900 rounded-[3rem] border-2 border-dashed border-slate-200 dark:border-slate-800">
            <p className="text-slate-400 font-bold uppercase tracking-widest italic">
              No collections have been added yet.
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
            {categories.map((category, i) => (
              <motion.div
                key={category.id}
                initial={{ opacity: 0, y: 20 }}
                whileInView={{ opacity: 1, y: 0 }}
                viewport={{ once: true }}
                transition={{ delay: i * 0.05 }}
              >
                <Link
                  href={`/shop?category_id=${category.id}`}
                  className="group relative flex flex-col justify-end aspect-video rounded-[2.5rem] overflow-hidden shadow-lg hover:shadow-2xl transition-shadow"
                >
                  {category.icon ? (
                    <img
                      src={category.icon}
                      alt={category.name}
                      className="absolute inset-0 w-full h-full object-cover group-hover:scale-110 transition-transform duration-700"
                    />
                  ) : (
                    <div
                      className={`absolute inset-0 bg-linear-to-br ${
                        FALLBACK_GRADIENTS[i % FALLBACK_GRADIENTS.length]
                      } flex items-center justify-center`}
                    >
                      <Layers className="w-16 h-16 text-white/20" />
                    </div>
                  )}

                  <div className="absolute inset-0 bg-linear-to-t from-black/70 via-black/10 to-transparent" />

                  <div className="relative z-10 p-8 flex items-end justify-between">
                    <h3 className="text-2xl font-black text-white tracking-tight italic uppercase">
                      {category.name}
                    </h3>
                    <div className="w-10 h-10 rounded-full bg-white/20 backdrop-blur-sm flex items-center justify-center group-hover:bg-white group-hover:text-slate-900 text-white transition-colors shrink-0">
                      <ArrowRight className="w-4 h-4 group-hover:translate-x-0.5 transition-transform" />
                    </div>
                  </div>
                </Link>
              </motion.div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
