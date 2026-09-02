'use client';

import { useState, useEffect, useCallback, Suspense } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { Search, SlidersHorizontal, X, ChevronRight } from 'lucide-react';
import api, { categoryService, type Category } from '@/services/api';
import type { Product } from '@/types/product';
import ProductCard from '@/components/product/ProductCard';

const SORT_OPTIONS = [
  { value: 'created_at', label: 'Newest' },
  { value: 'price', label: 'Price' },
  { value: 'name', label: 'Name (A-Z)' },
];

const LIMIT = 12;

function ShopContent() {
  const searchParams = useSearchParams();
  const router = useRouter();

  const initialCategoryId = searchParams.get('category_id') || '';

  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingMore, setIsLoadingMore] = useState(false);
  const [hasMore, setHasMore] = useState(false);
  const [page, setPage] = useState(1);

  const [search, setSearch] = useState('');
  const [searchInput, setSearchInput] = useState('');
  const [categoryId, setCategoryId] = useState(initialCategoryId);
  const [sortBy, setSortBy] = useState('created_at');
  const [isFilterOpen, setIsFilterOpen] = useState(false);

  // Categories are fetched once — they don't change while filtering products.
  useEffect(() => {
    categoryService
      .getAll()
      .then(setCategories)
      .catch((err) => console.error('Error fetching categories:', err));
  }, []);

  const fetchProducts = useCallback(
    async (targetPage: number, append: boolean) => {
      append ? setIsLoadingMore(true) : setIsLoading(true);
      try {
        const params: Record<string, string | number> = {
          page: targetPage,
          limit: LIMIT,
          sort_by: sortBy,
        };
        if (search) params.search = search;
        if (categoryId) params.category_id = categoryId;

        const response = await api.get('/products', { params });
        const data = response.data?.data?.products || response.data?.data || [];
        const list: Product[] = Array.isArray(data) ? data : [];

        setProducts((prev) => (append ? [...prev, ...list] : list));
        // Backend doesn't return a total count, so "has more" is inferred:
        // a full page came back, there's probably another one.
        setHasMore(list.length === LIMIT);
        setPage(targetPage);
      } catch (error) {
        console.error('Error fetching products:', error);
        if (!append) setProducts([]);
      } finally {
        append ? setIsLoadingMore(false) : setIsLoading(false);
      }
    },
    [search, categoryId, sortBy]
  );

  // Re-fetch page 1 whenever a filter changes.
  useEffect(() => {
    fetchProducts(1, false);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [search, categoryId, sortBy]);

  // Keep the URL in sync with the category filter so /collections links
  // and browser back/forward behave the way a shop page should.
  useEffect(() => {
    const params = new URLSearchParams();
    if (categoryId) params.set('category_id', categoryId);
    const query = params.toString();
    router.replace(query ? `/shop?${query}` : '/shop', { scroll: false });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [categoryId]);

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSearch(searchInput.trim());
  };

  const activeCategoryName = categories.find((c) => String(c.id) === categoryId)?.name;
  const hasActiveFilters = Boolean(search || categoryId);

  const clearFilters = () => {
    setSearch('');
    setSearchInput('');
    setCategoryId('');
    setSortBy('created_at');
  };

  return (
    <div className="min-h-screen pt-12 pb-24">
      <div className="container mx-auto px-6">
        {/* Header */}
        <div className="mb-12 space-y-4">
          <div className="flex items-center gap-2 text-primary">
            <div className="w-8 h-1 bg-primary" />
            <span className="text-xs font-black uppercase tracking-[0.3em]">Catalog</span>
          </div>
          <h1 className="text-5xl md:text-6xl font-black text-slate-900 dark:text-white tracking-tighter italic uppercase">
            {activeCategoryName || 'Shop All'}
          </h1>
        </div>

        {/* Controls */}
        <div className="flex flex-col md:flex-row gap-4 mb-10">
          <form onSubmit={handleSearchSubmit} className="relative flex-1">
            <Search className="absolute left-5 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" />
            <input
              type="text"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
              placeholder="Search products..."
              className="w-full h-14 pl-12 pr-4 bg-slate-50 dark:bg-slate-900 border-2 border-transparent focus:border-teal-500 rounded-md outline-none font-medium text-sm transition-colors"
            />
          </form>

          <button
            type="button"
            onClick={() => setIsFilterOpen((v) => !v)}
            className="md:hidden flex items-center justify-center gap-2 h-14 px-6 bg-slate-900 text-white rounded-md text-sm font-bold"
          >
            <SlidersHorizontal className="w-4 h-4" />
            Filters
          </button>

          <div className={`${isFilterOpen ? 'flex' : 'hidden'} md:flex flex-col sm:flex-row gap-4`}>
            <select
              value={categoryId}
              onChange={(e) => setCategoryId(e.target.value)}
              className="h-14 px-5 bg-slate-50 dark:bg-slate-900 border-2 border-transparent focus:border-teal-500 rounded-md outline-none font-bold text-sm text-slate-700 dark:text-slate-200 min-w-40"
            >
              <option value="">All Categories</option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>

            <select
              value={sortBy}
              onChange={(e) => setSortBy(e.target.value)}
              className="h-14 px-5 bg-slate-50 dark:bg-slate-900 border-2 border-transparent focus:border-teal-500 rounded-md outline-none font-bold text-sm text-slate-700 dark:text-slate-200 min-w-36"
            >
              {SORT_OPTIONS.map((opt) => (
                <option key={opt.value} value={opt.value}>
                  Sort: {opt.label}
                </option>
              ))}
            </select>
          </div>
        </div>

        {hasActiveFilters && (
          <div className="flex items-center gap-3 mb-8">
            <span className="text-xs font-bold text-slate-400 uppercase tracking-widest">Active filters</span>
            {search && (
              <span className="inline-flex items-center gap-2 px-3 py-1.5 bg-primary/10 text-primary rounded-md text-xs font-bold">
                "{search}"
              </span>
            )}
            {activeCategoryName && (
              <span className="inline-flex items-center gap-2 px-3 py-1.5 bg-primary/10 text-primary rounded-md text-xs font-bold">
                {activeCategoryName}
              </span>
            )}
            <button
              onClick={clearFilters}
              className="flex items-center gap-1 text-xs font-bold text-slate-400 hover:text-rose-500 transition-colors"
            >
              <X className="w-3.5 h-3.5" />
              Clear
            </button>
          </div>
        )}

        {/* Grid */}
        {isLoading ? (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8">
            {[...Array(8)].map((_, i) => (
              <div key={i} className="aspect-4/6 bg-slate-100 dark:bg-slate-800 rounded-3xl animate-pulse" />
            ))}
          </div>
        ) : products.length === 0 ? (
          <div className="py-24 text-center bg-slate-50 dark:bg-slate-900 rounded-lg border-2 border-dashed border-slate-200 dark:border-slate-800">
            <p className="text-slate-400 font-bold uppercase tracking-widest italic">
              No products match these filters.
            </p>
          </div>
        ) : (
          <>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-x-8 gap-y-12">
              {products.map((product) => (
                <ProductCard key={product.id} product={product} />
              ))}
            </div>

            {hasMore && (
              <div className="flex justify-center mt-16">
                <button
                  onClick={() => fetchProducts(page + 1, true)}
                  disabled={isLoadingMore}
                  className="flex items-center gap-3 h-14 px-10 bg-slate-900 text-white rounded-md text-sm font-black uppercase tracking-widest hover:bg-slate-800 transition-all active:scale-95 disabled:opacity-50"
                >
                  {isLoadingMore ? 'Loading...' : 'Load More'}
                  <ChevronRight className="w-4 h-4" />
                </button>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}

export default function ShopPage() {
  // useSearchParams needs a Suspense boundary in the app router.
  return (
    <Suspense fallback={null}>
      <ShopContent />
    </Suspense>
  );
}
