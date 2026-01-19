"use client";

import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { useAuth } from '@/context/AuthContext';
import { authService } from '@/services/api';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Card, CardContent, CardHeader } from '@/components/common/Card';
import { toast } from 'sonner';

const profileSchema = z.object({
  first_name: z.string().min(1, "First name is required"),
  last_name: z.string().min(1, "Last name is required"),
  username: z.string().min(3, "Username must be at least 3 characters"),
  email: z.string().email("Invalid email address"),
});

type ProfileFormValues = z.infer<typeof profileSchema>;

export default function ProfilePage() {
  const { user } = useAuth();
  const [loading, setLoading] = useState(false);

  const form = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      first_name: user?.first_name || '',
      last_name: user?.last_name || '',
      username: user?.username || '',
      email: user?.email || '',
    },
  });

  async function onSubmit(values: ProfileFormValues) {
    setLoading(true);
    try {
      await authService.updateProfile(values);
      toast.success("Profile updated successfully");
      // Could also refresh user context here if needed, but page reload will do for now
      // or authService.getCurrentUser() call in auth context would catch it on next mount
      // Ideally we should update local state too.
      // But for now let's show success.
      window.location.reload(); // Simple way to refresh data
    } catch (error: any) {
      console.error(error);
      toast.error(error.response?.data?.message || "Failed to update profile");
    } finally {
      setLoading(false);
    }
  }

  return (
    <Card className="shadow-sm border-slate-100">
      <CardHeader className="border-b border-slate-100">
        <h2 className="text-xl font-bold text-slate-900">Personal Information</h2>
        <p className="text-sm text-slate-500">Manage your personal details</p>
      </CardHeader>
      <CardContent className="pt-6">
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Input
              label="First Name"
              error={form.formState.errors.first_name?.message}
              {...form.register('first_name')}
            />
            <Input
              label="Last Name"
              error={form.formState.errors.last_name?.message}
              {...form.register('last_name')}
            />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Input
              label="Username"
              error={form.formState.errors.username?.message}
              {...form.register('username')}
            />
            <Input
              label="Email"
              type="email"
              error={form.formState.errors.email?.message}
              {...form.register('email')}
            />
          </div>

          <div className="flex justify-end pt-4">
            <Button type="submit" disabled={loading} className="min-w-[120px]">
              {loading ? "Saving..." : "Save Changes"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
