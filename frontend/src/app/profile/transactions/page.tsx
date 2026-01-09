"use client";

import { useEffect, useState } from 'react';
import { transactionService, Transaction } from '@/services/api';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { toast } from 'sonner';
import { Package, Calendar, ChevronRight, ShoppingBag } from 'lucide-react';
import Link from 'next/link';

export default function TransactionsPage() {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchTransactions();
  }, []);

  const fetchTransactions = async () => {
    try {
      const data = await transactionService.getAll();
      setTransactions(data || []);
    } catch (error) {
      console.error(error);
      toast.error("Failed to load transactions");
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
        case 'paid': return 'text-emerald-500 bg-emerald-50';
        case 'shipped': return 'text-blue-500 bg-blue-50';
        case 'completed': return 'text-teal-500 bg-teal-50';
        case 'cancelled': return 'text-red-500 bg-red-50';
        default: return 'text-amber-500 bg-amber-50';
    }
  };

  if (loading) return <div className="container mx-auto py-20 text-center">Loading history...</div>;

  return (
    <div className="container mx-auto px-6 py-12 max-w-4xl">
       <h1 className="text-3xl font-black italic mb-8">ORDER HISTORY</h1>

       {transactions.length === 0 ? (
           <div className="text-center py-20 bg-slate-50 dark:bg-slate-900 rounded-3xl">
              <Package className="w-12 h-12 text-slate-300 mx-auto mb-4" />
              <p className="text-slate-500 font-medium">No orders found.</p>
              <Link href="/">
                 <Button variant="link" className="mt-2">Start Shopping</Button>
              </Link>
           </div>
       ) : (
           <div className="space-y-6">
               {transactions.map((tx) => (
                   <Link href={`/profile/transactions/${tx.tx_id}`} key={tx.tx_id} className="block group">
                       <Card className="premium-card hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                           <CardContent className="p-6">
                               <div className="flex flex-col md:flex-row justify-between md:items-center gap-4">
                                   <div className="space-y-1">
                                       <div className="flex items-center gap-3">
                                            <span className="font-bold text-lg">#{tx.tx_id.slice(0,8)}</span>
                                            <span className={`px-2.5 py-0.5 text-[10px] font-black uppercase tracking-widest rounded-full ${getStatusColor(tx.status)}`}>
                                                {tx.status}
                                            </span>
                                       </div>
                                       <div className="flex items-center gap-2 text-slate-400 text-sm font-medium">
                                            <Calendar className="w-4 h-4" />
                                            {new Date(tx.created_at).toLocaleDateString(undefined, {
                                                year: 'numeric', month: 'long', day: 'numeric'
                                            })}
                                       </div>
                                   </div>

                                   <div className="flex items-center gap-6">
                                       <div className="text-right">
                                           <div className="text-xs text-slate-400 font-bold uppercase tracking-wider mb-0.5">Total Amount</div>
                                           <div className="font-black text-primary text-lg">
                                               {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(tx.total_price)}
                                           </div>
                                       </div>
                                       <div className="w-10 h-10 rounded-full bg-white dark:bg-slate-700 shadow-sm flex items-center justify-center text-slate-400 group-hover:text-primary group-hover:translate-x-1 transition-all">
                                           <ChevronRight className="w-5 h-5" />
                                       </div>
                                   </div>
                               </div>
                               
                               <div className="mt-6 pt-6 border-t border-slate-100 dark:border-slate-800 flex gap-4 overflow-hidden">
                                    {tx.orders?.slice(0, 3).map((order) => (
                                        <div key={order.id} className="w-16 h-16 bg-slate-100 rounded-lg shrink-0 overflow-hidden relative">
                                            {order.product?.images ? (
                                                <img src={order.product.images} alt="" className="w-full h-full object-cover" />
                                            ) : (
                                                <div className="w-full h-full flex items-center justify-center text-slate-300">
                                                    <ShoppingBag className="w-6 h-6" />
                                                </div>
                                            )}
                                        </div>
                                    ))}
                                    {tx.orders?.length > 3 && (
                                        <div className="w-16 h-16 bg-slate-50 rounded-lg flex items-center justify-center text-slate-400 font-bold text-xs">
                                            +{tx.orders.length - 3}
                                        </div>
                                    )}
                               </div>
                           </CardContent>
                       </Card>
                   </Link>
               ))}
           </div>
       )}
    </div>
  );
}
