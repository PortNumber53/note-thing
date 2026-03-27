import { useEffect } from 'react'
import { useParams } from 'react-router'
import { NoteList } from '@/components/layout/note-list'
import { EditorPane } from '@/components/layout/editor-pane'
import { useNotesStore } from '@/stores/notes-store'
import { useTagsStore } from '@/stores/tags-store'

export function TagNotesView() {
  const { tagId } = useParams()
  const fetchNotes = useNotesStore.getState().fetchNotes
  const { tags } = useTagsStore()
  const tag = tags.find((t) => t.id === tagId)

  useEffect(() => {
    if (tagId) fetchNotes({ tag_id: tagId })
  }, [tagId, fetchNotes])

  return (
    <>
      <NoteList title={tag ? `#${tag.name}` : 'Tag'} />
      <EditorPane />
    </>
  )
}
