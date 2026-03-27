import { useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router'
import { useAuthStore } from '@/stores/auth-store'

export function AuthCallback() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const setToken = useAuthStore((s) => s.setToken)
  const fetchUser = useAuthStore((s) => s.fetchUser)

  useEffect(() => {
    const token = searchParams.get('token')
    if (token) {
      setToken(token)
      fetchUser().then(() => navigate('/notes', { replace: true }))
    } else {
      navigate('/login', { replace: true })
    }
  }, [searchParams, setToken, fetchUser, navigate])

  return (
    <div className="flex min-h-screen items-center justify-center">
      <p className="text-muted-foreground">Signing you in...</p>
    </div>
  )
}
