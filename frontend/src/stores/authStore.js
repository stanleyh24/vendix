import { create } from 'zustand'
import { persist } from 'zustand/middleware'

export const useAuthStore = create(
  persist(
    (set) => ({
      user: null,
      accessToken: null,
      refreshToken: null,
      tenantId: 'demo',
      isAuthenticated: false,

      setAuth: (user, accessToken, refreshToken) => 
        set({ user, accessToken, refreshToken, isAuthenticated: true }),

      setTenant: (tenantId) => 
        set({ tenantId }),

      logout: () => 
        set({ user: null, accessToken: null, refreshToken: null, isAuthenticated: false }),
    }),
    {
      name: 'auth-storage',
    }
  )
)

