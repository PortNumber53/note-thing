import { create } from 'zustand'
import type { UserSettings } from '@/types'
import { api } from '@/lib/api'

type SettingsState = {
  settings: UserSettings | null
  isLoading: boolean
  fetchSettings: () => Promise<void>
  updateSettings: (data: Partial<UserSettings>) => Promise<void>
}

export const useSettingsStore = create<SettingsState>()((set, get) => ({
  settings: null,
  isLoading: false,

  fetchSettings: async () => {
    set({ isLoading: true })
    try {
      const settings = await api.get<UserSettings>('/api/settings')
      set({ settings })
    } finally {
      set({ isLoading: false })
    }
  },

  updateSettings: async (data: Partial<UserSettings>) => {
    const current = get().settings
    const merged = { ...current, ...data }
    const settings = await api.put<UserSettings>('/api/settings', merged)
    set({ settings })
  },
}))
