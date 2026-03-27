import type { Note } from '@/types'
import { cn } from '@/lib/utils'

type Props = {
  note: Note
  isActive: boolean
  onClick: () => void
}

export function NoteListItem({ note, isActive, onClick }: Props) {
  const preview = note.body.slice(0, 120).replace(/[#*_`~\[\]]/g, '')
  const date = new Date(note.updatedAt).toLocaleDateString(undefined, {
    month: 'short',
    day: 'numeric',
  })

  return (
    <button
      onClick={onClick}
      className={cn(
        "w-full text-left p-3 border-b transition-colors hover:bg-accent",
        isActive && "bg-accent"
      )}
    >
      <div className="flex items-baseline justify-between gap-2">
        <h3 className="font-medium truncate text-sm">
          {note.title || 'Untitled'}
        </h3>
        <span className="text-xs text-muted-foreground shrink-0">{date}</span>
      </div>
      <p className="mt-1 text-xs text-muted-foreground line-clamp-2">{preview}</p>
      {note.tags.length > 0 && (
        <div className="mt-1.5 flex gap-1 flex-wrap">
          {note.tags.map((tag) => (
            <span key={tag.id} className="rounded bg-secondary px-1.5 py-0.5 text-[10px] text-secondary-foreground">
              {tag.name}
            </span>
          ))}
        </div>
      )}
    </button>
  )
}
