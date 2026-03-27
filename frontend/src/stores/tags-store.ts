import { create } from 'zustand'
import type { Tag } from '@/types'
import { api } from '@/lib/api'

type TagsState = {
  tags: Tag[]
  isLoading: boolean
  fetchTags: () => Promise<void>
  createTag: (name: string) => Promise<Tag>
  updateTag: (id: string, name: string) => Promise<Tag>
  deleteTag: (id: string) => Promise<void>
  setNoteTags: (noteId: string, tagIds: string[]) => Promise<void>
}

export const useTagsStore = create<TagsState>()((set) => ({
  tags: [],
  isLoading: false,

  fetchTags: async () => {
    set({ isLoading: true })
    try {
      const tags = await api.get<Tag[]>('/api/tags/')
      set({ tags })
    } finally {
      set({ isLoading: false })
    }
  },

  createTag: async (name) => {
    const tag = await api.post<Tag>('/api/tags/', { name })
    set((s) => ({ tags: [...s.tags, tag] }))
    return tag
  },

  updateTag: async (id, name) => {
    const tag = await api.put<Tag>(`/api/tags/${id}`, { name })
    set((s) => ({ tags: s.tags.map((t) => (t.id === id ? tag : t)) }))
    return tag
  },

  deleteTag: async (id) => {
    await api.delete(`/api/tags/${id}`)
    set((s) => ({ tags: s.tags.filter((t) => t.id !== id) }))
  },

  setNoteTags: async (noteId, tagIds) => {
    await api.put(`/api/notes/${noteId}/tags`, { tagIds })
  },
}))
