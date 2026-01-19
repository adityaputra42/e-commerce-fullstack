"use client";

import React from 'react';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { User, MapPin, ShoppingBag, Lock, LogOut } from 'lucide-react';
import { useAuth } from '@/context/AuthContext';

export default function ProfileLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const { logout } = useAuth();

  const menuItems = [
    {
      label: "My Profile",
      icon: User,
      href: "/profile",
      exact: true
    },
    {
      label: "My Addresses",
      icon: MapPin,
      href: "/profile/addresses"
    },
    {
      label: "My Orders",
      icon: ShoppingBag,
      href: "/profile/transactions"
    },
    {
      label: "Change Password",
      icon: Lock,
      href: "/profile/security"
    }
  ];

  return (
    <div className="container mx-auto px-4 py-24 min-h-screen">
      <div className="flex flex-col md:flex-row gap-8">
        {/* Sidebar */}
        <aside className="w-full md:w-64 flex-shrink-0">
          <div className="bg-white rounded-2xl shadow-sm border border-slate-100 p-4 sticky top-24">
            <div className="space-y-1">
              {menuItems.map((item) => {
                const isActive = item.exact 
                  ? pathname === item.href 
                  : pathname.startsWith(item.href);
                
                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    className={`flex items-center gap-3 px-4 py-3 rounded-xl transition-colors font-medium ${
                      isActive 
                        ? 'bg-teal-50 text-teal-600' 
                        : 'text-slate-600 hover:bg-slate-50 hover:text-slate-900'
                    }`}
                  >
                    <item.icon className="w-5 h-5" />
                    {item.label}
                  </Link>
                );
              })}
              
              <hr className="my-2 border-slate-100" />
              
              <button
                onClick={logout}
                className="w-full flex items-center gap-3 px-4 py-3 rounded-xl transition-colors font-medium text-red-600 hover:bg-red-50"
              >
                <LogOut className="w-5 h-5" />
                Logout
              </button>
            </div>
          </div>
        </aside>

        {/* Content */}
        <div className="flex-1">
          {children}
        </div>
      </div>
    </div>
  );
}
