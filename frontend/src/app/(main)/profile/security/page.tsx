"use client";

import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { authService } from '@/services/api';
import { Button } from '@/components/common/Button';
import { Input } from '@/components/common/Input';
import { Card, CardContent, CardHeader } from '@/components/common/Card';
import { toast } from 'sonner';

const passwordSchema = z.object({
  current_password: z.string().min(1, "Current password is required"),
  new_password: z.string().min(8, "Password must be at least 8 characters"),
  confirm_password: z.string().min(8, "Password must be at least 8 characters"),
}).refine((data) => data.new_password === data.confirm_password, {
  message: "Passwords do not match",
  path: ["confirm_password"],
});

type PasswordFormValues = z.infer<typeof passwordSchema>;

export default function ChangePasswordPage() {
  const [loading, setLoading] = useState(false);

  const form = useForm<PasswordFormValues>({
    resolver: zodResolver(passwordSchema),
    defaultValues: {
      current_password: '',
      new_password: '',
      confirm_password: '',
    },
  });

  async function onSubmit(values: PasswordFormValues) {
    setLoading(true);
    try {
      await authService.changePassword(values);
      toast.success("Password changed successfully");
      form.reset();
    } catch (error: any) {
      console.error(error);
      toast.error(error.response?.data?.message || "Failed to change password");
    } finally {
      setLoading(false);
    }
  }

  return (
    <Card className="shadow-sm border-slate-100">
      <CardHeader className="border-b border-slate-100">
        <h2 className="text-xl font-bold text-slate-900">Change Password</h2>
        <p className="text-sm text-slate-500">Ensure your account is secure</p>
      </CardHeader>
      <CardContent className="pt-6">
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 max-w-md">
          <Input
            label="Current Password"
            type="password"
            error={form.formState.errors.current_password?.message}
            {...form.register('current_password')}
          />

          <Input
            label="New Password"
            type="password"
            error={form.formState.errors.new_password?.message}
            {...form.register('new_password')}
          />

          <Input
            label="Confirm New Password"
            type="password"
            error={form.formState.errors.confirm_password?.message}
            {...form.register('confirm_password')}
          />

          <div className="flex justify-end pt-4">
            <Button type="submit" disabled={loading} className="min-w-[120px]">
              {loading ? "Saving..." : "Update Password"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
