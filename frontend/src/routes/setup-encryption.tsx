import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router'
import { Shield, AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useAuthStore } from '@/stores/auth-store'
import { useCryptoStore } from '@/stores/crypto-store'

export function SetupEncryptionPage() {
  const navigate = useNavigate()
  const { isAuthenticated, init } = useAuthStore()

  useEffect(() => {
    init()
  }, [init])

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login', { replace: true })
    }
  }, [isAuthenticated, navigate])
  const setupEncryption = useCryptoStore((s) => s.setupEncryption)
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
      await setupEncryption(password)
      navigate('/app/notes', { replace: true })
    } catch {
      setError('Setup failed. Please try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="mx-auto w-full max-w-md space-y-6 p-6">
        <div className="flex flex-col items-center gap-3 text-center">
          <Shield className="h-12 w-12 text-primary" />
          <h1 className="text-2xl font-bold">Secure Your Notes</h1>
          <p className="text-muted-foreground">
            Note Thing uses end-to-end encryption. Your notes are encrypted on this device
            before being sent to the server. We can never read your notes.
          </p>
        </div>

        <div className="rounded-md bg-yellow-50 dark:bg-yellow-950 p-3 flex gap-2">
          <AlertTriangle className="h-4 w-4 text-yellow-600 shrink-0 mt-0.5" />
          <p className="text-sm text-yellow-800 dark:text-yellow-200">
            Choose a strong master password. If you forget it, your encrypted notes
            <strong> cannot be recovered</strong>.
          </p>
        </div>

        <div className="space-y-4">
          <div className="space-y-1.5">
            <Label htmlFor="master-pw">Master password</Label>
            <Input
              id="master-pw"
              type="password"
              value={password}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
              placeholder="At least 8 characters"
              autoFocus
            />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="confirm-pw">Confirm password</Label>
            <Input
              id="confirm-pw"
              type="password"
              value={confirm}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setConfirm(e.target.value)}
              placeholder="Repeat password"
            />
          </div>
          {error && <p className="text-sm text-destructive">{error}</p>}
          <Button
            className="w-full"
            size="lg"
            onClick={handleSetup}
            disabled={loading || !password || !confirm}
          >
            {loading ? 'Setting up encryption...' : 'Create Master Password'}
          </Button>
        </div>
      </div>
    </div>
  )
}
