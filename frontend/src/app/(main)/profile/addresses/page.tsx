"use client";

import React, { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Plus, Pencil, Trash2, MapPin } from 'lucide-react';
import { addressService, Address } from '@/services/api';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Card, CardContent, CardHeader } from '@/components/common/Card';
import { toast } from 'sonner';

// Address Schema
const addressSchema = z.object({
  recipient_name: z.string().min(3, "Name is too short"),
  recipient_phone_number: z.string().min(10, "Phone number is invalid"),
  full_address: z.string().min(10, "Address is too short"),
  city: z.string().min(1, "City is required"),
  province: z.string().min(1, "Province is required"),
  district: z.string().min(1, "District is required"),
  village: z.string().min(1, "Village is required"),
  postal_code: z.string().min(5, "Postal code is required"),
});

type AddressFormValues = z.infer<typeof addressSchema>;

export default function AddressesPage() {
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [loading, setLoading] = useState(true);
  const [isFormOpen, setIsFormOpen] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const form = useForm<AddressFormValues>({
    resolver: zodResolver(addressSchema),
    defaultValues: {
      recipient_name: '',
      recipient_phone_number: '',
      full_address: '',
      city: '',
      province: '',
      district: '',
      village: '',
      postal_code: '',
    },
  });

  const fetchAddresses = async () => {
    try {
      const data = await addressService.getAll();
      setAddresses(data);
    } catch (error) {
      console.error(error);
      toast.error("Failed to load addresses");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchAddresses();
  }, []);

  const handleEdit = (address: Address) => {
    form.reset({
      recipient_name: address.recipient_name,
      recipient_phone_number: address.recipient_phone_number,
      full_address: address.full_address,
      city: address.city,
      province: address.province,
      district: address.district,
      village: address.village,
      postal_code: address.postal_code,
    });
    setEditingId(address.id);
    setIsFormOpen(true);
    // Scroll to form
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this address?")) return;
    
    try {
      await addressService.delete(id);
      toast.success("Address deleted successfully");
      fetchAddresses();
    } catch (error) {
      toast.error("Failed to delete address");
    }
  };

  async function onSubmit(values: AddressFormValues) {
    setSubmitting(true);
    try {
      if (editingId) {
        await addressService.update(editingId, values);
        toast.success("Address updated successfully");
      } else {
        await addressService.create(values);
        toast.success("Address added successfully");
      }
      setIsFormOpen(false);
      setEditingId(null);
      form.reset();
      fetchAddresses();
    } catch (error: any) {
      console.error(error);
      toast.error(error.response?.data?.message || "Failed to save address");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-bold text-slate-900">My Addresses</h2>
          <p className="text-sm text-slate-500">Manage your shipping addresses</p>
        </div>
        {!isFormOpen && (
          <Button onClick={() => { 
            form.reset(); 
            setEditingId(null); 
            setIsFormOpen(true); 
          }}>
            <Plus className="w-4 h-4 mr-2" />
            Add New Address
          </Button>
        )}
      </div>

      {isFormOpen && (
        <Card className="shadow-sm border-slate-100 animate-in fade-in slide-in-from-top-4">
          <CardHeader className="border-b border-slate-100 flex flex-row items-center justify-between">
            <h3 className="font-bold text-slate-900">
              {editingId ? "Edit Address" : "Add New Address"}
            </h3>
            <button 
              onClick={() => setIsFormOpen(false)}
              className="text-slate-400 hover:text-slate-600 font-medium text-sm"
            >
              Cancel
            </button>
          </CardHeader>
          <CardContent className="pt-6">
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Input
                  label="Recipient Name"
                  placeholder="John Doe"
                  error={form.formState.errors.recipient_name?.message}
                  {...form.register('recipient_name')}
                />
                <Input
                  label="Phone Number"
                  placeholder="08123456789"
                  error={form.formState.errors.recipient_phone_number?.message}
                  {...form.register('recipient_phone_number')}
                />
              </div>

              <div className="space-y-1">
                 <label className="text-sm font-medium text-slate-700">Full Address</label>
                 <textarea
                   className={`w-full px-4 py-2 rounded-xl border ${
                     form.formState.errors.full_address 
                      ? 'border-red-500 focus:ring-red-500' 
                      : 'border-slate-200 focus:border-teal-500 focus:ring-teal-500'
                   } focus:outline-none focus:ring-1 transition-colors min-h-[80px]`}
                   placeholder="Street name, house number, etc."
                   {...form.register('full_address')}
                 />
                 {form.formState.errors.full_address && (
                   <p className="text-xs text-red-500 mt-1">{form.formState.errors.full_address.message}</p>
                 )}
              </div>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                 <Input
                  label="Province"
                  error={form.formState.errors.province?.message}
                  {...form.register('province')}
                />
                <Input
                  label="City"
                  error={form.formState.errors.city?.message}
                  {...form.register('city')}
                />
                 <Input
                  label="District"
                  error={form.formState.errors.district?.message}
                  {...form.register('district')}
                />
                <Input
                  label="Village"
                  error={form.formState.errors.village?.message}
                  {...form.register('village')}
                />
              </div>

              <div className="w-1/2 md:w-1/4">
                 <Input
                  label="Postal Code"
                  error={form.formState.errors.postal_code?.message}
                  {...form.register('postal_code')}
                />
              </div>

              <div className="flex justify-end pt-4">
                <Button type="submit" disabled={submitting} className="min-w-[120px]">
                  {submitting ? "Saving..." : "Save Address"}
                </Button>
              </div>
            </form>
          </CardContent>
        </Card>
      )}

      {loading ? (
        <div className="text-center py-12 text-slate-500">Loading addresses...</div>
      ) : addresses.length === 0 ? (
        !isFormOpen && (
          <div className="text-center py-12 bg-slate-50 rounded-2xl border border-dashed border-slate-200">
            <MapPin className="w-12 h-12 text-slate-300 mx-auto mb-4" />
            <h3 className="text-lg font-bold text-slate-900 mb-1">No Addresses Yet</h3>
            <p className="text-slate-500 mb-4">Add an address for shipping your orders.</p>
            <Button onClick={() => setIsFormOpen(true)} variant="outline">
              Add Address
            </Button>
          </div>
        )
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {addresses.map((address) => (
            <div 
              key={address.id} 
              className="bg-white p-6 rounded-2xl border border-slate-100 shadow-sm hover:shadow-md transition-shadow relative group"
            >
              <div className="absolute top-4 right-4 flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <button 
                  onClick={() => handleEdit(address)}
                  className="p-2 text-slate-400 hover:text-teal-500 hover:bg-teal-50 rounded-full transition-colors"
                >
                  <Pencil className="w-4 h-4" />
                </button>
                <button 
                  onClick={() => handleDelete(address.id)}
                  className="p-2 text-slate-400 hover:text-red-500 hover:bg-red-50 rounded-full transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>

              <h3 className="font-bold text-slate-900 mb-1">{address.recipient_name}</h3>
              <p className="text-sm text-slate-500 mb-4">{address.recipient_phone_number}</p>
              
              <div className="text-sm text-slate-600 space-y-1">
                <p>{address.full_address}</p>
                <p>{address.village}, {address.district}</p>
                <p>{address.city}, {address.province}, {address.postal_code}</p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
