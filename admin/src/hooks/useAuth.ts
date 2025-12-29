// import { create } from 'zustand';
// import { persist } from 'zustand/middleware';
// import type { User } from '../types/user';

// interface AuthState {
//   accessToken: string | null;
//   refreshToken: string | null;
//   user: User | null;
//   permissions: string[];
//   isAuthenticated: boolean;

//   setTokens: (accessToken: string, refreshToken: string) => void;
//   login: (accessToken: string, refreshToken: string, user: User) => void;
//   logout: () => void;
// }

// export const useAuthStore = create<AuthState>()(
//   persist(
//     (set) => ({
//       accessToken: null,
//       refreshToken: null,
//       user: null,
//       permissions: [],
//       isAuthenticated: false,

//       setTokens: (accessToken, refreshToken) =>
//         set((state) => ({
//           accessToken,
//           refreshToken,
//           isAuthenticated: Boolean(accessToken),
//           user: state.user,
//           permissions: state.permissions,
//         })),

//       login: (accessToken, refreshToken, user) =>
//         set({
//           accessToken,
//           refreshToken,
//           user,
//           permissions: user.permissions ?? [],
//           isAuthenticated: true,
//         }),

//       logout: () =>
//         set({
//           accessToken: null,
//           refreshToken: null,
//           user: null,
//           permissions: [],
//           isAuthenticated: false,
//         }),
//     }),
//     {
//       name: 'auth-storage',
//       partialize: (state) => ({
//         accessToken: state.accessToken,
//         refreshToken: state.refreshToken,
//         user: state.user,
//         permissions: state.permissions,
//         isAuthenticated: state.isAuthenticated,
//       }),
//     }
//   )
// );
import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import type { User } from '../types/user';

interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  user: User | null;
  isAuthenticated: boolean;
  permissions: string[];
  setTokens: (accessToken: string, refreshToken: string) => void;
  setUser: (user: User) => void;
  login: (accessToken: string, refreshToken: string, user: User) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      accessToken: null,
      refreshToken: null,
      user: null,
      isAuthenticated: false,
      permissions: [],
      setTokens: (accessToken, refreshToken) => set({ accessToken, refreshToken }),
      setUser: (user) => set({ user, isAuthenticated: true }),
      login: (accessToken, refreshToken, user) =>
        set({
          accessToken,
          refreshToken,
          user,
          isAuthenticated: true,
          permissions: user.permissions || [],
        }),
      logout: () =>
        set({
          accessToken: null,
          refreshToken: null,
          user: null,
          isAuthenticated: false,
          permissions: [],
        }),
    }),
    {
      name: 'auth-storage', // unique name for localStorage key
    }
  )
);
