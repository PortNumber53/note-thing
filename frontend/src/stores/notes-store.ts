import { create } from 'zustand'
import type { Note } from '@/types'
import { api } from '@/lib/api'

type NotesState = {
  notes: Note[]
  activeNoteId: string | null
  isLoading: boolean
  fetchNotes: (params?: { notebook_id?: string; tag_id?: string }) => Promise<void>
  fetchTrashed: () => Promise<void>
  createNote: (data: { title: string; body: string; notebookId?: string; tagIds?: string[] }) => Promise<Note>
  updateNote: (id: string, data: { title?: string; body?: string; notebookId?: string }) => Promise<Note>
  deleteNote: (id: string) => Promise<void>
  restoreNote: (id: string) => Promise<void>
  permanentDeleteNote: (id: string) => Promise<void>
  setActiveNoteId: (id: string | null) => void
}

export const useNotesStore = create<NotesState>()((set) => ({
  notes: [],
  activeNoteId: null,
  isLoading: false,

  fetchNotes: async (params) => {
    set({ isLoading: true })
    try {
      const query = new URLSearchParams()
      if (params?.notebook_id) query.set('notebook_id', params.notebook_id)
      if (params?.tag_id) query.set('tag_id', params.tag_id)
      const qs = query.toString()
      const notes = await api.get<Note[]>(`/api/notes/${qs ? `?${qs}` : ''}`)
      set({ notes })
    } finally {
      set({ isLoading: false })
    }
  },

  fetchTrashed: async () => {
    set({ isLoading: true })
    try {
      const notes = await api.get<Note[]>('/api/notes/trash')
      set({ notes })
    } finally {
      set({ isLoading: false })
    }
  },

  createNote: async (data) => {
    const note = await api.post<Note>('/api/notes', data)
    set((s) => ({ notes: [note, ...s.notes] }))
    return note
  },

  updateNote: async (id, data) => {
    const note = await api.put<Note>(`/api/notes/${id}`, data)
    set((s) => ({ notes: s.notes.map((n) => (n.id === id ? note : n)) }))
    return note
  },

  deleteNote: async (id) => {
    await api.delete(`/api/notes/${id}`)
    set((s) => ({ notes: s.notes.filter((n) => n.id !== id), activeNoteId: s.activeNoteId === id ? null : s.activeNoteId }))
  },

  restoreNote: async (id) => {
    await api.post(`/api/notes/${id}/restore`)
    set((s) => ({ notes: s.notes.filter((n) => n.id !== id) }))
  },

  permanentDeleteNote: async (id) => {
    await api.delete(`/api/notes/${id}/permanent`)
    set((s) => ({ notes: s.notes.filter((n) => n.id !== id), activeNoteId: s.activeNoteId === id ? null : s.activeNoteId }))
  },

  setActiveNoteId: (id) => set({ activeNoteId: id }),
}))
