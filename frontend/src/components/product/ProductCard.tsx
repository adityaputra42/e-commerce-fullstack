'use client';

import { Product } from '@/types/product';
import Link from 'next/link';
import { ShoppingBag, Heart } from 'lucide-react';
import { motion } from 'framer-motion';
import { useCart } from '@/context/CartContext';

interface ProductCardProps {
  product: Product;
}
const ProductCard = ({ product }: ProductCardProps) => {
  const { addToCart } = useCart();

  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      className="bg-white border border-neutral-200 rounded-lg overflow-hidden group flex flex-col h-full hover:border-neutral-300 hover:shadow-sm transition-all"
    >
      <div className="relative aspect-4/5 overflow-hidden bg-neutral-50">
        <div className="absolute top-3 left-3 z-10">
          <span className="px-2 py-1 bg-neutral-900 text-white text-[10px] font-bold uppercase tracking-wider rounded-sm">
            {product.category?.name || 'New'}
          </span>
        </div>

        <button
          className="absolute top-3 right-3 z-10 w-8 h-8 bg-white border border-neutral-200 rounded-sm flex items-center justify-center text-neutral-400 hover:text-rose-500 hover:border-neutral-300 transition-colors"
          aria-label="Add to wishlist"
        >
          <Heart className="w-4 h-4" />
        </button>

        <Link href={`/product/${product.id}`}>
          <img
            src={product.images || 'https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&q=80'}
            alt={product.name}
            className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
          />
        </Link>
      </div>

      <div className="p-4 flex flex-col gap-2 grow">
        <Link href={`/product/${product.id}`}>
          <h3 className="text-sm font-bold text-neutral-900 line-clamp-1 group-hover:text-primary transition-colors">
            {product.name}
          </h3>
        </Link>

        {/* Color Variants Dot Preview */}
        {product.color_varian && product.color_varian.length > 0 && (
          <div className="flex gap-1.5">
            {product.color_varian.slice(0, 4).map((v) => (
              <div
                key={v.id}
                className="w-3 h-3 rounded-full border border-neutral-200"
                style={{ backgroundColor: v.color }}
              />
            ))}
            {product.color_varian.length > 4 && (
              <span className="text-[10px] text-neutral-400 font-bold ml-0.5">+{product.color_varian.length - 4}</span>
            )}
          </div>
        )}

        <div className="mt-auto pt-1 flex items-center justify-between gap-2">
          <p className="text-base font-black text-neutral-900">
            {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(product.price)}
          </p>
          <button
            onClick={(e) => {
              e.preventDefault();
              e.stopPropagation();
              addToCart(product);
            }}
            className="w-9 h-9 shrink-0 bg-neutral-900 text-white rounded-sm flex items-center justify-center hover:bg-black active:scale-[0.96] transition-all"
            aria-label="Add to cart"
          >
            <ShoppingBag className="w-4 h-4" />
          </button>
        </div>
      </div>
    </motion.div>
  );
};

export default ProductCard;
