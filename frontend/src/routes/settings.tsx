import { useEffect, useState, useRef } from 'react'
import { Sun, Moon, Monitor, Check } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { useThemeStore } from '@/stores/theme-store'
import { useEditorPrefsStore } from '@/stores/editor-prefs-store'
import { useSettingsStore } from '@/stores/settings-store'
import { useNotebooksStore } from '@/stores/notebooks-store'

export function SettingsView() {
  const { setting: themeSetting, setTheme } = useThemeStore()
  const { previewMode, setPreviewMode } = useEditorPrefsStore()
  const { settings, fetchSettings, updateSettings } = useSettingsStore()
  const { notebooks, fetchNotebooks } = useNotebooksStore()

  useEffect(() => {
    fetchSettings()
    fetchNotebooks()
  }, [fetchSettings, fetchNotebooks])

  return (
    <div className="flex-1 overflow-auto p-8">
      <h1 className="text-2xl font-bold mb-6">Settings</h1>
      <div className="max-w-md space-y-8">

        <div className="space-y-3">
          <h2 className="text-sm font-medium uppercase text-muted-foreground">Appearance</h2>
          <div className="flex items-center justify-between">
            <Label>Theme</Label>
            <div className="flex gap-1">
              <Button
                variant={themeSetting === 'system' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setTheme('system')}
              >
                <Monitor className="mr-1.5 h-4 w-4" />
                System
              </Button>
              <Button
                variant={themeSetting === 'light' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setTheme('light')}
              >
                <Sun className="mr-1.5 h-4 w-4" />
                Light
              </Button>
              <Button
                variant={themeSetting === 'dark' ? 'default' : 'outline'}
                size="sm"
                onClick={() => setTheme('dark')}
              >
                <Moon className="mr-1.5 h-4 w-4" />
                Dark
              </Button>
            </div>
          </div>
        </div>

        <Separator />

        <div className="space-y-3">
          <h2 className="text-sm font-medium uppercase text-muted-foreground">Editor</h2>
          <div className="flex items-center justify-between">
            <Label>Preview mode</Label>
            <Select value={previewMode} onValueChange={(v) => setPreviewMode(v as 'split' | 'live' | 'source')}>
              <SelectTrigger className="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="live">Live</SelectItem>
                <SelectItem value="split">Split</SelectItem>
                <SelectItem value="source">Source</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <Separator />

        <div className="space-y-3">
          <h2 className="text-sm font-medium uppercase text-muted-foreground">Defaults</h2>
          <div className="space-y-1.5">
            <Label>Default notebook</Label>
            <NotebookAutocomplete
              notebooks={notebooks}
              selectedId={settings?.defaultNotebookId || null}
              onSelect={(id) => updateSettings({ defaultNotebookId: id })}
            />
          </div>
        </div>
      </div>
    </div>
  )
}

function NotebookAutocomplete({
  notebooks,
  selectedId,
  onSelect,
}: {
  notebooks: { id: string; name: string }[]
  selectedId: string | null
  onSelect: (id: string) => void
}) {
  const { createNotebook, fetchNotebooks } = useNotebooksStore()
  const selected = notebooks.find((nb) => nb.id === selectedId)
  const [query, setQuery] = useState(selected?.name || '')
  const [open, setOpen] = useState(false)
  const [creating, setCreating] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    setQuery(selected?.name || '')
  }, [selected])

  const filtered = notebooks.filter((nb) =>
    nb.name.toLowerCase().includes(query.toLowerCase())
  )
  const exactMatch = notebooks.find((nb) => nb.name.toLowerCase() === query.trim().toLowerCase())

  const handleSelect = (id: string, name: string) => {
    onSelect(id)
    setQuery(name)
    setOpen(false)
  }

  const handleCreate = async () => {
    if (!query.trim() || exactMatch) return
    setCreating(true)
    try {
      const nb = await createNotebook(query.trim())
      await fetchNotebooks()
      onSelect(nb.id)
      setQuery(nb.name)
      setOpen(false)
    } finally {
      setCreating(false)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      if (filtered.length === 1) {
        handleSelect(filtered[0].id, filtered[0].name)
      } else if (!exactMatch && query.trim()) {
        handleCreate()
      }
    } else if (e.key === 'Escape') {
      setOpen(false)
      setQuery(selected?.name || '')
    }
  }

  return (
    <div className="relative">
      <Input
        ref={inputRef}
        value={query}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
          setQuery(e.target.value)
          setOpen(true)
        }}
        onFocus={() => setOpen(true)}
        onBlur={() => setTimeout(() => setOpen(false), 150)}
        onKeyDown={handleKeyDown}
        placeholder="Type to search or create..."
        className="w-full"
      />
      {open && (query || filtered.length > 0) && (
        <div className="absolute z-50 mt-1 w-full rounded-md border bg-popover shadow-md">
          <div className="max-h-48 overflow-auto p-1">
            {filtered.map((nb) => (
              <button
                key={nb.id}
                className="flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm hover:bg-accent text-left"
                onMouseDown={(e) => e.preventDefault()}
                onClick={() => handleSelect(nb.id, nb.name)}
              >
                {nb.id === selectedId && <Check className="h-3.5 w-3.5 text-primary" />}
                {nb.id !== selectedId && <span className="w-3.5" />}
                {nb.name}
              </button>
            ))}
            {query.trim() && !exactMatch && (
              <button
                className="flex w-full items-center gap-2 rounded-sm px-2 py-1.5 text-sm hover:bg-accent text-left text-primary"
                onMouseDown={(e) => e.preventDefault()}
                onClick={handleCreate}
                disabled={creating}
              >
                <span className="w-3.5">+</span>
                Create &ldquo;{query.trim()}&rdquo;
              </button>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
