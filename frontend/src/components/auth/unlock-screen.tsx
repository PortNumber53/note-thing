import { useState } from 'react'
import { Lock } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useCryptoStore } from '@/stores/crypto-store'
import { useAuthStore } from '@/stores/auth-store'

export function UnlockScreen() {
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const unlock = useCryptoStore((s) => s.unlock)
  const logout = useAuthStore((s) => s.logout)

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (!password) return
    setLoading(true)
    setError('')
    try {
      const success = await unlock(password)
      if (!success) {
        setError('Incorrect master password')
        setPassword('')
      }
    } catch {
      setError('Decryption failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="mx-auto w-full max-w-sm space-y-6 text-center">
        <div className="flex flex-col items-center gap-2">
          <Lock className="h-12 w-12 text-primary" />
          <h1 className="text-2xl font-bold">Vault Locked</h1>
          <p className="text-muted-foreground">Enter your master password to unlock your notes.</p>
        </div>
        <form onSubmit={handleSubmit} className="space-y-4">
          <Input
            type="password"
            value={password}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
            placeholder="Master password"
            autoFocus
            disabled={loading}
          />
          {error && <p className="text-sm text-destructive">{error}</p>}
          <Button type="submit" className="w-full" disabled={loading || !password}>
            {loading ? 'Unlocking...' : 'Unlock'}
          </Button>
        </form>
        <Button variant="ghost" className="w-full text-muted-foreground" onClick={logout}>
          Sign out
        </Button>
      </div>
    </div>
  )
}
