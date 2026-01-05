import Link from 'next/link';
import { ShoppingBag, Facebook, Twitter, Instagram, Github } from 'lucide-react';

const Footer = () => {
  return (
    <footer className="bg-slate-50 dark:bg-slate-900 border-t border-slate-100 dark:border-slate-800 pt-20 pb-10">
      <div className="container mx-auto px-6">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-12 mb-16">
          <div className="space-y-6">
            <Link href="/" className="flex items-center gap-2 group">
              <div className="w-10 h-10 bg-primary rounded-xl flex items-center justify-center text-white shadow-lg shadow-indigo-600/30 transition-transform group-hover:rotate-12">
                <ShoppingBag className="w-5 h-5 fill-current" />
              </div>
              <span className="text-xl font-black tracking-tight text-slate-900 dark:text-white uppercase italic">Store<span className="text-primary">X</span></span>
            </Link>
            <p className="text-sm text-slate-500 font-medium leading-relaxed">
              Elevating your lifestyle with premium products curated for quality and modern aesthetics. Join our community and discover the difference.
            </p>
            <div className="flex gap-4">
              <a href="#" className="p-2 bg-white dark:bg-slate-800 rounded-xl text-slate-400 hover:text-primary shadow-sm transition-all hover:scale-110"><Facebook className="w-5 h-5" /></a>
              <a href="#" className="p-2 bg-white dark:bg-slate-800 rounded-xl text-slate-400 hover:text-primary shadow-sm transition-all hover:scale-110"><Twitter className="w-5 h-5" /></a>
              <a href="#" className="p-2 bg-white dark:bg-slate-800 rounded-xl text-slate-400 hover:text-primary shadow-sm transition-all hover:scale-110"><Instagram className="w-5 h-5" /></a>
              <a href="#" className="p-2 bg-white dark:bg-slate-800 rounded-xl text-slate-400 hover:text-primary shadow-sm transition-all hover:scale-110"><Github className="w-5 h-5" /></a>
            </div>
          </div>

          <div>
            <h4 className="text-sm font-black text-slate-900 dark:text-white uppercase tracking-widest mb-6">Quick Links</h4>
            <ul className="space-y-4">
              <li><Link href="/shop" className="text-sm text-slate-500 hover:text-primary font-medium transition-colors">All Products</Link></li>
              <li><Link href="/collections" className="text-sm text-slate-500 hover:text-primary font-medium transition-colors">New Arrivals</Link></li>
              <li><Link href="/cart" className="text-sm text-slate-500 hover:text-primary font-medium transition-colors">Shopping Cart</Link></li>
              <li><Link href="/account" className="text-sm text-slate-500 hover:text-primary font-medium transition-colors">User Account</Link></li>
            </ul>
          </div>

          <div>
            <h4 className="text-sm font-black text-slate-900 dark:text-white uppercase tracking-widest mb-6">Policies</h4>
            <ul className="space-y-4">
              <li><Link href="/privacy" className="text-sm text-slate-500 hover:text-primary font-medium transition-colors">Privacy Policy</Link></li>
              <li><Link href="/terms" className="text-sm text-slate-500 hover:text-primary font-medium transition-colors">Terms of Service</Link></li>
              <li><Link href="/shipping" className="text-sm text-slate-500 hover:text-primary font-medium transition-colors">Shipping Info</Link></li>
              <li><Link href="/returns" className="text-sm text-slate-500 hover:text-primary font-medium transition-colors">Refunds & Returns</Link></li>
            </ul>
          </div>

          <div>
            <h4 className="text-sm font-black text-slate-900 dark:text-white uppercase tracking-widest mb-6">Newsletter</h4>
            <p className="text-sm text-slate-500 font-medium mb-4">Subscribe to get special offers and limited stock alerts.</p>
            <form className="space-y-3">
              <input 
                type="email" 
                placeholder="Your email address" 
                className="w-full h-12 px-4 bg-white dark:bg-slate-800 border-2 border-slate-100 dark:border-slate-700 rounded-2xl text-sm focus:border-indigo-500 outline-none transition-all" 
              />
              <button className="w-full h-12 bg-primary text-white rounded-2xl text-sm font-black uppercase tracking-widest hover:bg-indigo-700 transition-colors shadow-lg shadow-indigo-600/20">Subscribe</button>
            </form>
          </div>
        </div>

        <div className="pt-10 border-t border-slate-100 dark:border-slate-800 flex flex-col md:flex-row items-center justify-between gap-4">
          <p className="text-xs text-slate-400 font-bold uppercase tracking-widest">
            © {new Date().getFullYear()} STOREX COMMERCE. ALL RIGHTS RESERVED.
          </p>
          <div className="flex gap-8">
            <span className="text-[10px] text-slate-300 font-black uppercase italic tracking-tighter"></span>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
