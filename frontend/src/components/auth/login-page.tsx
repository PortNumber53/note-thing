import { useState } from 'react'
import { useNavigate } from 'react-router'
import { BookOpen, Eye, EyeOff } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { useAuthStore } from '@/stores/auth-store'
import { useCryptoStore } from '@/stores/crypto-store'
import { api } from '@/lib/api'
import type { User } from '@/types'

type AuthResponse = { token: string; user: User }

export function LoginPage() {
  const [mode, setMode] = useState<'login' | 'signup'>('login')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [name, setName] = useState('')
  const [error, setError] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const { setToken } = useAuthStore()

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const endpoint = mode === 'signup' ? '/auth/signup' : '/auth/login'
      const body = mode === 'signup'
        ? { email, password, name }
        : { email, password }

      const data = await api.post<AuthResponse>(endpoint, body)
      setToken(data.token)
      localStorage.setItem('note-thing-user', JSON.stringify(data.user))
      useAuthStore.setState({ user: data.user })

      await useCryptoStore.getState().fetchEncryptionStatus()
      const { isEncryptionEnabled } = useCryptoStore.getState()
      navigate(isEncryptionEnabled ? '/notes' : '/setup-encryption', { replace: true })
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Authentication failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="mx-auto w-full max-w-sm space-y-6 p-6">
        <div className="flex flex-col items-center gap-2 text-center">
          <BookOpen className="h-12 w-12 text-primary" />
          <h1 className="text-2xl font-bold">Note Thing</h1>
          <p className="text-muted-foreground">Your notes, organized.</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {mode === 'signup' && (
            <div className="space-y-1.5">
              <Label htmlFor="name">Name</Label>
              <Input
                id="name"
                value={name}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setName(e.target.value)}
                placeholder="Your name"
                required
              />
            </div>
          )}
          <div className="space-y-1.5">
            <Label htmlFor="email">Email</Label>
            <Input
              id="email"
              type="email"
              value={email}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
              placeholder="you@example.com"
              required
            />
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="password">Password</Label>
            <div className="relative">
              <Input
                id="password"
                type={showPassword ? 'text' : 'password'}
                value={password}
                onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
                placeholder={mode === 'signup' ? 'At least 8 characters' : 'Your password'}
                required
                minLength={mode === 'signup' ? 8 : undefined}
                className="pr-10"
              />
              <button
                type="button"
                onClick={() => setShowPassword(!showPassword)}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              >
                {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
              </button>
            </div>
          </div>
          {error && <p className="text-sm text-destructive">{error}</p>}
          <Button type="submit" className="w-full" size="lg" disabled={loading}>
            {loading ? 'Please wait...' : mode === 'signup' ? 'Create Account' : 'Sign In'}
          </Button>
        </form>

        <div className="text-center text-sm">
          {mode === 'login' ? (
            <p className="text-muted-foreground">
              Don&apos;t have an account?{' '}
              <button className="text-primary underline" onClick={() => { setMode('signup'); setError('') }}>
                Sign up
              </button>
            </p>
          ) : (
            <p className="text-muted-foreground">
              Already have an account?{' '}
              <button className="text-primary underline" onClick={() => { setMode('login'); setError('') }}>
                Sign in
              </button>
            </p>
          )}
        </div>

        <div className="relative">
          <Separator />
          <span className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 bg-background px-2 text-xs text-muted-foreground">
            or
          </span>
        </div>

        <Button
          variant="outline"
          className="w-full"
          size="lg"
          onClick={() => { window.location.href = '/auth/google/login' }}
        >
          Continue with Google
        </Button>
      </div>
    </div>
  )
}
