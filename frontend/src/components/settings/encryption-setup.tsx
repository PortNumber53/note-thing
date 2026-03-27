import { useState } from 'react'
import { Shield, Lock, AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'
import { useCryptoStore } from '@/stores/crypto-store'
import { useNotesStore } from '@/stores/notes-store'
import { api } from '@/lib/api'
import { encryptNote } from '@/lib/crypto'
import type { Note } from '@/types'

export function EncryptionSetup() {
  const { isEncryptionEnabled, isUnlocked, kek, keyVersion, lock, setupEncryption } = useCryptoStore()

  if (!isEncryptionEnabled) {
    return <SetupForm onSetup={setupEncryption} />
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-3">
        <Shield className="h-5 w-5 text-green-600" />
        <div>
          <p className="font-medium">End-to-end encryption is enabled</p>
          <p className="text-sm text-muted-foreground">Your notes are encrypted before leaving this device.</p>
        </div>
      </div>

      {isUnlocked && (
        <div className="flex gap-2">
          <Button variant="outline" size="sm" onClick={lock}>
            <Lock className="mr-1.5 h-3.5 w-3.5" />
            Lock now
          </Button>
          <MigrateButton kek={kek} keyVersion={keyVersion} />
        </div>
      )}
    </div>
  )
}

function SetupForm({ onSetup }: { onSetup: (password: string) => Promise<void> }) {
  const [open, setOpen] = useState(false)
  const [password, setPassword] = useState('')
  const [confirm, setConfirm] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  const handleSetup = async () => {
    if (password.length < 8) {
      setError('Password must be at least 8 characters')
      return
    }
    if (password !== confirm) {
      setError('Passwords do not match')
      return
    }
    setLoading(true)
    setError('')
    try {
      await onSetup(password)
      setOpen(false)
    } catch {
      setError('Setup failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-3">
      <div className="flex items-center gap-3">
        <Shield className="h-5 w-5 text-muted-foreground" />
        <div>
          <p className="font-medium">End-to-end encryption</p>
          <p className="text-sm text-muted-foreground">Encrypt your notes so only you can read them.</p>
        </div>
      </div>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          <Button>Enable Encryption</Button>
        </DialogTrigger>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Set Up Encryption</DialogTitle>
          </DialogHeader>
          <div className="rounded-md bg-yellow-50 dark:bg-yellow-950 p-3 flex gap-2">
            <AlertTriangle className="h-4 w-4 text-yellow-600 shrink-0 mt-0.5" />
            <p className="text-sm text-yellow-800 dark:text-yellow-200">
              Choose a strong master password. If you forget it, your encrypted notes cannot be recovered.
            </p>
          </div>
          <div className="space-y-3">
            <div className="space-y-1">
              <Label>Master password</Label>
              <Input
                type="password"
                value={password}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
                placeholder="At least 8 characters"
              />
            </div>
            <div className="space-y-1">
              <Label>Confirm password</Label>
              <Input
                type="password"
                value={confirm}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setConfirm(e.target.value)}
                placeholder="Repeat password"
              />
            </div>
            {error && <p className="text-sm text-destructive">{error}</p>}
            <Button className="w-full" onClick={handleSetup} disabled={loading}>
              {loading ? 'Setting up...' : 'Enable Encryption'}
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}

function MigrateButton({ kek, keyVersion }: { kek: CryptoKey | null; keyVersion: number | null }) {
  const [migrating, setMigrating] = useState(false)
  const [progress, setProgress] = useState({ done: 0, total: 0 })

  const handleMigrate = async () => {
    if (!kek || !keyVersion) return
    setMigrating(true)

    try {
      // Fetch all unencrypted notes
      const allNotes = await api.get<Note[]>('/api/notes/')
      const unencrypted = allNotes.filter((n) => !n.isEncrypted)

      if (unencrypted.length === 0) {
        setMigrating(false)
        return
      }

      setProgress({ done: 0, total: unencrypted.length })

      // Encrypt in batches of 10
      for (let i = 0; i < unencrypted.length; i += 10) {
        const batch = unencrypted.slice(i, i + 10)
        await Promise.all(
          batch.map(async (note) => {
            const encrypted = await encryptNote(kek, note.title, note.body)
            await api.put(`/api/notes/${note.id}`, {
              title: '',
              body: '',
              encryptedTitle: encrypted.encryptedTitle,
              encryptedBody: encrypted.encryptedBody,
              noteKeyWrapped: encrypted.noteKeyWrapped,
              keyVersion,
              isEncrypted: true,
            })
          }),
        )
        setProgress((p) => ({ ...p, done: Math.min(i + 10, unencrypted.length) }))
      }

      // Refresh notes
      useNotesStore.getState().fetchNotes()
    } finally {
      setMigrating(false)
    }
  }

  return (
    <>
      <Button variant="outline" size="sm" onClick={handleMigrate} disabled={migrating}>
        {migrating
          ? `Encrypting ${progress.done}/${progress.total}...`
          : 'Encrypt existing notes'}
      </Button>
    </>
  )
}
