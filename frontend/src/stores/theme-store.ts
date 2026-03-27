import { create } from 'zustand'

type ThemeSetting = 'system' | 'light' | 'dark'

type ThemeState = {
  setting: ThemeSetting
  resolved: 'light' | 'dark'
  setTheme: (setting: ThemeSetting) => void
}

const STORAGE_KEY = 'note-thing-theme'

function resolveTheme(setting: ThemeSetting): 'light' | 'dark' {
  if (setting === 'system') {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
  }
  return setting
}

function applyTheme(resolved: 'light' | 'dark') {
  document.documentElement.classList.toggle('dark', resolved === 'dark')
}

function getInitialSetting(): ThemeSetting {
  const stored = localStorage.getItem(STORAGE_KEY)
  if (stored === 'light' || stored === 'dark' || stored === 'system') return stored
  return 'system'
}

const initialSetting = getInitialSetting()
const initialResolved = resolveTheme(initialSetting)
applyTheme(initialResolved)

// Listen for system theme changes when set to 'system'
window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
  const state = useThemeStore.getState()
  if (state.setting === 'system') {
    const resolved = resolveTheme('system')
    applyTheme(resolved)
    useThemeStore.setState({ resolved })
  }
})

export const useThemeStore = create<ThemeState>()((set) => ({
  setting: initialSetting,
  resolved: initialResolved,
  setTheme: (setting: ThemeSetting) => {
    localStorage.setItem(STORAGE_KEY, setting)
    const resolved = resolveTheme(setting)
    applyTheme(resolved)
    set({ setting, resolved })
  },
}))
