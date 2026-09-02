"use client";

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { productService } from '@/services/api';
import { Product, ColorVariant, SizeVariant } from '@/types/product';
import { Button } from '@/components/common/Button';
import { toast } from 'sonner';
import { ArrowLeft, ShoppingBag, ShieldCheck, Truck } from 'lucide-react';
import Link from 'next/link';
import { useCart } from '@/context/CartContext';
import { ColorSelector } from '@/components/product/ColorSelector';
import { SizeSelector } from '@/components/product/SizeSelector';

function ProductSkeleton() {
  return (
    <div className="container mx-auto px-6 py-12 animate-pulse">
      <div className="h-6 bg-neutral-200 w-1/4 mb-10 rounded-md"></div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-12">
        <div className="aspect-square bg-neutral-200 rounded-lg"></div>
        <div className="space-y-5">
          <div className="h-10 bg-neutral-200 w-3/4 rounded-md"></div>
          <div className="h-6 bg-neutral-200 w-1/4 rounded-md"></div>
          <div className="h-28 bg-neutral-200 w-full rounded-md"></div>
        </div>
      </div>
    </div>
  );
}

export default function ProductDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const { addToCart } = useCart();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [selectedColorVariant, setSelectedColorVariant] = useState<ColorVariant | undefined>();
  const [selectedSizeVariant, setSelectedSizeVariant] = useState<SizeVariant | undefined>();
  const [currentImage, setCurrentImage] = useState<string>('');

  useEffect(() => {
    if (id) {
      fetchProduct();
    } else {
      setLoading(false);
    }
  }, [id]);

  const fetchProduct = async () => {
    try {
      const data = await productService.getById(id);
      setProduct(data);
      setCurrentImage(data.images || '');

      if (data.color_varian && data.color_varian.length > 0) {
        setSelectedColorVariant(data.color_varian[0]);
      }
    } catch (error) {
      console.error(error);
      toast.error('Failed to load product');
    } finally {
      setLoading(false);
    }
  };

  const handleColorSelect = (variant: ColorVariant) => {
    setSelectedColorVariant(variant);
    setSelectedSizeVariant(undefined);

    if (variant.images) {
      setCurrentImage(variant.images);
    } else if (product?.images) {
      setCurrentImage(product.images);
    }
  };

  const handleAddToCart = () => {
    if (!product) return;

    const hasColorVariants = product.color_varian && product.color_varian.length > 0;
    const hasSizeVariants = selectedColorVariant?.size_varian && selectedColorVariant.size_varian.length > 0;

    if (hasColorVariants && !selectedColorVariant) {
      toast.error('Please select a color variant');
      return;
    }

    if (hasSizeVariants && !selectedSizeVariant) {
      toast.error('Please select a size variant');
      return;
    }

    addToCart(product, selectedColorVariant, selectedSizeVariant);
  };

  if (loading) return <ProductSkeleton />;
  if (!product) return <div className="min-h-screen flex items-center justify-center text-neutral-500 font-medium">Product not found</div>;

  return (
    <div className="min-h-screen flex flex-col">
      <div className="container mx-auto px-6 py-10 grow">
        <Link href="/shop" className="inline-flex items-center gap-2 text-neutral-500 hover:text-primary transition-colors mb-8 text-sm font-semibold">
          <ArrowLeft className="w-4 h-4" />
          Back to Shop
        </Link>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-12 lg:gap-16">
          {/* Product Image — plain bordered box, no rotated blob shadow */}
          <div className="aspect-square rounded-lg overflow-hidden bg-neutral-50 border border-neutral-200">
            {currentImage ? (
              <img
                src={currentImage}
                alt={product.name}
                className="w-full h-full object-cover"
              />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-neutral-300 text-sm font-bold">
                No image available
              </div>
            )}
          </div>

          {/* Product Details */}
          <div className="flex flex-col justify-center space-y-6">
            <div>
              <span className="inline-block px-2.5 py-1 bg-primary/10 text-primary text-[10px] font-bold uppercase tracking-widest rounded-sm mb-4">
                {product.category?.name || 'Product'}
              </span>

              <h1 className="text-3xl md:text-4xl font-black text-neutral-900 leading-tight mb-4">
                {product.name}
              </h1>
              <div className="text-2xl font-black text-neutral-900 mb-6">
                {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(product.price)}
              </div>

              <p className="text-neutral-500 leading-relaxed">
                {product.description}
              </p>
            </div>

            {/* Variant Selectors */}
            <div className="space-y-5">
              {product.color_varian && product.color_varian.length > 0 && (
                <ColorSelector
                  variants={product.color_varian}
                  selectedVariant={selectedColorVariant}
                  onSelect={handleColorSelect}
                />
              )}

              {selectedColorVariant?.size_varian && selectedColorVariant.size_varian.length > 0 && (
                <SizeSelector
                  variants={selectedColorVariant.size_varian}
                  selectedVariant={selectedSizeVariant}
                  onSelect={setSelectedSizeVariant}
                />
              )}
            </div>

            <div className="space-y-5 pt-6 border-t border-neutral-200">
              <div className="grid grid-cols-2 gap-4">
                <div className="flex items-center gap-2.5">
                  <div className="w-9 h-9 bg-neutral-50 border border-neutral-200 rounded-md flex items-center justify-center text-neutral-500 shrink-0">
                    <ShieldCheck className="w-4 h-4" />
                  </div>
                  <span className="text-xs font-bold text-neutral-700">Authentic Guarantee</span>
                </div>
                <div className="flex items-center gap-2.5">
                  <div className="w-9 h-9 bg-neutral-50 border border-neutral-200 rounded-md flex items-center justify-center text-neutral-500 shrink-0">
                    <Truck className="w-4 h-4" />
                  </div>
                  <span className="text-xs font-bold text-neutral-700">Fast Shipping</span>
                </div>
              </div>

              <Button
                className="w-full h-14 text-base gap-2"
                onClick={handleAddToCart}
              >
                <ShoppingBag className="w-4 h-4" />
                Add to Cart
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
