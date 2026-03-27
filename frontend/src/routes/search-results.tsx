import { ScrollArea } from '@/components/ui/scroll-area'
import { NoteListItem } from '@/components/layout/note-list-item'
import { EditorPane } from '@/components/layout/editor-pane'
import { SearchBar } from '@/components/search/search-bar'
import { useNotesStore } from '@/stores/notes-store'

export function SearchView() {
  const { notes, isLoading, activeNoteId, setActiveNoteId } = useNotesStore()

  return (
    <>
      <div className="flex h-full w-72 flex-col border-r">
        <div className="border-b p-2">
          <SearchBar />
        </div>
        <ScrollArea className="flex-1">
          {isLoading && (
            <div className="p-4 text-center text-sm text-muted-foreground">Searching...</div>
          )}
          {!isLoading && notes.length === 0 && (
            <div className="p-4 text-center text-sm text-muted-foreground">No results</div>
          )}
          {notes.map((note) => (
            <NoteListItem
              key={note.id}
              note={note}
              isActive={note.id === activeNoteId}
              onClick={() => setActiveNoteId(note.id)}
            />
          ))}
        </ScrollArea>
      </div>
      <EditorPane />
    </>
  )
}
