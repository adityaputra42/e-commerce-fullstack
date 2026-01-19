'use client';

import { Product } from '@/types/product';
import Link from 'next/link';
import { Button } from '@/components/common/Button';
import { Eye, ShoppingBag, Heart } from 'lucide-react';
import { motion } from 'framer-motion';
import { useCart } from '@/context/CartContext';

interface ProductCardProps {
  product: Product;
}

const ProductCard = ({ product }: ProductCardProps) => {
  const { addToCart } = useCart();

  return (
    <motion.div 
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      className="premium-card bg-white dark:bg-slate-900 rounded-3xl overflow-hidden group flex flex-col h-full"
    >
      <div className="relative aspect-4/5 overflow-hidden">
        {/* Badge */}
        <div className="absolute top-4 left-4 z-10">
          <span className="px-3 py-1 bg-primary/90 text-white text-[10px] font-black uppercase tracking-widest rounded-full backdrop-blur-sm shadow-lg shadow-teal-600/20">
            {product.category?.name || 'NEW ARRIVAL'}
          </span>
        </div>

        {/* Quick Actions */}
        <div className="absolute top-4 right-4 z-10 flex flex-col gap-2 translate-x-12 group-hover:translate-x-0 transition-transform duration-500">
           <button className="w-10 h-10 bg-white rounded-full flex items-center justify-center text-slate-400 hover:text-rose-500 shadow-xl transition-colors">
             <Heart className="w-5 h-5" />
           </button>
           <Link href={`/product/${product.id}`} className="w-10 h-10 bg-white rounded-full flex items-center justify-center text-slate-400 hover:text-primary shadow-xl transition-colors">
             <Eye className="w-5 h-5" />
           </Link>
        </div>

        {/* Image */}
        <Link href={`/product/${product.id}`}>
          <img 
            src={product.images || 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&q=80'} 
            alt={product.name}
            className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-700"
          />
        </Link>

        {/* Add to Cart Overlay */}
        <div className="absolute inset-x-0 bottom-0 p-4 translate-y-full group-hover:translate-y-0 transition-transform duration-500">
          <Button 
            className="w-full h-12 rounded-2xl flex items-center justify-center gap-2 text-sm font-bold shadow-2xl transition-colors"
            onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                addToCart(product);
            }}
          >
            <ShoppingBag className="w-4 h-4" />
            Add to Cart
          </Button>
        </div>
      </div>

      <div className="p-6 flex flex-col gap-2 grow">
        <div className="flex justify-between items-start gap-2">
            <Link href={`/product/${product.id}`} className="grow">
              <h3 className="text-lg font-black text-slate-900 dark:text-white line-clamp-1 group-hover:text-primary transition-colors">{product.name}</h3>
            </Link>
            <p className="text-lg font-black text-primary">
               {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(product.price)}
            </p>
        </div>
        <p className="text-xs text-slate-400 font-medium line-clamp-2 leading-relaxed">
           {product.description || 'Modern and elegant design crafted with premium materials for your comfort and style.'}
        </p>

        {/* Color Variants Dot Preview */}
        {product.color_varian && product.color_varian.length > 0 && (
          <div className="flex gap-1.5 mt-2">
            {product.color_varian.slice(0, 4).map((v) => (
              <div 
                key={v.id} 
                className="w-3 h-3 rounded-full border border-slate-200"
                style={{ backgroundColor: v.color }}
              />
            ))}
            {product.color_varian.length > 4 && (
              <span className="text-[10px] text-slate-400 font-bold ml-0.5">+{product.color_varian.length - 4}</span>
            )}
          </div>
        )}
      </div>
    </motion.div>
  );
};

export default ProductCard;
