"use client";

import { useEffect, useState } from 'react';
import { useCart } from '@/context/CartContext';
import { useAuth } from '@/context/AuthContext';
import { useRouter } from 'next/navigation';
import { addressService, shippingService, paymentService, transactionService, Address, ShippingMethod, PaymentMethod } from '@/services/api';
import { Button } from '@/components/common/Button';
import { Card, CardContent } from '@/components/common/Card';
import { RadioGroup, RadioGroupItem } from "@/components/common/RadioGroup";
import { Label } from "@/components/common/Label";
import { toast } from 'sonner';
import { MapPin, Truck, CreditCard, ArrowRight, Plus } from 'lucide-react';
import Link from 'next/link';
import { Skeleton } from '@/components/common/Button';

export default function CheckoutPage() {
  const { items, cartTotal, clearCart } = useCart();
  const { isAuthenticated, user } = useAuth();
  const router = useRouter();

  const [addresses, setAddresses] = useState<Address[]>([]);
  const [shippingMethods, setShippingMethods] = useState<ShippingMethod[]>([]);
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);

  const [selectedAddressId, setSelectedAddressId] = useState<number | null>(null);
  const [selectedShippingId, setSelectedShippingId] = useState<number | null>(null);
  const [selectedPaymentId, setSelectedPaymentId] = useState<number | null>(null);

  const [loading, setLoading] = useState(true);
  const [processing, setProcessing] = useState(false);

  useEffect(() => {
    if (!isAuthenticated) {
        router.push('/login?redirect=/checkout');
        return;
    }
    if (items.length === 0) {
        router.push('/cart');
        return;
    }
    fetchData();
  }, [isAuthenticated, items]);

  const fetchData = async () => {
    try {
        const [addrData, shipData, payData] = await Promise.all([
            addressService.getAll(),
            shippingService.getAll(),
            paymentService.getAllMethods()
        ]);
        setAddresses(addrData);
        if (addrData.length > 0) setSelectedAddressId(addrData[0].id);
        
        setShippingMethods(shipData);
        if (shipData.length > 0) setSelectedShippingId(shipData[0].id);

        setPaymentMethods(payData);
        if (payData.length > 0) setSelectedPaymentId(payData[0].id);

    } catch (error) {
        console.error(error);
        toast.error("Failed to load checkout data");
    } finally {
        setLoading(false);
    }
  };

  const handlePlaceOrder = async () => {
    if (!selectedAddressId || !selectedShippingId || !selectedPaymentId) {
        toast.error("Please select address, shipping, and payment method");
        return;
    }

    setProcessing(true);
    try {
        const productOrders = items.map(item => ({
            product_id: item.id,
            quantity: item.quantity
        }));

        await transactionService.createTransaction({
            address_id: selectedAddressId,
            shipping_id: selectedShippingId,
            payment_method_id: selectedPaymentId,
            product_orders: productOrders
        });

        toast.success("Order placed successfully!");
        clearCart();
        router.push('/profile/transactions');

    } catch (error) {
        console.error(error);
        toast.error("Failed to place order");
    } finally {
        setProcessing(false);
    }
  };

  const selectedShipping = shippingMethods.find(s => s.id === selectedShippingId);
  const shippingCost = selectedShipping ? selectedShipping.cost : 0;
  const grandTotal = cartTotal + shippingCost;

  if (loading) return <div className="container mx-auto py-20 text-center">Loading checkout...</div>;

  return (
    <div className="container mx-auto px-6 py-12">
        <h1 className="text-3xl font-black italic mb-8">It's almost yours!</h1>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-12">
            <div className="lg:col-span-2 space-y-8">
                {/* 1. Address Selection */}
                <section>
                    <div className="flex justify-between items-center mb-4">
                        <h2 className="text-xl font-bold flex items-center gap-2">
                            <MapPin className="text-primary" /> Shipping Address
                        </h2>
                        <Link href="/profile/address/new" className="text-sm font-bold text-primary hover:underline flex items-center gap-1">
                            <Plus className="w-4 h-4" /> Add New
                        </Link>
                    </div>
                    
                    {addresses.length === 0 ? (
                         <div className="p-8 border-2 border-dashed border-slate-200 rounded-2xl text-center">
                            <p className="text-slate-500 mb-4">You have no saved addresses.</p>
                            <Link href="/profile/address/new">
                                <Button>Add Address</Button>
                            </Link>
                         </div>
                    ) : (
                        <RadioGroup value={selectedAddressId?.toString()} onValueChange={(v: string) => setSelectedAddressId(Number(v))} className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            {addresses.map((addr) => (
                                <div key={addr.id} className={`relative`}>
                                    <RadioGroupItem value={addr.id.toString()} id={`addr-${addr.id}`} className="peer sr-only" />
                                    <Label 
                                        htmlFor={`addr-${addr.id}`}
                                        className="flex flex-col p-6 rounded-2xl border-2 border-slate-100 bg-white dark:bg-slate-900 cursor-pointer transition-all peer-checked:border-primary peer-checked:bg-primary/5 hover:border-slate-200"
                                    >
                                        <span className="font-bold text-lg mb-1">{addr.name}</span>
                                        <span className="text-slate-500 text-sm mb-2">{addr.phone}</span>
                                        <span className="text-slate-600 dark:text-slate-300 text-sm leading-relaxed">
                                            {addr.street}, {addr.city}, {addr.state} {addr.zip_code}
                                        </span>
                                    </Label>
                                </div>
                            ))}
                        </RadioGroup>
                    )}
                </section>

                {/* 2. Shipping Method */}
                <section>
                    <h2 className="text-xl font-bold flex items-center gap-2 mb-4">
                        <Truck className="text-primary" /> Shipping Method
                    </h2>
                    <RadioGroup value={selectedShippingId?.toString()} onValueChange={(v: string) => setSelectedShippingId(Number(v))} className="space-y-3">
                         {shippingMethods.map((ship) => (
                             <div key={ship.id} className="flex items-center space-x-2">
                                 <RadioGroupItem value={ship.id.toString()} id={`ship-${ship.id}`} className="peer sr-only" />
                                 <Label 
                                    htmlFor={`ship-${ship.id}`}
                                    className="flex items-center justify-between w-full p-4 rounded-xl border border-slate-200 cursor-pointer peer-checked:border-primary peer-checked:bg-primary/5 hover:bg-slate-50 dark:hover:bg-slate-800"
                                 >
                                     <div>
                                         <div className="font-bold">{ship.name}</div>
                                         <div className="text-xs text-slate-500">Est: {ship.estimated_delivery}</div>
                                     </div>
                                     <div className="font-bold text-primary">
                                        {new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(ship.cost)}
                                     </div>
                                 </Label>
                             </div>
                         ))}
                    </RadioGroup>
                </section>

                {/* 3. Payment Method */}
                <section>
                    <h2 className="text-xl font-bold flex items-center gap-2 mb-4">
                        <CreditCard className="text-primary" /> Payment Method
                    </h2>
                    <RadioGroup value={selectedPaymentId?.toString()} onValueChange={(v: string) => setSelectedPaymentId(Number(v))} className="grid grid-cols-2 md:grid-cols-3 gap-4">
                        {paymentMethods.map((pay) => (
                            <div key={pay.id}>
                                <RadioGroupItem value={pay.id.toString()} id={`pay-${pay.id}`} className="peer sr-only" />
                                <Label 
                                    htmlFor={`pay-${pay.id}`}
                                    className="flex flex-col items-center justify-center p-6 rounded-xl border-2 border-slate-100 cursor-pointer peer-checked:border-primary peer-checked:bg-primary/5 hover:border-slate-200 transition-all text-center h-full"
                                >
                                    <span className="font-bold">{pay.name}</span>
                                    <span className="text-xs text-slate-500 mt-1 capitalize">{pay.type.replace('_', ' ')}</span>
                                </Label>
                            </div>
                        ))}
                    </RadioGroup>
                </section>
            </div>

            {/* Order Summary */}
            <div className="lg:col-span-1">
                <Card className="premium-card sticky top-24">
                   <CardContent className="p-8 space-y-6">
                      <h2 className="text-xl font-bold">Order Summary</h2>
                      
                      <div className="space-y-4">
                         {items.map(item => (
                             <div key={item.id} className="flex justify-between text-sm">
                                 <span className="text-slate-500">{item.name} x{item.quantity}</span>
                                 <span>{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(item.price * item.quantity)}</span>
                             </div>
                         ))}
                         <hr />
                         <div className="flex justify-between text-slate-500">
                             <span>Subtotal</span>
                             <span>{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(cartTotal)}</span>
                         </div>
                         <div className="flex justify-between text-slate-500">
                             <span>Shipping</span>
                             <span>{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(shippingCost)}</span>
                         </div>
                      </div>
                      
                      <div className="border-t pt-4 flex justify-between font-black text-lg">
                         <span>Total</span>
                         <span className="text-primary">{new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', maximumFractionDigits: 0 }).format(grandTotal)}</span>
                      </div>

                      <Button 
                        className="w-full h-14 rounded-xl font-bold text-lg shadow-xl shadow-primary/20" 
                        onClick={handlePlaceOrder}
                        disabled={processing}
                      >
                         {processing ? "Processing..." : "Place Order"} <ArrowRight className="w-5 h-5 ml-2" />
                      </Button>
                   </CardContent>
                </Card>
            </div>
        </div>
    </div>
  );
}
