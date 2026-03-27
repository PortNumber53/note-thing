import { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router'
import { Search } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { useDebounce } from '@/hooks/use-debounce'
import { api } from '@/lib/api'
import { useNotesStore } from '@/stores/notes-store'
import { useCryptoStore } from '@/stores/crypto-store'
import { searchNotes as searchLocal } from '@/lib/search-index'
import type { Note } from '@/types'

export function SearchBar() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [query, setQuery] = useState(searchParams.get('q') || '')
  const debouncedQuery = useDebounce(query, 300)
  const isEncryptionEnabled = useCryptoStore((s) => s.isEncryptionEnabled)

  useEffect(() => {
    if (debouncedQuery) {
      setSearchParams({ q: debouncedQuery })

      if (isEncryptionEnabled) {
        // Client-side search on decrypted notes
        const matchingIds = searchLocal(debouncedQuery)
        const allNotes = useNotesStore.getState().notes
        // If notes aren't loaded yet, load them first
        if (allNotes.length === 0) {
          useNotesStore.getState().fetchNotes().then(() => {
            const loaded = useNotesStore.getState().notes
            const filtered = matchingIds.map((id) => loaded.find((n) => n.id === id)).filter(Boolean) as Note[]
            useNotesStore.setState({ notes: filtered, isLoading: false })
          })
        } else {
          const filtered = matchingIds.map((id) => allNotes.find((n) => n.id === id)).filter(Boolean) as Note[]
          useNotesStore.setState({ notes: filtered, isLoading: false })
        }
      } else {
        // Server-side search
        api.get<Note[]>(`/api/search?q=${encodeURIComponent(debouncedQuery)}`).then((notes) => {
          useNotesStore.setState({ notes, isLoading: false })
        })
      }
    } else {
      useNotesStore.setState({ notes: [], isLoading: false })
    }
  }, [debouncedQuery, setSearchParams, isEncryptionEnabled])

  return (
    <div className="relative">
      <Search className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
      <Input
        value={query}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setQuery(e.target.value)}
        placeholder="Search notes..."
        className="pl-8"
        autoFocus
      />
    </div>
  )
}
