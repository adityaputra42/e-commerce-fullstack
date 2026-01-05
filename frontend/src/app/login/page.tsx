'use client';

import { useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import Link from 'next/link';
import { Mail, Lock, LogIn, Sparkles, ShoppingBag, ArrowLeft } from 'lucide-react';
import { motion } from 'framer-motion';
import api from '@/lib/api';

export default function LoginPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const redirect = searchParams.get('redirect') || '/';
  
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      // Assuming backend login endpoint
      const response = await api.post('/auth/login', { email, password });
      const { access_token } = response.data.data;
      
      localStorage.setItem('token', access_token);
      router.push(redirect);
    } catch (err: any) {
      setError(err.response?.data?.message || 'Invalid credentials. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-slate-50 relative overflow-hidden p-6">
      {/* Decorative Orbs */}
      <div className="absolute top-0 right-0 w-[500px] h-[500px] bg-indigo-500/10 rounded-full blur-3xl -translate-y-1/2 translate-x-1/2" />
      <div className="absolute bottom-0 left-0 w-[400px] h-[400px] bg-rose-500/10 rounded-full blur-3xl translate-y-1/2 -translate-x-1/2" />

      <motion.div 
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        className="w-full max-w-md"
      >
        <Link href="/" className="inline-flex items-center gap-2 text-[10px] font-black uppercase tracking-widest text-slate-400 hover:text-primary transition-colors mb-8 group">
           <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
           Back to Store
        </Link>

        <div className="bg-white rounded-[3rem] p-10 md:p-12 shadow-2xl shadow-indigo-500/5 relative overflow-hidden">
           <div className="absolute top-0 right-0 p-8 opacity-5">
              <ShoppingBag className="w-32 h-32 rotate-12" />
           </div>

           <div className="space-y-8 relative z-10">
              <div className="space-y-2">
                 <div className="w-16 h-1 bg-primary" />
                 <h1 className="text-4xl font-black text-slate-900 tracking-tighter italic uppercase">Welcome Back.</h1>
                 <p className="text-sm text-slate-400 font-bold uppercase tracking-widest">Sign in to your account</p>
              </div>

              {error && (
                <div className="p-4 bg-rose-50 text-rose-600 rounded-2xl text-xs font-bold border border-rose-100 animate-shake">
                   {error}
                </div>
              )}

              <form onSubmit={handleLogin} className="space-y-6">
                 <div className="space-y-2">
                    <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Email Address</label>
                    <div className="relative group">
                       <Mail className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-300 group-focus-within:text-primary transition-colors" />
                       <input 
                         type="email" 
                         required
                         value={email}
                         onChange={(e) => setEmail(e.target.value)}
                         placeholder="your@email.com" 
                         className="w-full h-16 pl-12 pr-4 bg-slate-50 border-2 border-transparent focus:bg-white focus:border-indigo-600 rounded-2xl text-sm font-bold outline-none transition-all"
                       />
                    </div>
                 </div>

                 <div className="space-y-2">
                    <div className="flex justify-between items-center ml-1">
                       <label className="text-[10px] font-black text-slate-400 uppercase tracking-widest">Password</label>
                       <a href="#" className="text-[10px] font-black text-primary uppercase tracking-widest hover:underline">Forgot?</a>
                    </div>
                    <div className="relative group">
                       <Lock className="absolute left-4 top-1/2 -translate-y-1/2 w-5 h-5 text-slate-300 group-focus-within:text-primary transition-colors" />
                       <input 
                         type="password" 
                         required
                         value={password}
                         onChange={(e) => setPassword(e.target.value)}
                         placeholder="••••••••" 
                         className="w-full h-16 pl-12 pr-4 bg-slate-50 border-2 border-transparent focus:bg-white focus:border-indigo-600 rounded-2xl text-sm font-bold outline-none transition-all"
                       />
                    </div>
                 </div>

                 <button 
                  disabled={isLoading}
                  className="w-full h-16 bg-slate-900 text-white rounded-2xl flex items-center justify-center gap-3 text-sm font-black uppercase tracking-widest hover:bg-black transition-all shadow-xl active:scale-95 disabled:opacity-50"
                 >
                    {isLoading ? (
                      <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                    ) : (
                      <>
                        <LogIn className="w-5 h-5" />
                        Enter Store
                      </>
                    )}
                 </button>
              </form>

              <div className="pt-6 border-t border-slate-50 flex flex-col items-center gap-4 text-center">
                 <p className="text-xs text-slate-400 font-bold uppercase tracking-widest">Don't have an account?</p>
                 <Link href="/register" className="inline-flex items-center gap-2 text-primary font-black uppercase tracking-widest text-[10px] hover:underline">
                    <Sparkles className="w-4 h-4" />
                    Create New Account
                 </Link>
              </div>
           </div>
        </div>
      </motion.div>
    </div>
  );
}
