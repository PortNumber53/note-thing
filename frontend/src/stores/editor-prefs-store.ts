import { create } from 'zustand'
import type { EditorPreviewMode } from '@/types'

type EditorPrefsState = {
  previewMode: EditorPreviewMode
  setPreviewMode: (mode: EditorPreviewMode) => void
}

const STORAGE_KEY = 'note-thing-editor-preview'

export const useEditorPrefsStore = create<EditorPrefsState>()((set) => ({
  previewMode: (localStorage.getItem(STORAGE_KEY) as EditorPreviewMode) || 'live',
  setPreviewMode: (mode: EditorPreviewMode) => {
    localStorage.setItem(STORAGE_KEY, mode)
    set({ previewMode: mode })
  },
}))
