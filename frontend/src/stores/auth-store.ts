import { create } from 'zustand'
import type { User } from '@/types'
import { api } from '@/lib/api'

const STORAGE_KEY_TOKEN = 'note-thing-token'
const STORAGE_KEY_USER = 'note-thing-user'

type AuthState = {
  token: string | null
  user: User | null
  isAuthenticated: boolean
  isLoading: boolean
  init: () => void
  setToken: (token: string) => void
  fetchUser: () => Promise<void>
  logout: () => void
}

export const useAuthStore = create<AuthState>()((set, get) => ({
  token: null,
  user: null,
  isAuthenticated: false,
  isLoading: true,

  init: () => {
    const token = localStorage.getItem(STORAGE_KEY_TOKEN)
    if (token) {
      const storedUser = localStorage.getItem(STORAGE_KEY_USER)
      let user: User | null = null
      if (storedUser) {
        try {
          user = JSON.parse(storedUser)
        } catch {
          localStorage.removeItem(STORAGE_KEY_USER)
        }
      }
      set({ token, user, isAuthenticated: true, isLoading: false })
      if (!user) {
        get().fetchUser()
      }
    } else {
      set({ isLoading: false })
    }
  },

  setToken: (token: string) => {
    localStorage.setItem(STORAGE_KEY_TOKEN, token)
    set({ token, isAuthenticated: true })
  },

  fetchUser: async () => {
    try {
      const user = await api.get<User>('/api/me')
      localStorage.setItem(STORAGE_KEY_USER, JSON.stringify(user))
      set({ user })
    } catch {
      get().logout()
    }
  },

  logout: () => {
    localStorage.removeItem(STORAGE_KEY_TOKEN)
    localStorage.removeItem(STORAGE_KEY_USER)
    set({ token: null, user: null, isAuthenticated: false })
  },
}))
