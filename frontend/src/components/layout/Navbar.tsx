'use client';

import Link from 'next/link';
import { ShoppingCart, User, Search, Menu, X } from 'lucide-react';
import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

const Navbar = () => {
  const [isScrolled, setIsScrolled] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 20);
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <nav className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
      isScrolled ? 'glass-nav h-16 shadow-lg shadow-indigo-500/5' : 'bg-transparent h-20'
    }`}>
      <div className="container mx-auto px-6 h-full flex items-center justify-between">
        <Link href="/" className="flex items-center gap-2 group">
          <div className="w-10 h-10 bg-indigo-600 rounded-xl flex items-center justify-center text-white scale-100 group-hover:scale-110 transition-transform shadow-lg shadow-indigo-600/30">
            <ShoppingCart className="w-5 h-5 fill-current" />
          </div>
          <span className="text-xl font-black tracking-tight text-slate-900 dark:text-white">STORE<span className="text-indigo-600">X</span></span>
        </Link>

        {/* Desktop Menu */}
        <div className="hidden md:flex items-center gap-8">
          <Link href="/" className="text-sm font-bold text-slate-600 hover:text-indigo-600 transition-colors">Home</Link>
          <Link href="/shop" className="text-sm font-bold text-slate-600 hover:text-indigo-600 transition-colors">Shop</Link>
          <Link href="/collections" className="text-sm font-bold text-slate-600 hover:text-indigo-600 transition-colors">Collections</Link>
        </div>

        <div className="flex items-center gap-4">
          <button className="p-2 text-slate-600 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-full transition-colors">
            <Search className="w-5 h-5" />
          </button>
          <Link href="/cart" className="relative p-2 text-slate-600 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-full transition-colors">
            <ShoppingCart className="w-5 h-5" />
            <span className="absolute top-0 right-0 w-4 h-4 bg-indigo-600 text-[10px] font-bold text-white flex items-center justify-center rounded-full">0</span>
          </Link>
          <Link href="/login" className="hidden md:flex items-center gap-2 py-2 px-4 bg-slate-900 text-white dark:bg-white dark:text-slate-900 rounded-full text-xs font-bold hover:scale-105 active:scale-95 transition-all">
            <User className="w-4 h-4" />
            Login
          </Link>
          
          <button 
            className="md:hidden p-2 text-slate-600"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          >
            {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>
      </div>

      {/* Mobile Menu */}
      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div 
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            className="md:hidden glass-nav border-b border-slate-100 overflow-hidden"
          >
            <div className="container mx-auto px-6 py-6 flex flex-col gap-4">
              <Link href="/" onClick={() => setIsMobileMenuOpen(false)} className="text-lg font-bold text-slate-900">Home</Link>
              <Link href="/shop" onClick={() => setIsMobileMenuOpen(false)} className="text-lg font-bold text-slate-900">Shop</Link>
              <Link href="/collections" onClick={() => setIsMobileMenuOpen(false)} className="text-lg font-bold text-slate-900">Collections</Link>
              <hr className="border-slate-100" />
              <Link href="/login" onClick={() => setIsMobileMenuOpen(false)} className="flex items-center gap-2 font-bold text-indigo-600">
                <User className="w-5 h-5" />
                Login / Register
              </Link>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </nav>
  );
};

export default Navbar;
