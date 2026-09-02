'use client';

import Link from 'next/link';
import { ShoppingCart, User, Search, Menu, X, LogOut, UserCircle, Truck } from 'lucide-react';
import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';

import { useCart } from '@/context/CartContext';
import { useAuth } from '@/context/AuthContext';

const Navbar = () => {
  const { cartCount } = useCart();
  const { isAuthenticated, user, logout } = useAuth();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isUserMenuOpen, setIsUserMenuOpen] = useState(false);

  return (
    <div className="fixed top-0 left-0 right-0 z-50">
      <div className="hidden sm:flex items-center justify-center gap-2 h-9 bg-neutral-900 text-neutral-100 text-xs font-medium tracking-wide">
        <Truck className="w-3.5 h-3.5" />
        <span>Free shipping on orders over $50</span>
      </div>

      <nav className="bg-white border-b border-neutral-200 h-16">
        <div className="container mx-auto px-6 h-full flex items-center justify-between">
          <Link href="/" className="flex items-center gap-2.5">
            <div className="w-9 h-9 bg-primary rounded-md flex items-center justify-center text-white">
              <ShoppingCart className="w-4.5 h-4.5" />
            </div>
            <span className="text-lg font-black tracking-tight text-neutral-900">
              Barrakallah <span className="text-primary">Hijab Store</span>
            </span>
          </Link>

          {/* Desktop Menu */}
          <div className="hidden md:flex items-center gap-8">
            <Link href="/" className="text-sm font-semibold text-neutral-600 hover:text-neutral-900 transition-colors">Home</Link>
            <Link href="/shop" className="text-sm font-semibold text-neutral-600 hover:text-neutral-900 transition-colors">Shop</Link>
            <Link href="/collections" className="text-sm font-semibold text-neutral-600 hover:text-neutral-900 transition-colors">Collections</Link>
          </div>

          <div className="flex items-center gap-1">
            <button className="p-2.5 text-neutral-600 hover:bg-neutral-100 rounded-md transition-colors">
              <Search className="w-4.5 h-4.5" />
            </button>
            <Link href="/cart" className="relative p-2.5 text-neutral-600 hover:bg-neutral-100 rounded-md transition-colors">
              <ShoppingCart className="w-4.5 h-4.5" />
              {cartCount > 0 && (
                <span className="absolute top-1 right-1 min-w-4 h-4 px-1 bg-primary text-[10px] font-bold text-white flex items-center justify-center rounded-full">
                  {cartCount}
                </span>
              )}
            </Link>

            {/* User Profile or Login Button */}
            {isAuthenticated ? (
              <div className="relative hidden md:block ml-1">
                <button
                  onClick={() => setIsUserMenuOpen(!isUserMenuOpen)}
                  className="flex items-center gap-2 p-2.5 text-neutral-600 hover:bg-neutral-100 rounded-md transition-colors"
                >
                  <UserCircle className="w-5.5 h-5.5" />
                </button>

                <AnimatePresence>
                  {isUserMenuOpen && (
                    <>
                      <div
                        className="fixed inset-0 z-40"
                        onClick={() => setIsUserMenuOpen(false)}
                      />

                      <motion.div
                        initial={{ opacity: 0, y: -8 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -8 }}
                        transition={{ duration: 0.15 }}
                        className="absolute right-0 mt-2 w-56 bg-white rounded-md shadow-lg border border-neutral-200 overflow-hidden z-50"
                      >
                        <div className="p-4 border-b border-neutral-100">
                          <p className="font-bold text-neutral-900 text-sm">{user?.first_name} {user?.last_name}</p>
                          <p className="text-xs text-neutral-500">{user?.email}</p>
                        </div>

                        <div className="py-1.5">
                          <Link
                            href="/profile"
                            onClick={() => setIsUserMenuOpen(false)}
                            className="flex items-center gap-3 px-4 py-2 text-sm text-neutral-700 hover:bg-neutral-50 transition-colors"
                          >
                            <User className="w-4 h-4" />
                            Profile
                          </Link>
                          <Link
                            href="/profile/transactions"
                            onClick={() => setIsUserMenuOpen(false)}
                            className="flex items-center gap-3 px-4 py-2 text-sm text-neutral-700 hover:bg-neutral-50 transition-colors"
                          >
                            <ShoppingCart className="w-4 h-4" />
                            My Orders
                          </Link>
                        </div>

                        <div className="border-t border-neutral-100 py-1.5">
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
              <Link href="/login" className="hidden md:flex items-center gap-2 h-9 ml-2 px-4 bg-neutral-900 text-white rounded-md text-xs font-bold hover:bg-black active:scale-[0.98] transition-all">
                <User className="w-3.5 h-3.5" />
                Login
              </Link>
            )}

            <button
              className="md:hidden p-2.5 text-neutral-600"
              onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
            >
              {isMobileMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
            </button>
          </div>
        </div>
      </nav>

      {/* Mobile Menu */}
      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: 'auto' }}
            exit={{ opacity: 0, height: 0 }}
            className="md:hidden bg-white border-b border-neutral-200 overflow-hidden"
          >
            <div className="container mx-auto px-6 py-6 flex flex-col gap-4">
              <Link href="/" onClick={() => setIsMobileMenuOpen(false)} className="text-base font-bold text-neutral-900">Home</Link>
              <Link href="/shop" onClick={() => setIsMobileMenuOpen(false)} className="text-base font-bold text-neutral-900">Shop</Link>
              <Link href="/collections" onClick={() => setIsMobileMenuOpen(false)} className="text-base font-bold text-neutral-900">Collections</Link>
              <hr className="border-neutral-100" />
              {isAuthenticated ? (
                <>
                  <Link href="/profile" onClick={() => setIsMobileMenuOpen(false)} className="flex items-center gap-2 font-bold text-primary">
                    <UserCircle className="w-5 h-5" />
                    Profile
                  </Link>
                  <Link href="/profile/transactions" onClick={() => setIsMobileMenuOpen(false)} className="flex items-center gap-2 font-bold text-neutral-700">
                    <ShoppingCart className="w-5 h-5" />
                    My Orders
                  </Link>
                  <button onClick={() => { logout(); setIsMobileMenuOpen(false); }} className="flex items-center gap-2 font-bold text-red-600">
                    <LogOut className="w-5 h-5" />
                    Logout
                  </button>
                </>
              ) : (
                <Link href="/login" onClick={() => setIsMobileMenuOpen(false)} className="flex items-center gap-2 font-bold text-primary">
                  <User className="w-5 h-5" />
                  Login / Register
                </Link>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default Navbar;
