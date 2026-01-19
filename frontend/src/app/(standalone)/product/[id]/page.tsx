"use client";

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useAuth } from '@/context/AuthContext';
import { productService } from '@/services/api';
import { Product, ColorVariant, SizeVariant } from '@/types/product';
import { Button } from '@/components/common/Button';
import { toast } from 'sonner';
import { ArrowLeft, Star, ShoppingBag, ShieldCheck, Zap } from 'lucide-react';
import Link from 'next/link';
import { useCart } from '@/context/CartContext';
import { ColorSelector } from '@/components/product/ColorSelector';
import { SizeSelector } from '@/components/product/SizeSelector';

function ProductSkeleton() {
  return (
    <div className="container mx-auto px-6 py-20 animate-pulse">
      <div className="h-8 bg-slate-200 w-1/4 mb-12 rounded-xl"></div>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-16">
        <div className="aspect-square bg-slate-200 rounded-[3rem]"></div>
        <div className="space-y-6">
            <div className="h-12 bg-slate-200 w-3/4 rounded-xl"></div>
            <div className="h-6 bg-slate-200 w-1/4 rounded-xl"></div>
            <div className="h-32 bg-slate-200 w-full rounded-xl"></div>
        </div>
      </div>
    </div>
  )
}

export default function ProductDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const { isAuthenticated } = useAuth();
  const { addToCart } = useCart();
  const router = useRouter();
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
      setCurrentImage(data.images || "https://images.unsplash.com/photo-1523275335684-37898b6baf30");
      
      // Auto-select first color variant if available
      if (data.color_varian && data.color_varian.length > 0) {
        setSelectedColorVariant(data.color_varian[0]);
      }
    } catch (error) {
      console.error(error);
      toast.error("Failed to load product");
    } finally {
      setLoading(false);
    }
  };

  const handleColorSelect = (variant: ColorVariant) => {
    setSelectedColorVariant(variant);
    setSelectedSizeVariant(undefined); // Reset size selection when color changes
    
    // Update image if variant has an image
    if (variant.images) {
      setCurrentImage(variant.images);
    } else if (product?.images) {
      setCurrentImage(product.images);
    }
  };

  const handleAddToCart = () => {
    if (!product) return;

    // Check if product has variants
    const hasColorVariants = product.color_varian && product.color_varian.length > 0;
    const hasSizeVariants = selectedColorVariant?.size_varian && selectedColorVariant.size_varian.length > 0;

    // Validate variant selection
    if (hasColorVariants && !selectedColorVariant) {
      toast.error("Please select a color variant");
      return;
    }

    if (hasSizeVariants && !selectedSizeVariant) {
      toast.error("Please select a size variant");
      return;
    }

    // Add to cart with selected variants
    addToCart(product, selectedColorVariant, selectedSizeVariant);
  };

  if (loading) return <ProductSkeleton />;
  if (!product) return <div className="min-h-screen flex items-center justify-center">Product not found</div>;

  return (
    <div className="min-h-screen flex flex-col">
        <div className="container mx-auto px-6 py-12 grow">
             <Link href="/" className="inline-flex items-center gap-2 text-slate-500 hover:text-teal-500 transition-colors mb-8 font-medium">
                <ArrowLeft className="w-4 h-4" />
                Back to Collection
             </Link>

             <div className="grid grid-cols-1 md:grid-cols-2 gap-12 lg:gap-24">
                {/* Product Image */}
                <div className="relative group">
                    <div className="aspect-square rounded-[3rem] overflow-hidden bg-slate-100 shadow-2xl relative z-10 border border-slate-200">
                        <img 
                            src={currentImage} 
                            alt={product.name}
                            className="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
                        />
                    </div>
                    {/* Decorative Blob */}
                    <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[120%] h-[120%] bg-teal-500/5 blur-3xl rounded-full z-0" />
                </div>

                {/* Product Details */}
                <div className="flex flex-col justify-center space-y-8">
                    <div>
                        <div className="flex items-center gap-2 mb-4">
                            <span className="px-3 py-1 bg-teal-50 text-teal-600 text-[10px] font-black uppercase tracking-widest rounded-full">
                                {product.category?.name || "Premium"}
                            </span>
                            <div className="flex items-center gap-1 text-amber-400">
                                <Star className="w-4 h-4 fill-current" />
                                <span className="text-slate-700 font-bold text-sm">4.9 (128 reviews)</span>
                            </div>
                        </div>

                        <h1 className="text-4xl md:text-5xl lg:text-6xl font-black text-slate-900 leading-[1.1] mb-6 italic tracking-tight">
                            {product.name}
                        </h1>
                        <div className="text-3xl font-medium text-teal-500 mb-8">
                            Rp {product.price.toLocaleString()}
                        </div>
                        
                        <p className="text-slate-500 leading-relaxed text-lg mb-8">
                            {product.description}
                        </p>
                    </div>

                    {/* Variant Selectors */}
                    <div className="space-y-6">
                        {/* Color Selector */}
                        {product.color_varian && product.color_varian.length > 0 && (
                            <ColorSelector
                                variants={product.color_varian}
                                selectedVariant={selectedColorVariant}
                                onSelect={handleColorSelect}
                            />
                        )}

                        {/* Size Selector */}
                        {selectedColorVariant?.size_varian && selectedColorVariant.size_varian.length > 0 && (
                            <SizeSelector
                                variants={selectedColorVariant.size_varian}
                                selectedVariant={selectedSizeVariant}
                                onSelect={setSelectedSizeVariant}
                            />
                        )}
                    </div>

                    <div className="space-y-6 pt-8 border-t border-slate-100">
                        {/* Features */}
                        <div className="grid grid-cols-2 gap-4">
                            <div className="flex items-center gap-3">
                                <div className="p-2 bg-slate-50 rounded-lg text-slate-600">
                                    <ShieldCheck className="w-5 h-5" />
                                </div>
                                <span className="text-sm font-bold text-slate-700">Authentic Guarantee</span>
                            </div>
                            <div className="flex items-center gap-3">
                                <div className="p-2 bg-slate-50 rounded-lg text-slate-600">
                                    <Zap className="w-5 h-5" />
                                </div>
                                <span className="text-sm font-bold text-slate-700">Express Shipping</span>
                            </div>
                        </div>

                        <Button 
                            className="w-full h-16 text-lg rounded-2xl gap-3 font-black uppercase tracking-widest shadow-xl shadow-teal-500/20 hover:shadow-teal-500/40 transition-all hover:-translate-y-1"
                            onClick={handleAddToCart}
                        >
                            <ShoppingBag className="w-5 h-5" />
                            Add to Cart
                        </Button>
                        
                        <p className="text-center text-xs text-slate-400 font-medium">
                            Free shipping on orders over Rp 5.000.000
                        </p>
                    </div>
                </div>
             </div>
        </div>
    </div>
  );
}
