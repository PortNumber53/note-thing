import { useState, useEffect } from 'react'
import { useNotesStore } from '@/stores/notes-store'
import { useThemeStore } from '@/stores/theme-store'
import { TrashList } from '@/components/trash/trash-list'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { RotateCcw, Trash2, FileText } from 'lucide-react'
import MDEditor from '@uiw/react-md-editor'

export function TrashView() {
  const { notes, fetchTrashed, restoreNote, permanentDeleteNote } = useNotesStore()
  const [activeId, setActiveId] = useState<string | null>(null)
  const activeNote = notes.find((n) => n.id === activeId)
  const colorMode = useThemeStore((s) => s.resolved)

  useEffect(() => {
    fetchTrashed()
  }, [fetchTrashed])

  // Clear selection if active note is removed (restored/deleted)
  useEffect(() => {
    if (activeId && !notes.find((n) => n.id === activeId)) {
      setActiveId(null)
    }
  }, [notes, activeId])

  return (
    <>
      <div className="flex h-full w-72 flex-col border-r">
        <div className="border-b px-3 py-2">
          <h2 className="text-sm font-semibold">Trash</h2>
        </div>
        <TrashList activeId={activeId} onSelect={setActiveId} />
      </div>

      {activeNote ? (
        <div className="flex flex-1 flex-col">
          <div className="flex items-center gap-2 border-b px-4 py-2">
            <h1 className="flex-1 text-lg font-semibold truncate">
              {activeNote.title || 'Untitled'}
            </h1>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => {
                restoreNote(activeNote.id)
                setActiveId(null)
              }}
            >
              <RotateCcw className="mr-2 h-3.5 w-3.5" />
              Restore
            </Button>
            <Button
              variant="ghost"
              size="sm"
              className="text-destructive"
              onClick={() => {
                permanentDeleteNote(activeNote.id)
                setActiveId(null)
              }}
            >
              <Trash2 className="mr-2 h-3.5 w-3.5" />
              Delete
            </Button>
          </div>
          <div className="flex-1 overflow-auto" data-color-mode={colorMode}>
            <MDEditor
              value={activeNote.body}
              height="100%"
              visibleDragbar={false}
              preview="preview"
              hideToolbar
            />
          </div>
        </div>
      ) : (
        <div className="flex flex-1 items-center justify-center text-muted-foreground">
          <div className="text-center">
            <FileText className="mx-auto h-12 w-12 mb-2 opacity-50" />
            <p>Select a note to view</p>
          </div>
        </div>
      )}
    </>
  )
}
