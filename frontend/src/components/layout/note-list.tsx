import { ScrollArea } from '@/components/ui/scroll-area'
import { NoteListItem } from './note-list-item'
import { useNotesStore } from '@/stores/notes-store'
import { Button } from '@/components/ui/button'
import { Plus } from 'lucide-react'

type Props = {
  title: string
  onNewNote?: () => void
}

export function NoteList({ title, onNewNote }: Props) {
  const { notes, isLoading, activeNoteId, setActiveNoteId } = useNotesStore()

  return (
    <div className="flex h-full w-72 flex-col border-r">
      <div className="flex items-center justify-between border-b px-3 py-2">
        <h2 className="text-sm font-semibold">{title}</h2>
        {onNewNote && (
          <Button variant="ghost" size="icon" className="h-7 w-7" onClick={onNewNote}>
            <Plus className="h-4 w-4" />
          </Button>
        )}
      </div>

      <ScrollArea className="flex-1">
        {isLoading && (
          <div className="p-4 text-center text-sm text-muted-foreground">Loading...</div>
        )}
        {!isLoading && notes.length === 0 && (
          <div className="p-4 text-center text-sm text-muted-foreground">No notes</div>
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
  )
}
