"use client";

import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod'; // Import zod correctly
import { addressService } from '@/services/api';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { toast } from 'sonner';
import { useRouter } from 'next/navigation';
import { ArrowLeft, MapPin } from 'lucide-react';
import Link from 'next/link';

const addressSchema = z.object({
  name: z.string().min(2, "Name is required"),
  phone: z.string().min(10, "Valid phone number is required"),
  street: z.string().min(5, "Street address is required"),
  city: z.string().min(2, "City is required"),
  state: z.string().min(2, "State is required"),
  zip_code: z.string().min(5, "Zip code is required"),
});

export default function NewAddressPage() {
  const router = useRouter();
  const form = useForm<z.infer<typeof addressSchema>>({
    resolver: zodResolver(addressSchema),
    defaultValues: {
      name: "",
      phone: "",
      street: "",
      city: "",
      state: "",
      zip_code: "",
    },
  });

  async function onSubmit(values: z.infer<typeof addressSchema>) {
    try {
      await addressService.create({
          ...values,
          is_primary: true // Default to primary for now or logic to check if first address
      });
      toast.success("Address added successfully");
      router.push('/profile/address');
    } catch (error) {
      console.error(error);
      toast.error("Failed to add address");
    }
  }

  return (
    <div className="container mx-auto px-6 py-12 max-w-2xl">
      <Link href="/profile/address" className="inline-flex items-center gap-2 text-slate-500 hover:text-primary mb-8 transition-colors">
         <ArrowLeft className="w-4 h-4" />
         Back to Addresses
      </Link>

      <div className="mb-8">
         <div className="w-12 h-12 bg-primary/10 rounded-2xl flex items-center justify-center text-primary mb-4">
            <MapPin className="w-6 h-6" />
         </div>
         <h1 className="text-3xl font-black italic">ADD NEW ADDRESS</h1>
         <p className="text-slate-500">Fill in the details for your shipping destination.</p>
      </div>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <FormField
                control={form.control}
                name="name"
                render={({ field }) => (
                <FormItem>
                    <FormLabel>Full Name</FormLabel>
                    <FormControl>
                    <Input placeholder="John Doe" {...field} className="h-12 rounded-xl" />
                    </FormControl>
                    <FormMessage />
                </FormItem>
                )}
            />
            <FormField
                control={form.control}
                name="phone"
                render={({ field }) => (
                <FormItem>
                    <FormLabel>Phone Number</FormLabel>
                    <FormControl>
                    <Input placeholder="0812..." {...field} className="h-12 rounded-xl" />
                    </FormControl>
                    <FormMessage />
                </FormItem>
                )}
            />
          </div>

          <FormField
            control={form.control}
            name="street"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Street Address</FormLabel>
                <FormControl>
                  <Input placeholder="Jl. Sudirman No. 123" {...field} className="h-12 rounded-xl" />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
             <FormField
                control={form.control}
                name="city"
                render={({ field }) => (
                <FormItem>
                    <FormLabel>City</FormLabel>
                    <FormControl>
                    <Input placeholder="Jakarta Selatan" {...field} className="h-12 rounded-xl" />
                    </FormControl>
                    <FormMessage />
                </FormItem>
                )}
            />
             <FormField
                control={form.control}
                name="state"
                render={({ field }) => (
                <FormItem>
                    <FormLabel>State/Province</FormLabel>
                    <FormControl>
                    <Input placeholder="DKI Jakarta" {...field} className="h-12 rounded-xl" />
                    </FormControl>
                    <FormMessage />
                </FormItem>
                )}
            />
             <FormField
                control={form.control}
                name="zip_code"
                render={({ field }) => (
                <FormItem>
                    <FormLabel>Zip Code</FormLabel>
                    <FormControl>
                    <Input placeholder="12190" {...field} className="h-12 rounded-xl" />
                    </FormControl>
                    <FormMessage />
                </FormItem>
                )}
            />
          </div>

          <Button type="submit" className="w-full h-12 rounded-xl text-lg font-bold shadow-xl shadow-primary/20" disabled={form.formState.isSubmitting}>
            {form.formState.isSubmitting ? "Saving..." : "Save Address"}
          </Button>
        </form>
      </Form>
    </div>
  );
}
