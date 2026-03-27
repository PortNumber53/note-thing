import { useEffect } from 'react'
import { useParams } from 'react-router'
import { NoteList } from '@/components/layout/note-list'
import { EditorPane } from '@/components/layout/editor-pane'
import { useNotesStore } from '@/stores/notes-store'
import { useNotebooksStore } from '@/stores/notebooks-store'

export function NotebookNotesView() {
  const { notebookId } = useParams()
  const { fetchNotes, createNote, setActiveNoteId } = useNotesStore()
  const notebooks = useNotebooksStore((s) => s.notebooks)
  const notebook = notebooks.find((nb) => nb.id === notebookId)

  useEffect(() => {
    if (notebookId) fetchNotes({ notebook_id: notebookId })
  }, [notebookId, fetchNotes])

  const handleNewNote = async () => {
    const note = await createNote({ title: '', body: '', notebookId })
    setActiveNoteId(note.id)
  }

  return (
    <>
      <NoteList title={notebook?.name || 'Notebook'} onNewNote={handleNewNote} />
      <EditorPane />
    </>
  )
}
