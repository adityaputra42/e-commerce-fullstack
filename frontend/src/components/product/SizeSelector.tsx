import React from 'react';
import { SizeVariant } from '@/types/product';
import { Check } from 'lucide-react';

interface SizeSelectorProps {
  variants: SizeVariant[];
  selectedVariant?: SizeVariant;
  onSelect: (variant: SizeVariant) => void;
}

export const SizeSelector: React.FC<SizeSelectorProps> = ({
  variants,
  selectedVariant,
  onSelect,
}) => {
  if (!variants || variants.length === 0) return null;

  return (
    <div className="space-y-3">
      <h3 className="text-sm font-bold text-slate-900 uppercase tracking-wide">
        Select Size
      </h3>
      <div className="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-3">
        {variants.map((variant) => {
          const isSelected = selectedVariant?.id === variant.id;
          const isOutOfStock = variant.stock === 0;
          
          return (
            <button
              key={variant.id}
              onClick={() => !isOutOfStock && onSelect(variant)}
              disabled={isOutOfStock}
              className={`
                relative px-4 py-3 rounded-xl border-2 font-bold text-sm
                transition-all duration-200
                ${isOutOfStock
                  ? 'border-slate-200 bg-slate-50 text-slate-400 cursor-not-allowed opacity-50'
                  : isSelected
                    ? 'border-teal-500 bg-teal-500 text-white shadow-md shadow-teal-200'
                    : 'border-slate-200 bg-white text-slate-700 hover:border-teal-300 hover:shadow-sm active:scale-95'
                }
              `}
            >
              <div className="flex flex-col items-center gap-1">
                <span className="text-base">{variant.size}</span>
                {!isOutOfStock && (
                  <span className={`text-[10px] ${isSelected ? 'text-teal-100' : 'text-slate-500'}`}>
                    Stock: {variant.stock}
                  </span>
                )}
                {isOutOfStock && (
                  <span className="text-[10px] text-red-500 font-bold">
                    Out of Stock
                  </span>
                )}
              </div>
              
              {isSelected && !isOutOfStock && (
                <div className="absolute -top-1 -right-1 w-5 h-5 bg-white rounded-full flex items-center justify-center shadow-md">
                  <Check className="w-3 h-3 text-teal-500" strokeWidth={3} />
                </div>
              )}
            </button>
          );
        })}
      </div>
    </div>
  );
};
