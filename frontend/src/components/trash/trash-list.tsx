import { useNotesStore } from '@/stores/notes-store'
import { Button } from '@/components/ui/button'
import { RotateCcw, Trash2 } from 'lucide-react'
import { ScrollArea } from '@/components/ui/scroll-area'
import { cn } from '@/lib/utils'

function stripMarkdown(text: string) {
  return text
    .replace(/!\[([^\]]*)\]\([^)]*\)/g, '$1')
    .replace(/\[([^\]]*)\]\([^)]*\)/g, '$1')
    .replace(/^#{1,6}\s+/gm, '')
    .replace(/^>\s?/gm, '')
    .replace(/^[-*+]\s+/gm, '')
    .replace(/^\d+\\?\.\s+/gm, '')
    .replace(/^---+$/gm, '')
    .replace(/\\([.!#*_`~[\]()>+\-])/g, '$1')
    .replace(/[*_`~[\]]/g, '')
    .trim()
}

export function TrashList({ activeId, onSelect }: { activeId: string | null; onSelect: (id: string) => void }) {
  const { notes, restoreNote, permanentDeleteNote } = useNotesStore()

  return (
    <ScrollArea className="flex-1">
      {notes.length === 0 && (
        <div className="p-4 text-center text-sm text-muted-foreground">Trash is empty</div>
      )}
      {notes.map((note) => (
        <button
          key={note.id}
          onClick={() => onSelect(note.id)}
          className={cn(
            "w-full text-left p-3 border-b transition-colors hover:bg-accent",
            activeId === note.id && "bg-accent"
          )}
        >
          <div className="flex items-baseline justify-between gap-2">
            <h3 className="text-sm font-medium truncate">{note.title || 'Untitled'}</h3>
            <span className="text-xs text-muted-foreground shrink-0">
              {new Date(note.updatedAt).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })}
            </span>
          </div>
          <p className="mt-1 text-xs text-muted-foreground line-clamp-2">
            {stripMarkdown(note.body.slice(0, 160)).slice(0, 80)}
          </p>
        </button>
      ))}
    </ScrollArea>
  )
}
