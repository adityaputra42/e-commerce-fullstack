"use client";

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { transactionService, Transaction } from '@/services/api';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { toast } from 'sonner';
import { ArrowLeft, MapPin, Truck, CreditCard, ShoppingBag } from 'lucide-react';
import Link from 'next/link';

export default function TransactionDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const [transaction, setTransaction] = useState<Transaction | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (id) fetchTransaction();
  }, [id]);

  const fetchTransaction = async () => {
    try {
      const data = await transactionService.getById(id);
      setTransaction(data);
    } catch (error) {
      console.error(error);
      toast.error("Failed to load transaction details");
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div className="container mx-auto py-20 text-center">Loading details...</div>;
  if (!transaction) return <div className="container mx-auto py-20 text-center">Transaction not found</div>;

  return (
    <div className="container mx-auto px-6 py-12 max-w-4xl">
        <Link href="/profile/transactions" className="inline-flex items-center gap-2 text-slate-500 hover:text-primary mb-8 transition-colors">
            <ArrowLeft className="w-4 h-4" />
            Back to History
        </Link>
        
        <div className="flex flex-col md:flex-row justify-between md:items-start gap-4 mb-8">
            <div>
                <h1 className="text-3xl font-black italic mb-2">ORDER #{transaction.tx_id.slice(0,8)}</h1>
                <p className="text-slate-500">{new Date(transaction.created_at).toLocaleString()}</p>
            </div>
            <div className="px-4 py-2 bg-slate-900 text-white dark:bg-white dark:text-slate-900 font-bold uppercase tracking-widest rounded-lg text-sm">
                {transaction.status}
            </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {/* Main Content */}
            <div className="lg:col-span-2 space-y-8">
                {/* Items */}
                <Card className="premium-card">
                    <CardContent className="p-6 space-y-6">
                        <h2 className="font-bold text-lg flex items-center gap-2">
                            <ShoppingBag className="w-5 h-5 text-primary" /> Items
                        </h2>
                        {transaction.orders?.map((order) => (
                            <div key={order.id} className="flex gap-4">
                                <div className="w-20 h-20 bg-slate-100 rounded-lg shrink-0 overflow-hidden">
                                     <img src={order.product?.images} alt="" className="w-full h-full object-cover" />
                                </div>
                                <div className="grow">
                                    <h3 className="font-bold text-slate-900 dark:text-white line-clamp-1">{order.product?.name}</h3>
                                    <p className="text-sm text-slate-500 mb-1">Qty: {order.quantity}</p>
                                    <p className="font-bold text-primary">
                                        {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(order.unit_price)}
                                    </p>
                                </div>
                            </div>
                        ))}
                    </CardContent>
                </Card>

                {/* Shipping Info */}
                <Card className="premium-card">
                    <CardContent className="p-6 space-y-6">
                        <h2 className="font-bold text-lg flex items-center gap-2">
                            <MapPin className="w-5 h-5 text-primary" /> Shipping Details
                        </h2>
                        {transaction.address && (
                            <div className="p-4 bg-slate-50 dark:bg-slate-900 rounded-xl">
                                <p className="font-bold">{transaction.address.name}</p>
                                <p className="text-sm text-slate-500 mb-2">{transaction.address.phone}</p>
                                <p className="text-sm text-slate-600 dark:text-slate-300">
                                    {transaction.address.street}<br />
                                    {transaction.address.city}, {transaction.address.state} {transaction.address.zip_code}
                                </p>
                            </div>
                        )}
                        <div className="flex items-center gap-4">
                            <div className="flex items-center gap-2 text-sm font-medium text-slate-600">
                                <Truck className="w-4 h-4" />
                                {transaction.shipping?.name}
                            </div>
                            <div className="flex items-center gap-2 text-sm font-medium text-slate-600">
                                <CreditCard className="w-4 h-4" />
                                {transaction.payment_method?.name}
                            </div>
                        </div>
                    </CardContent>
                </Card>
            </div>

            {/* Summary */}
            <div className="lg:col-span-1">
                <Card className="premium-card sticky top-24">
                    <CardContent className="p-6 space-y-4">
                        <h2 className="font-bold text-lg">Order Summary</h2>
                        <div className="flex justify-between text-slate-500 text-sm">
                             <span>Shipping Cost</span>
                             <span>{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(transaction.shipping_price)}</span>
                        </div>
                        <div className="border-t pt-4 flex justify-between font-black text-lg">
                             <span>Total</span>
                             <span className="text-primary">{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(transaction.total_price)}</span>
                        </div>
                    </CardContent>
                </Card>
            </div>
        </div>
    </div>
  );
}
