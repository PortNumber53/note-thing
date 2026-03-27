import { useState } from 'react'
import { useAuthStore } from '@/stores/auth-store'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'

export function ProfileView() {
  const { user, updateUser, deleteAccount } = useAuthStore()
  const [name, setName] = useState(user?.name || '')
  const [saving, setSaving] = useState(false)
  const [saved, setSaved] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const [deleteConfirm, setDeleteConfirm] = useState('')

  if (!user) return null

  const memberSince = new Date(user.createdAt).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })

  const handleSave = async () => {
    if (!name.trim()) return
    setSaving(true)
    try {
      await updateUser({ name: name.trim() })
      setSaved(true)
      setTimeout(() => setSaved(false), 2000)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async () => {
    await deleteAccount()
  }

  return (
    <div className="flex-1 overflow-auto p-8">
      <h1 className="text-2xl font-bold mb-6">Profile</h1>
      <div className="max-w-md space-y-8">
        <div className="flex items-center gap-4">
          <img
            src={user.avatarUrl}
            alt={user.name}
            className="h-16 w-16 rounded-full"
          />
          <div>
            <p className="text-sm text-muted-foreground">{user.email}</p>
            <p className="text-xs text-muted-foreground">Member since {memberSince}</p>
          </div>
        </div>

        <div className="space-y-2">
          <Label htmlFor="name">Display name</Label>
          <div className="flex gap-2">
            <Input
              id="name"
              value={name}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setName(e.target.value)}
              placeholder="Your name"
            />
            <Button onClick={handleSave} disabled={saving || name.trim() === user.name}>
              {saved ? 'Saved' : 'Save'}
            </Button>
          </div>
        </div>

        <Separator />

        <div className="space-y-2">
          <h2 className="text-lg font-semibold text-destructive">Danger Zone</h2>
          <p className="text-sm text-muted-foreground">
            Permanently delete your account and all associated data. This action cannot be undone.
          </p>
          <Dialog open={deleteOpen} onOpenChange={setDeleteOpen}>
            <DialogTrigger asChild>
              <Button variant="destructive">Delete Account</Button>
            </DialogTrigger>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Delete Account</DialogTitle>
              </DialogHeader>
              <p className="text-sm text-muted-foreground">
                This will permanently delete your account, all notebooks, notes, and tags. Type <strong>delete</strong> to confirm.
              </p>
              <Input
                value={deleteConfirm}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setDeleteConfirm(e.target.value)}
                placeholder='Type "delete" to confirm'
              />
              <Button
                variant="destructive"
                className="w-full"
                disabled={deleteConfirm !== 'delete'}
                onClick={handleDelete}
              >
                Permanently Delete My Account
              </Button>
            </DialogContent>
          </Dialog>
        </div>
      </div>
    </div>
  )
}
