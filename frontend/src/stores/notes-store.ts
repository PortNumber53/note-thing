import { create } from 'zustand'
import type { Note } from '@/types'
import { api } from '@/lib/api'
import { useCryptoStore } from './crypto-store'
import { encryptNote, decryptNote } from '@/lib/crypto'
import { rebuildIndex } from '@/lib/search-index'

async function decryptNotes(notes: Note[]): Promise<Note[]> {
  const crypto = useCryptoStore.getState()
  if (!crypto.isUnlocked || !crypto.kek) return notes

  const result: Note[] = []
  for (const note of notes) {
    if (note.isEncrypted && note.encryptedTitle && note.encryptedBody && note.noteKeyWrapped) {
      try {
        const cached = crypto.getCachedNoteKey(note.id)
        const decrypted = await decryptNote(crypto.kek, note.encryptedTitle, note.encryptedBody, note.noteKeyWrapped)
        if (!cached) crypto.cacheNoteKey(note.id, decrypted.dek)
        result.push({ ...note, title: decrypted.title, body: decrypted.body })
      } catch {
        result.push({ ...note, title: '[Decryption failed]', body: '' })
      }
    } else {
      result.push(note)
    }
  }

  // Rebuild client-side search index with decrypted content
  if (crypto.isEncryptionEnabled) {
    rebuildIndex(result.map((n) => ({ id: n.id, title: n.title, body: n.body })))
  }

  return result
}

async function encryptNoteData(
  data: { title: string; body: string; notebookId?: string; tagIds?: string[] },
  noteId?: string,
): Promise<Record<string, unknown>> {
  const crypto = useCryptoStore.getState()
  if (!crypto.isEncryptionEnabled || !crypto.isUnlocked || !crypto.kek) {
    return data
  }

  const cachedDEK = noteId ? crypto.getCachedNoteKey(noteId) : undefined
  const encrypted = await encryptNote(crypto.kek, data.title || '', data.body || '', cachedDEK)

  if (noteId) crypto.cacheNoteKey(noteId, encrypted.dek)

  return {
    title: '',
    body: '',
    encryptedTitle: encrypted.encryptedTitle,
    encryptedBody: encrypted.encryptedBody,
    noteKeyWrapped: encrypted.noteKeyWrapped,
    keyVersion: crypto.keyVersion,
    isEncrypted: true,
    notebookId: data.notebookId,
    tagIds: data.tagIds,
  }
}

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
      const raw = await api.get<Note[]>(`/api/notes/${qs ? `?${qs}` : ''}`)
      const notes = await decryptNotes(raw)
      set({ notes })
    } finally {
      set({ isLoading: false })
    }
  },

  fetchTrashed: async () => {
    set({ isLoading: true })
    try {
      const raw = await api.get<Note[]>('/api/notes/trash')
      const notes = await decryptNotes(raw)
      set({ notes })
    } finally {
      set({ isLoading: false })
    }
  },

  createNote: async (data) => {
    const payload = await encryptNoteData(data)
    const raw = await api.post<Note>('/api/notes', payload)
    const [note] = await decryptNotes([raw])
    set((s) => ({ notes: [note, ...s.notes] }))
    return note
  },

  updateNote: async (id, data) => {
    const payload = await encryptNoteData({ title: data.title || '', body: data.body || '', notebookId: data.notebookId }, id)
    const raw = await api.put<Note>(`/api/notes/${id}`, payload)
    const [note] = await decryptNotes([raw])
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
