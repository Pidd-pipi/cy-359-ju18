import { create } from 'zustand'
import * as authApi from '../api/user'
import { getToken, setAuth, clearAuth, getUser, type AuthUser } from '../utils/auth'

interface AuthState {
  token: string | null
  user: AuthUser | null
  loading: boolean
  init: () => void
  login: (username: string, password: string) => Promise<void>
  register: (data: { username: string; password: string; nickname: string; email?: string; phone?: string }) => Promise<void>
  logout: () => void
  refreshProfile: () => Promise<void>
}

export const useAuthStore = create<AuthState>((set) => ({
  token: getToken(),
  user: getUser(),
  loading: false,

  init: () => {
    set({ token: getToken(), user: getUser() })
  },

  login: async (username, password) => {
    const res = await authApi.login({ username, password })
    setAuth(res.data.token, res.data.user)
    set({ token: res.data.token, user: res.data.user })
  },

  register: async (data) => {
    await authApi.register(data)
  },

  logout: () => {
    clearAuth()
    set({ token: null, user: null })
  },

  refreshProfile: async () => {
    const res = await authApi.getProfile()
    const current = getUser()
    if (current) {
      const merged: AuthUser = {
        id: res.data.id,
        username: res.data.username,
        nickname: res.data.nickname || current.nickname,
        role: res.data.role,
        points: res.data.points,
      }
      setAuth(getToken() || '', merged)
      set({ user: merged })
    }
  },
}))
