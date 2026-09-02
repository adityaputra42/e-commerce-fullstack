'use client';

import { useState } from 'react';
import Link from 'next/link';
import { toast } from 'sonner';
import { ShoppingBag, Facebook, Instagram } from 'lucide-react';

const Footer = () => {
  const [email, setEmail] = useState('');

  const handleSubscribe = (e: React.FormEvent) => {
    e.preventDefault();
    toast.info("Newsletter signup isn't connected yet — coming soon.");
    setEmail('');
  };

  return (
    <footer className="bg-neutral-50 border-t border-neutral-200 pt-16 pb-8">
      <div className="container mx-auto px-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-12 mb-14">
          <div className="space-y-5">
            <Link href="/" className="flex items-center gap-2.5">
              <div className="w-9 h-9 bg-primary rounded-md flex items-center justify-center text-white">
                <ShoppingBag className="w-4.5 h-4.5" />
              </div>
              <span className="text-base font-black tracking-tight text-neutral-900">
                Barrakallah <span className="text-primary">Hijab Store</span>
              </span>
            </Link>
            <p className="text-sm text-neutral-500 font-medium leading-relaxed max-w-xs">
              Modest fashion made with quality materials and honest pricing. Thoughtfully
              curated hijabs and modest wear for everyday life.
            </p>
            <div className="flex gap-3">
              <a href="#" className="w-9 h-9 bg-white border border-neutral-200 rounded-md flex items-center justify-center text-neutral-400 hover:text-primary hover:border-neutral-300 transition-colors"><Facebook className="w-4 h-4" /></a>
              <a href="#" className="w-9 h-9 bg-white border border-neutral-200 rounded-md flex items-center justify-center text-neutral-400 hover:text-primary hover:border-neutral-300 transition-colors"><Instagram className="w-4 h-4" /></a>
            </div>
          </div>

          <div>
            <h4 className="text-sm font-bold text-neutral-900 mb-5">Quick Links</h4>
            <ul className="space-y-3">
              <li><Link href="/shop" className="text-sm text-neutral-500 hover:text-primary font-medium transition-colors">All Products</Link></li>
              <li><Link href="/collections" className="text-sm text-neutral-500 hover:text-primary font-medium transition-colors">Collections</Link></li>
              <li><Link href="/cart" className="text-sm text-neutral-500 hover:text-primary font-medium transition-colors">Shopping Cart</Link></li>
              <li><Link href="/profile" className="text-sm text-neutral-500 hover:text-primary font-medium transition-colors">My Account</Link></li>
            </ul>
          </div>

          <div>
            <h4 className="text-sm font-bold text-neutral-900 mb-5">Policies</h4>
            <ul className="space-y-3">
              <li><Link href="/privacy" className="text-sm text-neutral-500 hover:text-primary font-medium transition-colors">Privacy Policy</Link></li>
              <li><Link href="/terms" className="text-sm text-neutral-500 hover:text-primary font-medium transition-colors">Terms of Service</Link></li>
              <li><Link href="/shipping" className="text-sm text-neutral-500 hover:text-primary font-medium transition-colors">Shipping Info</Link></li>
              <li><Link href="/returns" className="text-sm text-neutral-500 hover:text-primary font-medium transition-colors">Refunds & Returns</Link></li>
            </ul>
          </div>

          <div>
            <h4 className="text-sm font-bold text-neutral-900 mb-5">Newsletter</h4>
            <p className="text-sm text-neutral-500 font-medium mb-4">Get notified about new arrivals and product updates.</p>
            <form onSubmit={handleSubscribe} className="space-y-2.5">
              <input
                type="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Your email address"
                className="w-full h-11 px-4 bg-white border border-neutral-300 rounded-md text-sm focus:border-primary outline-none transition-all"
              />
              <button type="submit" className="w-full h-11 bg-neutral-900 text-white rounded-md text-sm font-bold hover:bg-black transition-colors">
                Subscribe
              </button>
            </form>
          </div>
        </div>

        <div className="pt-8 border-t border-neutral-200 flex flex-col md:flex-row items-center justify-between gap-4">
          <p className="text-xs text-neutral-400 font-medium">
            © {new Date().getFullYear()} Barrakallah Hijab Store. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
