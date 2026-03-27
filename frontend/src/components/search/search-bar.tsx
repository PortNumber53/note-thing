import { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router'
import { Search } from 'lucide-react'
import { Input } from '@/components/ui/input'
import { useDebounce } from '@/hooks/use-debounce'
import { api } from '@/lib/api'
import { useNotesStore } from '@/stores/notes-store'
import type { Note } from '@/types'

export function SearchBar() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [query, setQuery] = useState(searchParams.get('q') || '')
  const debouncedQuery = useDebounce(query, 300)

  useEffect(() => {
    if (debouncedQuery) {
      setSearchParams({ q: debouncedQuery })
      api.get<Note[]>(`/api/search?q=${encodeURIComponent(debouncedQuery)}`).then((notes) => {
        useNotesStore.setState({ notes, isLoading: false })
      })
    } else {
      useNotesStore.setState({ notes: [], isLoading: false })
    }
  }, [debouncedQuery, setSearchParams])

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
