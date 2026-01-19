"use client";

import { useState } from 'react';
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
  name: z.string().min(2),
  email: z.string().email(),
  password: z.string().min(6),
});

export default function RegisterPage() {
  const { login } = useAuth();
  const [loading, setLoading] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      email: "",
      password: "",
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {
    setLoading(true);
    try {
      const response = await authService.register(values);
      login(response.token, response.user);
      toast.success("Account created successfully");
    } catch (error: any) {
      console.error(error);
      toast.error(error.response?.data?.message || "Failed to create account.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="flex min-h-[80vh] items-center justify-center container mx-auto px-4">
      <Card className="w-full max-w-md shadow-xl">
        <CardContent className="pt-8">
          <div className="space-y-2 mb-6">
            <h1 className="text-3xl font-black text-slate-900">Create an Account</h1>
            <p className="text-slate-500">
              Enter your details below to create your account
            </p>
          </div>
          
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <Input
                label="Name"
                placeholder="John Doe"
                error={form.formState.errors.name?.message}
                {...form.register('name')}
              />
            </div>
            
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
              {loading ? "Creating account..." : "Register"}
            </Button>
          </form>
          
          <div className="mt-6 text-center text-sm text-slate-600">
            Already have an account?{" "}
            <Link href="/login" className="text-teal-500 hover:text-teal-600 font-bold">
              Login
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
