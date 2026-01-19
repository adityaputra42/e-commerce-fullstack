import React from 'react';
import { cn } from '@/lib/utils';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  children: React.ReactNode;
}

export const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'primary', size = 'md', children, ...props }, ref) => {
    const baseStyles = 'font-bold rounded-xl transition-all duration-200 inline-flex items-center justify-center disabled:opacity-50 disabled:cursor-not-allowed';
    
    const variants = {
      primary: 'bg-teal-500 text-white hover:bg-teal-600 active:scale-95 shadow-md hover:shadow-lg',
      secondary: 'bg-purple-500 text-white hover:bg-purple-600 active:scale-95 shadow-md hover:shadow-lg',
      outline: 'bg-transparent border-2 border-slate-300 text-slate-700 hover:border-teal-500 hover:text-teal-600 active:scale-95',
      ghost: 'bg-transparent text-slate-700 hover:bg-slate-100 active:scale-95',
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
