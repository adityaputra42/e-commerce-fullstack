'use client';

import Link from 'next/link';
import { ShoppingCart, User, Search, Menu, X, LogOut, UserCircle } from 'lucide-react';
import { useState, useEffect } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

import { useCart } from '@/context/CartContext';
import { useAuth } from '@/context/AuthContext';

const Navbar = () => {
  const { cartCount } = useCart();
  const { isAuthenticated, user, logout } = useAuth();
  const [isScrolled, setIsScrolled] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 20);
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  return (
    <nav className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
      isScrolled ? 'glass-nav h-16 shadow-lg shadow-teal-500/5' : 'bg-transparent h-20'
    }`}>
      <div className="container mx-auto px-6 h-full flex items-center justify-between">
        <Link href="/" className="flex items-center gap-2 group">
          <div className="w-10 h-10 bg-teal-500 rounded-xl flex items-center justify-center text-white scale-100 group-hover:scale-110 transition-transform shadow-lg shadow-teal-600/30">
            <ShoppingCart className="w-5 h-5 fill-current" />
          </div>
          <span className="text-xl font-black tracking-tight text-slate-900">STORE<span className="text-teal-500">X</span></span>
        </Link>

        {/* Desktop Menu */}
        <div className="hidden md:flex items-center gap-8">
          <Link href="/" className="text-sm font-bold text-slate-600 hover:text-teal-500 transition-colors">Home</Link>
          <Link href="/shop" className="text-sm font-bold text-slate-600 hover:text-teal-500 transition-colors">Shop</Link>
          <Link href="/collections" className="text-sm font-bold text-slate-600 hover:text-teal-500 transition-colors">Collections</Link>
        </div>

        <div className="flex items-center gap-4">
          <button className="p-2 text-slate-600 hover:bg-slate-100 rounded-full transition-colors">
            <Search className="w-5 h-5" />
          </button>
          <Link href="/cart" className="relative p-2 text-slate-600 hover:bg-slate-100 rounded-full transition-colors">
            <ShoppingCart className="w-5 h-5" />
            {cartCount > 0 && (
              <span className="absolute top-0 right-0 w-4 h-4 bg-teal-500 text-[10px] font-bold text-white flex items-center justify-center rounded-full">
                {cartCount}
              </span>
            )}
          </Link>
          
          {/* User Profile or Login Button */}
          {isAuthenticated ? (
            <div className="relative hidden md:block">
              <button 
                onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                className="flex items-center gap-2 p-2 text-slate-600 hover:bg-slate-100 rounded-full transition-colors"
              >
                <UserCircle className="w-6 h-6" />
              </button>
              
              <AnimatePresence>
                {isUserMenuOpen && (
                  <>
                    {/* Backdrop */}
                    <div 
                      className="fixed inset-0 z-40" 
                      onClick={() => setIsUserMenuOpen(false)}
                    />
                    
                    {/* Dropdown Menu */}
                    <motion.div
                      initial={{ opacity: 0, y: -10 }}
                      animate={{ opacity: 1, y: 0 }}
                      exit={{ opacity: 0, y: -10 }}
                      transition={{ duration: 0.2 }}
                      className="absolute right-0 mt-2 w-56 bg-white rounded-2xl shadow-xl border border-slate-200 overflow-hidden z-50"
                    >
                      <div className="p-4 border-b border-slate-100">
                        <p className="font-bold text-slate-900">{user?.first_name} {user?.last_name}</p>
                        <p className="text-xs text-slate-500">{user?.email}</p>
                      </div>
                      
                      <div className="py-2">
                        <Link 
                          href="/profile" 
                          onClick={() => setIsUserMenuOpen(false)}
                          className="flex items-center gap-3 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 transition-colors"
                        >
                          <User className="w-4 h-4" />
                          Profile
                        </Link>
                        <Link 
                          href="/profile/transactions" 
                          onClick={() => setIsUserMenuOpen(false)}
                          className="flex items-center gap-3 px-4 py-2 text-sm text-slate-700 hover:bg-slate-50 transition-colors"
                        >
                          <ShoppingCart className="w-4 h-4" />
                          My Orders
                        </Link>
                      </div>
                      
                      <div className="border-t border-slate-100 py-2">
                        <button 
                          onClick={() => { logout(); setIsUserMenuOpen(false); }}
                          className="flex items-center gap-3 px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors w-full"
                        >
                          <LogOut className="w-4 h-4" />
                          Logout
                        </button>
                      </div>
                    </motion.div>
                  </>
                )}
              </AnimatePresence>
            </div>
          ) : (
            <Link href="/login" className="hidden md:flex items-center gap-2 py-2 px-4 bg-slate-900 text-white rounded-full text-xs font-bold hover:scale-105 active:scale-95 transition-all">
              <User className="w-4 h-4" />
              Login
            </Link>
          )}
          
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
              {isAuthenticated ? (
                <>
                  <Link href="/profile" onClick={() => setIsMobileMenuOpen(false)} className="flex items-center gap-2 font-bold text-teal-500">
                    <UserCircle className="w-5 h-5" />
                    Profile
                  </Link>
                  <Link href="/profile/transactions" onClick={() => setIsMobileMenuOpen(false)} className="flex items-center gap-2 font-bold text-slate-700">
                    <ShoppingCart className="w-5 h-5" />
                    My Orders
                  </Link>
                  <button onClick={() => { logout(); setIsMobileMenuOpen(false); }} className="flex items-center gap-2 font-bold text-red-600">
                    <LogOut className="w-5 h-5" />
                    Logout
                  </button>
                </>
              ) : (
                <Link href="/login" onClick={() => setIsMobileMenuOpen(false)} className="flex items-center gap-2 font-bold text-teal-500">
                  <User className="w-5 h-5" />
                  Login / Register
                </Link>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </nav>
  );
};

export default Navbar;
