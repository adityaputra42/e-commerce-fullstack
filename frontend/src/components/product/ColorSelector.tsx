import React from 'react';
import { ColorVariant } from '@/types/product';
import { Check } from 'lucide-react';

interface ColorSelectorProps {
  variants: ColorVariant[];
  selectedVariant?: ColorVariant;
  onSelect: (variant: ColorVariant) => void;
}

export const ColorSelector: React.FC<ColorSelectorProps> = ({
  variants,
  selectedVariant,
  onSelect,
}) => {
  if (!variants || variants.length === 0) return null;

  return (
    <div className="space-y-3">
      <h3 className="text-sm font-bold text-slate-900 uppercase tracking-wide">
        Select Color
      </h3>
      <div className="flex flex-wrap gap-3">
        {variants.map((variant) => {
          const isSelected = selectedVariant?.id === variant.id;
          
          return (
            <button
              key={variant.id}
              onClick={() => onSelect(variant)}
              className={`
                group relative flex items-center gap-3 px-4 py-3 rounded-xl border-2 
                transition-all duration-200
                ${isSelected 
                  ? 'border-teal-500 bg-teal-50 shadow-md' 
                  : 'border-slate-200 bg-white hover:border-teal-300 hover:shadow-sm'
                }
              `}
            >
              {/* Color Swatch */}
              <div className="relative">
                <div
                  className="w-8 h-8 rounded-lg border-2 border-white shadow-md"
                  style={{ backgroundColor: variant.color }}
                />
                {isSelected && (
                  <div className="absolute inset-0 flex items-center justify-center">
                    <Check className="w-5 h-5 text-white drop-shadow-md" strokeWidth={3} />
                  </div>
                )}
              </div>
              
              {/* Color Name */}
              <span className={`text-sm font-bold ${isSelected ? 'text-teal-700' : 'text-slate-700'}`}>
                {variant.name}
              </span>
            </button>
          );
        })}
      </div>
    </div>
  );
};
