"use client";

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { addressService, Address, User } from '@/services/api';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Plus, MapPin, Trash2 } from 'lucide-react';
import Link from 'next/link';
import { toast } from 'sonner';

export default function AddressPage() {
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [loading, setLoading] = useState(true);
  const router = useRouter();

  useEffect(() => {
    fetchAddresses();
  }, []);

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

  const handleDelete = async (id: number) => {
    if (!confirm("Are you sure you want to delete this address?")) return;
    try {
      await addressService.delete(id);
      setAddresses(addresses.filter((a) => a.id !== id));
      toast.success("Address deleted");
    } catch (error) {
       console.error(error);
       toast.error("Failed to delete address");
    }
  };

  return (
    <div className="container mx-auto px-6 py-12 max-w-4xl">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-black italic">MY ADDRESSES</h1>
        <Link href="/profile/address/new">
          <Button className="gap-2 font-bold">
            <Plus className="w-4 h-4" />
            Add New Address
          </Button>
        </Link>
      </div>

      {loading ? (
        <div className="space-y-4">
             {[1,2].map(i => <div key={i} className="h-32 bg-slate-100 dark:bg-slate-800 rounded-xl animate-pulse" />)}
        </div>
      ) : addresses.length === 0 ? (
        <div className="text-center py-20 bg-slate-50 dark:bg-slate-900 rounded-3xl">
           <MapPin className="w-12 h-12 text-slate-300 mx-auto mb-4" />
           <p className="text-slate-500 font-medium">No addresses found. Add one to speed up checkout.</p>
        </div>
      ) : (
        <div className="grid gap-4">
          {addresses.map((address) => (
            <Card key={address.id} className="premium-card relative group">
              <CardContent className="p-6 flex justify-between items-start">
                 <div className="space-y-1">
                    <div className="flex items-center gap-3">
                        <span className="font-bold text-lg">{address.name}</span>
                        {address.is_primary && <span className="px-2 py-0.5 bg-primary/10 text-primary text-[10px] uppercase font-bold tracking-wider rounded-full">Primary</span>}
                    </div>
                    <p className="text-slate-500 text-sm">{address.phone}</p>
                    <p className="text-slate-600 dark:text-slate-300 mt-2 font-medium">
                        {address.street}<br />
                        {address.city}, {address.state} {address.zip_code}
                    </p>
                 </div>
                 
                 <Button 
                    variant="ghost" 
                    size="icon" 
                    className="text-slate-400 hover:text-red-500"
                    onClick={() => handleDelete(address.id)}
                 >
                    <Trash2 className="w-4 h-4" />
                 </Button>
              </CardContent>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
