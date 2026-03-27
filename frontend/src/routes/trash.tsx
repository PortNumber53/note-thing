import { useEffect } from 'react'
import { useNotesStore } from '@/stores/notes-store'
import { TrashList } from '@/components/trash/trash-list'

export function TrashView() {
  const { fetchTrashed } = useNotesStore()

  useEffect(() => {
    fetchTrashed()
  }, [fetchTrashed])

  return (
    <div className="flex h-full w-80 flex-col border-r">
      <div className="border-b px-3 py-2">
        <h2 className="text-sm font-semibold">Trash</h2>
      </div>
      <TrashList />
    </div>
  )
}
