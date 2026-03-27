import { create } from 'zustand'
import type { Notebook } from '@/types'
import { api } from '@/lib/api'

type NotebooksState = {
  notebooks: Notebook[]
  isLoading: boolean
  fetchNotebooks: () => Promise<void>
  createNotebook: (name: string) => Promise<Notebook>
  updateNotebook: (id: string, name: string) => Promise<Notebook>
  deleteNotebook: (id: string) => Promise<void>
}

export const useNotebooksStore = create<NotebooksState>()((set) => ({
  notebooks: [],
  isLoading: false,

  fetchNotebooks: async () => {
    set({ isLoading: true })
    try {
      const notebooks = await api.get<Notebook[]>('/api/notebooks/')
      set({ notebooks })
    } finally {
      set({ isLoading: false })
    }
  },

  createNotebook: async (name) => {
    const nb = await api.post<Notebook>('/api/notebooks/', { name })
    set((s) => ({ notebooks: [...s.notebooks, nb] }))
    return nb
  },

  updateNotebook: async (id, name) => {
    const nb = await api.put<Notebook>(`/api/notebooks/${id}`, { name })
    set((s) => ({ notebooks: s.notebooks.map((n) => (n.id === id ? nb : n)) }))
    return nb
  },

  deleteNotebook: async (id) => {
    await api.delete(`/api/notebooks/${id}`)
    set((s) => ({ notebooks: s.notebooks.filter((n) => n.id !== id) }))
  },
}))
