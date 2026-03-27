import { useState } from 'react'
import { X, Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useTagsStore } from '@/stores/tags-store'
import type { Tag } from '@/types'

type Props = {
  noteId: string
  tags: Tag[]
  onTagsChanged: () => void
}

export function NoteTags({ noteId, tags, onTagsChanged }: Props) {
  const [adding, setAdding] = useState(false)
  const [newTagName, setNewTagName] = useState('')
  const { tags: allTags, createTag, setNoteTags } = useTagsStore()

  const handleAddTag = async () => {
    if (!newTagName.trim()) return

    let tag = allTags.find((t) => t.name.toLowerCase() === newTagName.trim().toLowerCase())
    if (!tag) {
      tag = await createTag(newTagName.trim())
    }

    const currentIds = tags.map((t) => t.id)
    if (!currentIds.includes(tag.id)) {
      await setNoteTags(noteId, [...currentIds, tag.id])
      onTagsChanged()
    }

    setNewTagName('')
    setAdding(false)
  }

  const handleRemoveTag = async (tagId: string) => {
    const newIds = tags.filter((t) => t.id !== tagId).map((t) => t.id)
    await setNoteTags(noteId, newIds)
    onTagsChanged()
  }

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault()
      handleAddTag()
    } else if (e.key === 'Escape') {
      setAdding(false)
      setNewTagName('')
    }
  }

  return (
    <div className="flex items-center gap-1.5 flex-wrap px-4 py-1.5 border-b">
      {tags.map((tag) => (
        <span
          key={tag.id}
          className="inline-flex items-center gap-1 rounded-md bg-secondary px-2 py-0.5 text-xs text-secondary-foreground"
        >
          {tag.name}
          <button
            onClick={() => handleRemoveTag(tag.id)}
            className="hover:text-destructive"
          >
            <X className="h-3 w-3" />
          </button>
        </span>
      ))}
      {adding ? (
        <Input
          value={newTagName}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => setNewTagName(e.target.value)}
          onKeyDown={handleKeyDown}
          onBlur={() => { if (!newTagName.trim()) setAdding(false) }}
          placeholder="Tag name"
          className="h-6 w-24 border-0 px-1 text-xs shadow-none focus-visible:ring-0"
          autoFocus
        />
      ) : (
        <Button
          variant="ghost"
          size="sm"
          className="h-6 px-1.5 text-xs text-muted-foreground"
          onClick={() => setAdding(true)}
        >
          <Plus className="mr-1 h-3 w-3" />
          Add tag
        </Button>
      )}
    </div>
  )
}
