import { useEffect } from 'react'
import { NoteList } from '@/components/layout/note-list'
import { EditorPane } from '@/components/layout/editor-pane'
import { useNotesStore } from '@/stores/notes-store'

export function AllNotesView() {
  const { fetchNotes, createNote, setActiveNoteId } = useNotesStore()

  useEffect(() => {
    fetchNotes()
  }, [fetchNotes])

  const handleNewNote = async () => {
    const note = await createNote({ title: '', body: '' })
    setActiveNoteId(note.id)
  }

  return (
    <>
      <NoteList title="All Notes" onNewNote={handleNewNote} />
      <EditorPane />
    </>
  )
}
