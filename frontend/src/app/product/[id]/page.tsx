'use client';

import { useState, useEffect, use } from 'react';
import api from '@/lib/api';
import type { Product, ColorVariant, SizeVariant } from '@/types/product';
import { ShoppingBag, Star, Share2, Heart, ShieldCheck, ChevronRight, Truck, RotateCcw } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { useRouter } from 'next/navigation';

export default function ProductDetail({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const router = useRouter();
  const [product, setProduct] = useState<Product | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [activeImage, setActiveImage] = useState('');
  const [selectedVariant, setSelectedVariant] = useState<ColorVariant | null>(null);
  const [selectedSize, setSelectedSize] = useState<string | null>(null);
  const [quantity, setQuantity] = useState(1);

  useEffect(() => {
    fetchProductDetail();
  }, [id]);

  const fetchProductDetail = async () => {
    try {
      const response = await api.get(`/products/${id}`);
      const data = response.data?.data;
      setProduct(data);
      if (data) {
        setActiveImage(data.images);
        if (data.color_varian && data.color_varian.length > 0) {
           setSelectedVariant(data.color_varian[0]);
           const firstVariant = data.color_varian[0];
           if (firstVariant.size_varian && firstVariant.size_varian.length > 0) {
              setSelectedSize(firstVariant.size_varian[0].size);
           }
        }
      }
    } catch (error) {
      console.error('Error fetching product detail:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleOrder = () => {
    // Check for token in localStorage (simplest way for now)
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    if (!token) {
       // Save intended action? or just redirect to login
       router.push('/login?redirect=/product/' + id);
       return;
    }
    // Proceed to checkout logic (to be implemented)
    alert('Proceeding to order... (Login authenticated)');
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
         <div className="w-16 h-16 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin" />
      </div>
    );
  }

  if (!product) {
    return (
      <div className="flex flex-col items-center justify-center min-h-screen gap-4">
         <XCircle className="w-20 h-20 text-rose-500" />
         <h2 className="text-2xl font-black text-slate-900">Product not found.</h2>
         <button onClick={() => router.push('/')} className="text-primary font-bold uppercase tracking-widest text-sm underline underline-offset-8 transition-all hover:scale-105">Back to Home</button>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-6 py-12 md:py-24">
      {/* Breadcrumbs */}
      <div className="flex items-center gap-2 text-[10px] font-black uppercase tracking-widest text-slate-400 mb-12">
         <button onClick={() => router.push('/')} className="hover:text-primary">Home</button>
         <ChevronRight className="w-3 h-3" />
         <span className="text-slate-900">{product.category?.name}</span>
         <ChevronRight className="w-3 h-3" />
         <span className="text-primary">{product.name}</span>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-16">
        {/* Left: Images */}
        <div className="space-y-6">
           <div className="aspect-4/5 bg-slate-50 dark:bg-slate-900 rounded-[3rem] overflow-hidden group shadow-2xl relative">
              <AnimatePresence mode="wait">
                <motion.img 
                  key={activeImage}
                  initial={{ opacity: 0, scale: 1.1 }}
                  animate={{ opacity: 1, scale: 1 }}
                  exit={{ opacity: 0 }}
                  transition={{ duration: 0.5 }}
                  src={activeImage || 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&q=80'} 
                  alt={product.name}
                  className="w-full h-full object-cover"
                />
              </AnimatePresence>
              
              <div className="absolute top-6 right-6">
                 <button className="w-12 h-12 bg-white rounded-full flex items-center justify-center text-slate-400 hover:text-rose-500 shadow-xl transition-all active:scale-90">
                    <Heart className="w-6 h-6" />
                 </button>
              </div>
           </div>

           <div className="flex gap-4 p-2 overflow-x-auto custom-scrollbar">
              <button 
                onClick={() => setActiveImage(product.images)}
                className={`w-24 h-24 rounded-2xl overflow-hidden border-2 shrink-0 transition-all ${
                  activeImage === product.images ? 'border-indigo-600 scale-105 shadow-lg' : 'border-transparent opacity-60 hover:opacity-100'
                }`}
              >
                 <img src={product.images} className="w-full h-full object-cover" />
              </button>
              {product.color_varian?.map((v: ColorVariant) => (
                 v.images && (
                  <button 
                    key={v.id}
                    onClick={() => {
                        setActiveImage(v.images!);
                        setSelectedVariant(v);
                    }}
                    className={`w-24 h-24 rounded-2xl overflow-hidden border-2 shrink-0 transition-all ${
                      activeImage === v.images ? 'border-indigo-600 scale-105 shadow-lg' : 'border-transparent opacity-60 hover:opacity-100'
                    }`}
                  >
                     <img src={v.images} className="w-full h-full object-cover" />
                  </button>
                 )
              ))}
           </div>
        </div>

        {/* Right: Info */}
        <div className="flex flex-col gap-10">
           <div className="space-y-4">
              <div className="flex items-center gap-2">
                 <div className="flex text-amber-500">
                    {[...Array(5)].map((_, i) => <Star key={i} className="w-4 h-4 fill-current" />)}
                 </div>
                 <span className="text-xs font-black text-slate-400 uppercase tracking-widest">(128 Reviews)</span>
              </div>
              <h1 className="text-4xl md:text-6xl font-black text-slate-900 dark:text-white leading-tight italic tracking-tighter uppercase">{product.name}</h1>
              <p className="text-3xl font-black text-primary font-display italic">
                 {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(product.price)}
              </p>
           </div>

           <p className="text-base text-slate-500 font-medium leading-[1.8] max-w-xl">
              {product.description || 'Premium craftsmanship meets modern design in this exclusive piece from our latest collection. Every detail is meticulously thought out to provide unparalleled quality and style.'}
           </p>

           {/* Variants */}
           <div className="space-y-8">
              {product.color_varian && product.color_varian.length > 0 && (
                <div className="space-y-4">
                   <h4 className="text-[10px] font-black text-slate-400 uppercase tracking-[0.3em]">Select Color: <span className="text-slate-900">{selectedVariant?.name}</span></h4>
                   <div className="flex gap-3">
                      {product.color_varian.map((v: ColorVariant) => (
                         <button 
                            key={v.id}
                            onClick={() => {
                                setSelectedVariant(v);
                                if (v.images) setActiveImage(v.images);
                                if (v.size_varian && v.size_varian.length > 0 && !v.size_varian.some((s: SizeVariant) => s.size === selectedSize)) {
                                     setSelectedSize(v.size_varian[0].size);
                                }
                            }}
                            className={`w-10 h-10 rounded-full border-2 transition-all p-1 ${
                                selectedVariant?.id === v.id ? 'border-indigo-600 scale-110 shadow-lg' : 'border-slate-100'
                            }`}
                         >
                            <div className="w-full h-full rounded-full" style={{ backgroundColor: v.color }} />
                         </button>
                      ))}
                   </div>
                </div>
              )}

              {selectedVariant?.size_varian && selectedVariant.size_varian.length > 0 && (
                <div className="space-y-4">
                   <h4 className="text-[10px] font-black text-slate-400 uppercase tracking-[0.3em]">Select Size</h4>
                   <div className="flex flex-wrap gap-2">
                      {selectedVariant.size_varian.map((s: SizeVariant) => (
                         <button 
                            key={s.id}
                            disabled={s.stock === 0}
                            onClick={() => setSelectedSize(s.size)}
                            className={`min-w-[50px] h-12 flex items-center justify-center px-4 rounded-xl text-xs font-black uppercase tracking-tighter transition-all ${
                                selectedSize === s.size 
                                ? 'bg-primary text-white shadow-xl shadow-indigo-600/20' 
                                : s.stock === 0 ? 'bg-slate-50 text-slate-300 cursor-not-allowed border-dashed border-2 border-slate-200' : 'bg-slate-50 text-slate-600 hover:bg-slate-100'
                            }`}
                         >
                            {s.size}
                         </button>
                      ))}
                   </div>
                   <p className="text-[10px] text-slate-400 font-bold italic">
                      {selectedVariant.size_varian.find((s: SizeVariant) => s.size === selectedSize)?.stock || 0} items available in stock
                   </p>
                </div>
              )}
           </div>

           {/* Quantity & Actions */}
           <div className="flex flex-col sm:flex-row items-center gap-4 pt-4 border-t border-slate-100 mt-4">
              <div className="flex items-center bg-slate-50 rounded-2xl h-16 px-4">
                 <button 
                    onClick={() => setQuantity(Math.max(1, quantity - 1))}
                    className="w-10 h-10 flex items-center justify-center text-slate-500 hover:text-slate-900 transition-colors"
                 >-</button>
                 <span className="w-12 text-center text-sm font-black italic">{quantity}</span>
                 <button 
                    onClick={() => setQuantity(quantity + 1)}
                    className="w-10 h-10 flex items-center justify-center text-slate-500 hover:text-slate-900 transition-colors"
                 >+</button>
              </div>

              <button 
                onClick={handleOrder}
                className="grow h-16 w-full px-10 bg-primary text-white rounded-2xl flex items-center justify-center gap-3 text-sm font-black uppercase tracking-widest hover:bg-indigo-700 transition-all shadow-2xl shadow-indigo-600/30 active:scale-95 group"
              >
                 <ShoppingBag className="w-5 h-5 group-hover:-rotate-12 transition-transform" />
                 Buy & Deliver Now
              </button>
              
              <button className="h-16 w-16 bg-white border-2 border-slate-100 rounded-2xl flex items-center justify-center text-slate-400 hover:text-primary transition-all hover:border-indigo-600 shrink-0">
                 <Share2 className="w-5 h-5" />
              </button>
           </div>

           {/* Trust Badges */}
           <div className="grid grid-cols-2 gap-4 mt-8">
              <div className="flex items-center gap-3 p-4 rounded-4xl bg-indigo-50/50 border border-indigo-100/50">
                 <Truck className="w-6 h-6 text-primary" />
                 <div>
                    <p className="text-[10px] font-black text-slate-900 uppercase">Premium Shipping</p>
                    <p className="text-[8px] font-bold text-slate-400 uppercase tracking-widest">Across Indonesia</p>
                 </div>
              </div>
              <div className="flex items-center gap-3 p-4 rounded-4xl bg-rose-50/50 border border-rose-100/50">
                 <RotateCcw className="w-6 h-6 text-rose-600" />
                 <div>
                    <p className="text-[10px] font-black text-slate-900 uppercase">30-Day Returns</p>
                    <p className="text-[8px] font-bold text-slate-400 uppercase tracking-widest">Genuine Guarantee</p>
                 </div>
              </div>
           </div>
        </div>
      </div>
    </div>
  );
}

function XCircle(props: React.SVGProps<SVGSVGElement>) {
  return (
    <svg {...props} xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="lucide lucide-circle-x">
      <circle cx="12" cy="12" r="10"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/>
    </svg>
  );
}
