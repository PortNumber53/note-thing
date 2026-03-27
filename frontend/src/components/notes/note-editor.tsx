import { useState, useEffect, useCallback } from 'react'
import MDEditor from '@uiw/react-md-editor'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Save, Trash2 } from 'lucide-react'
import { useNotesStore } from '@/stores/notes-store'
import { useThemeStore } from '@/stores/theme-store'
import { useDebounce } from '@/hooks/use-debounce'
import { NoteTags } from './note-tags'
import type { Note } from '@/types'

type Props = {
  note: Note
}

export function NoteEditor({ note }: Props) {
  const [title, setTitle] = useState(note.title)
  const [body, setBody] = useState(note.body)
  const [saving, setSaving] = useState(false)
  const { updateNote, deleteNote } = useNotesStore()
  const colorMode = useThemeStore((s) => s.resolved)

  useEffect(() => {
    setTitle(note.title)
    setBody(note.body)
  }, [note.id, note.title, note.body])

  const debouncedBody = useDebounce(body, 1000)
  const debouncedTitle = useDebounce(title, 1000)

  useEffect(() => {
    if (debouncedTitle !== note.title || debouncedBody !== note.body) {
      updateNote(note.id, { title: debouncedTitle, body: debouncedBody })
    }
  }, [debouncedTitle, debouncedBody, note.id, note.title, note.body, updateNote])

  const handleSave = useCallback(async () => {
    setSaving(true)
    try {
      await updateNote(note.id, { title, body })
    } finally {
      setSaving(false)
    }
  }, [note.id, title, body, updateNote])

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center gap-2 border-b px-4 py-2">
        <Input
          value={title}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => setTitle(e.target.value)}
          placeholder="Note title"
          className="border-0 text-lg font-semibold shadow-none focus-visible:ring-0"
        />
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8"
          onClick={handleSave}
          disabled={saving}
        >
          <Save className="h-4 w-4" />
        </Button>
        <Button
          variant="ghost"
          size="icon"
          className="h-8 w-8 text-destructive"
          onClick={() => deleteNote(note.id)}
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>
      <NoteTags
        noteId={note.id}
        tags={note.tags}
        onTagsChanged={() => useNotesStore.getState().fetchNotes()}
      />
      <div className="flex-1 overflow-auto" data-color-mode={colorMode}>
        <MDEditor
          value={body}
          onChange={(val) => setBody(val || '')}
          height="100%"
          visibleDragbar={false}
          preview="live"
        />
      </div>
    </div>
  )
}
