import React from 'react';
import { cn } from '@/lib/utils';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  children: React.ReactNode;
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', size = 'md', children, ...props }, ref) => {
    const baseStyles = 'font-bold rounded-md transition-all duration-150 inline-flex items-center justify-center disabled:opacity-50 disabled:cursor-not-allowed';

    // Flat, boxy variants: no drop shadows, no bouncy hover-scale — a
    // calm color swap on hover and a very slight press-down on click.
    // `secondary` moved from purple-filled to ink-filled (matches the
    // black CTA buttons — "Quick Add", "Explore Collection" — in the
    // boxy reference) instead of introducing a second brand hue.
    const variants = {
      primary: 'bg-primary text-white hover:bg-primary-dark active:scale-[0.98]',
      secondary: 'bg-neutral-900 text-white hover:bg-black active:scale-[0.98]',
      outline: 'bg-transparent border border-neutral-300 text-neutral-800 hover:border-neutral-900 hover:text-neutral-900 active:scale-[0.98]',
      ghost: 'bg-transparent text-neutral-700 hover:bg-neutral-100 active:scale-[0.98]',
    };

    const sizes = {
      sm: 'px-4 py-2 text-sm',
      md: 'px-6 py-3 text-base',
      lg: 'px-8 py-4 text-lg',
    };

    return (
      <button
        ref={ref}
        className={cn(baseStyles, variants[variant], sizes[size], className)}
        {...props}
      >
        {children}
      </button>
    );
  }
);

Button.displayName = 'Button';
