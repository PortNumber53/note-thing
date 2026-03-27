import { useNotesStore } from '@/stores/notes-store'
import { Button } from '@/components/ui/button'
import { RotateCcw, Trash2 } from 'lucide-react'
import { ScrollArea } from '@/components/ui/scroll-area'

export function TrashList() {
  const { notes, restoreNote, permanentDeleteNote } = useNotesStore()

  return (
    <ScrollArea className="flex-1">
      {notes.length === 0 && (
        <div className="p-4 text-center text-sm text-muted-foreground">Trash is empty</div>
      )}
      {notes.map((note) => (
        <div key={note.id} className="flex items-center gap-2 border-b p-3">
          <div className="flex-1 min-w-0">
            <h3 className="text-sm font-medium truncate">{note.title || 'Untitled'}</h3>
            <p className="text-xs text-muted-foreground truncate">{note.body.slice(0, 80)}</p>
          </div>
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7 shrink-0"
            onClick={() => restoreNote(note.id)}
            title="Restore"
          >
            <RotateCcw className="h-3.5 w-3.5" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            className="h-7 w-7 shrink-0 text-destructive"
            onClick={() => permanentDeleteNote(note.id)}
            title="Delete permanently"
          >
            <Trash2 className="h-3.5 w-3.5" />
          </Button>
        </div>
      ))}
    </ScrollArea>
  )
}
