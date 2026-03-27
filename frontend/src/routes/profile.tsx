import { useAuthStore } from '@/stores/auth-store'

export function ProfileView() {
  const user = useAuthStore((s) => s.user)

  if (!user) return null

  return (
    <div className="flex-1 p-8">
      <h1 className="text-2xl font-bold mb-6">Profile</h1>
      <div className="max-w-md space-y-6">
        <div className="flex items-center gap-4">
          <img
            src={user.avatarUrl}
            alt={user.name}
            className="h-16 w-16 rounded-full"
          />
          <div>
            <h2 className="text-lg font-semibold">{user.name}</h2>
            <p className="text-sm text-muted-foreground">{user.email}</p>
          </div>
        </div>
      </div>
    </div>
  )
}
