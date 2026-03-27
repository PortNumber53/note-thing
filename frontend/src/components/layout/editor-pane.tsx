import { useNotesStore } from '@/stores/notes-store'
import { NoteEditor } from '@/components/notes/note-editor'
import { FileText } from 'lucide-react'

export function EditorPane() {
  const { notes, activeNoteId } = useNotesStore()
  const activeNote = notes.find((n) => n.id === activeNoteId)

  if (!activeNote) {
    return (
      <div className="flex flex-1 items-center justify-center text-muted-foreground">
        <div className="text-center">
          <FileText className="mx-auto h-12 w-12 mb-2 opacity-50" />
          <p>Select a note to start editing</p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex-1">
      <NoteEditor key={activeNote.id} note={activeNote} />
    </div>
  )
}
