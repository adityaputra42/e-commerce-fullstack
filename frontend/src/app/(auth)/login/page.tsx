"use client";

import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { useAuth } from '@/context/AuthContext';
import { authService } from '@/services/api';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Card, CardContent } from '@/components/common/Card';
import { toast } from 'sonner';
import Link from 'next/link';

const formSchema = z.object({
  email: z.string().email(),
  password: z.string().min(6),
});

export default function LoginPage() {
  const { login } = useAuth();
  const [loading, setLoading] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    setLoading(true);
    try {
      const response = await authService.login(values.email, values.password);
      login(response.token, response.user);
      toast.success("Logged in successfully");
    } catch (error: any) {
      console.error(error);
      toast.error(error.response?.data?.message || "Failed to login. Please check your credentials.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-[80vh] items-center justify-center container mx-auto px-4">
      <Card className="w-full max-w-md shadow-xl">
        <CardContent className="pt-8">
          <div className="space-y-2 mb-6">
            <h1 className="text-3xl font-black text-slate-900">Login</h1>
            <p className="text-slate-500">
              Enter your email and password to access your account
            </p>
          </div>
          
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <Input
                label="Email"
                type="email"
                placeholder="email@example.com"
                error={form.formState.errors.email?.message}
                {...form.register('email')}
              />
            </div>
            
            <div>
              <Input
                label="Password"
                type="password"
                error={form.formState.errors.password?.message}
                {...form.register('password')}
              />
            </div>
            
            <Button type="submit" className="w-full" disabled={loading}>
              {loading ? "Logging in..." : "Login"}
            </Button>
          </form>
          
          <div className="mt-6 text-center text-sm text-slate-600">
            Don&apos;t have an account?{" "}
            <Link href="/register" className="text-teal-500 hover:text-teal-600 font-bold">
              Register
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
