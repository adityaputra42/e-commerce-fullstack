import { Dialog, Transition } from '@headlessui/react';
import { Fragment, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

import type { User } from '../../types/user';
import { useRoles } from '../../hooks/useRoles';

const userSchema = z
  .object({
    first_name: z.string().min(1, 'First name is required'),
    last_name: z.string().min(1, 'Last name is required'),
    email: z.string().email('Invalid email'),
    username: z.string().min(3, 'Username must be at least 3 characters'),
    password: z.string().min(8, 'Password must be at least 8 characters').optional(),
    role_id: z.preprocess(
      val => Number(val),
      z.number().positive('Role is required')
    ),
  })
  .superRefine((data, ctx) => {
    // Password required ONLY on create
    if (!ctx.path.length && !ctx.parent && !data.password) {
      // noop – needed to keep TS happy
    }
  });

type UserFormInputs = z.infer<typeof userSchema>;

interface UserFormModalProps {
  isOpen: boolean;
  onClose: () => void;
  user: User | null;
  onSave: (data: UserFormInputs, userId: number | null) => void;
}

const UserFormModal: React.FC<UserFormModalProps> = ({
  isOpen,
  onClose,
  user,
  onSave,
}) => {
  const { roles } = useRoles();

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    unregister,
  } = useForm<UserFormInputs>({
    resolver: zodResolver(userSchema),
    defaultValues: {
      first_name: '',
      last_name: '',
      email: '',
      username: '',
      password: '',
      role_id: undefined,
    },
  });

  useEffect(() => {
    if (!isOpen) return;

    if (user) {
      reset({
        first_name: user.first_name,
        last_name: user.last_name,
        email: user.email,
        username: user.username,
        role_id: user.role.id,
      });

      // remove password completely in edit mode
      unregister('password');
    } else {
      reset({
        first_name: '',
        last_name: '',
        email: '',
        username: '',
        password: '',
        role_id: undefined,
      });
    }
  }, [isOpen, user, reset, unregister]);

  const onSubmit = (data: UserFormInputs) => {
    const payload =
      user && !data.password
        ? { ...data, password: undefined }
        : data;

    onSave(payload, user ? user.id : null);
  };

  return (
    <Transition appear show={isOpen} as={Fragment}>
      <Dialog as="div" className="relative z-10" onClose={onClose}>
        <Transition.Child
          as={Fragment}
          enter="ease-out duration-300"
          enterFrom="opacity-0"
          enterTo="opacity-100"
          leave="ease-in duration-200"
          leaveFrom="opacity-100"
          leaveTo="opacity-0"
        >
          <div className="fixed inset-0 bg-black/25" />
        </Transition.Child>

        <div className="fixed inset-0 overflow-y-auto">
          <div className="flex min-h-full items-center justify-center p-4">
            <Transition.Child
              as={Fragment}
              enter="ease-out duration-300"
              enterFrom="opacity-0 scale-95"
              enterTo="opacity-100 scale-100"
              leave="ease-in duration-200"
              leaveFrom="opacity-100 scale-100"
              leaveTo="opacity-95"
            >
              <Dialog.Panel className="w-full max-w-md rounded-2xl bg-white p-6 shadow-xl">
                <Dialog.Title className="text-lg font-semibold">
                  {user ? 'Edit User' : 'Create User'}
                </Dialog.Title>

                <form
                  onSubmit={handleSubmit(onSubmit)}
                  className="mt-4 space-y-4"
                >
                  {['first_name', 'last_name', 'email', 'username'].map(
                    field => (
                      <div key={field}>
                        <label className="block text-sm font-medium capitalize">
                          {field.replace('_', ' ')}
                        </label>
                        <input
                          {...register(field as keyof UserFormInputs)}
                          className="mt-1 w-full rounded-md border px-3 py-2"
                        />
                        {errors[field as keyof UserFormInputs] && (
                          <p className="text-sm text-red-500">
                            {
                              errors[field as keyof UserFormInputs]
                                ?.message as string
                            }
                          </p>
                        )}
                      </div>
                    )
                  )}

                  {!user && (
                    <div>
                      <label className="block text-sm font-medium">
                        Password
                      </label>
                      <input
                        type="password"
                        {...register('password')}
                        className="mt-1 w-full rounded-md border px-3 py-2"
                      />
                      {errors.password && (
                        <p className="text-sm text-red-500">
                          {errors.password.message}
                        </p>
                      )}
                    </div>
                  )}

                  <div>
                    <label className="block text-sm font-medium">Role</label>
                    <select
                      {...register('role_id')}
                      className="mt-1 w-full rounded-md border px-3 py-2"
                    >
                      <option value="">Select role</option>
                      {roles?.map(role => (
                        <option key={role.id} value={role.id}>
                          {role.name}
                        </option>
                      ))}
                    </select>
                    {errors.role_id && (
                      <p className="text-sm text-red-500">
                        {errors.role_id.message}
                      </p>
                    )}
                  </div>

                  <div className="flex justify-end gap-2 pt-4">
                    <button
                      type="button"
                      onClick={onClose}
                      className="rounded-md bg-gray-100 px-4 py-2"
                    >
                      Cancel
                    </button>
                    <button
                      type="submit"
                      className="rounded-md bg-indigo-600 px-4 py-2 text-white"
                    >
                      Save
                    </button>
                  </div>
                </form>
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </div>
      </Dialog>
    </Transition>
  );
};

export default UserFormModal;
