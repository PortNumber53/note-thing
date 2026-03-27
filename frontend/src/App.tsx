import { useEffect } from 'react'
import { Outlet, useNavigate } from 'react-router'
import { Sidebar } from '@/components/layout/sidebar'
import { useAuthStore } from '@/stores/auth-store'
import { useCryptoStore } from '@/stores/crypto-store'
import { TooltipProvider } from '@/components/ui/tooltip'
import { UnlockScreen } from '@/components/auth/unlock-screen'

export function AppLayout() {
  const { isAuthenticated, isLoading, init } = useAuthStore()
  const { hasCheckedEncryption, isEncryptionEnabled, isUnlocked, fetchEncryptionStatus } = useCryptoStore()
  const navigate = useNavigate()

  useEffect(() => {
    init()
  }, [init])

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      fetchEncryptionStatus()
    }
  }, [isLoading, isAuthenticated, fetchEncryptionStatus])

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      navigate('/login', { replace: true })
    }
  }, [isLoading, isAuthenticated, navigate])

  // Reset auto-lock timer on user activity
  useEffect(() => {
    const reset = () => useCryptoStore.getState().resetActivity()
    window.addEventListener('mousedown', reset)
    window.addEventListener('keydown', reset)
    return () => {
      window.removeEventListener('mousedown', reset)
      window.removeEventListener('keydown', reset)
    }
  }, [])

  useEffect(() => {
    if (!isLoading && isAuthenticated && hasCheckedEncryption && !isEncryptionEnabled) {
      navigate('/setup-encryption', { replace: true })
    }
  }, [isLoading, isAuthenticated, hasCheckedEncryption, isEncryptionEnabled, navigate])

  if (isLoading || !isAuthenticated) return null
  if (!hasCheckedEncryption) return null
  if (!isEncryptionEnabled) return null

  if (isEncryptionEnabled && !isUnlocked) {
    return <UnlockScreen />
  }

  return (
    <TooltipProvider>
      <div className="flex h-screen overflow-hidden">
        <Sidebar />
        <Outlet />
      </div>
    </TooltipProvider>
  )
}
