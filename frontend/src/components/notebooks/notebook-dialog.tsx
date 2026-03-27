import { useState } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { useNotebooksStore } from '@/stores/notebooks-store'

type Props = {
  children: React.ReactNode
  notebookId?: string
  initialName?: string
}

export function NotebookDialog({ children, notebookId, initialName = '' }: Props) {
  const [open, setOpen] = useState(false)
  const [name, setName] = useState(initialName)
  const { createNotebook, updateNotebook } = useNotebooksStore()

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (!name.trim()) return
    if (notebookId) {
      await updateNotebook(notebookId, name.trim())
    } else {
      await createNotebook(name.trim())
    }
    setName('')
    setOpen(false)
  }

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{notebookId ? 'Rename Notebook' : 'New Notebook'}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Notebook name"
            autoFocus
          />
          <Button type="submit" className="w-full">
            {notebookId ? 'Rename' : 'Create'}
          </Button>
        </form>
      </DialogContent>
    </Dialog>
  )
}
