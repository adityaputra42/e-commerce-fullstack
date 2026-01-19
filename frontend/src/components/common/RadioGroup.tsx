import React from 'react';
import { cn } from '@/lib/utils';

interface RadioGroupProps {
  value?: string;
  onValueChange?: (value: string) => void;
  className?: string;
  children: React.ReactNode;
}

export const RadioGroup: React.FC<RadioGroupProps> = ({
  value,
  onValueChange,
  className,
  children,
}) => {
  return (
    <div className={cn('space-y-2', className)} role="radiogroup">
      {React.Children.map(children, (child) => {
        if (React.isValidElement(child)) {
          return React.cloneElement(child as React.ReactElement<any>, {
            checked: child.props.value === value,
            onChange: () => onValueChange?.(child.props.value),
          });
        }
        return child;
      })}
    </div>
  );
};

interface RadioGroupItemProps extends React.InputHTMLAttributes<HTMLInputElement> {
  value: string;
  checked?: boolean;
}

export const RadioGroupItem = React.forwardRef<HTMLInputElement, RadioGroupItemProps>(
  ({ className, value, checked, onChange, ...props }, ref) => {
    return (
      <input
        ref={ref}
        type="radio"
        value={value}
        checked={checked}
        onChange={onChange}
        className={cn(
          'w-4 h-4 text-teal-500 border-slate-300 focus:ring-2 focus:ring-teal-500',
          className
        )}
        {...props}
      />
    );
  }
);

RadioGroupItem.displayName = 'RadioGroupItem';
