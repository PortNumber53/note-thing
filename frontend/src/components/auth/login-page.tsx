import { Button } from '@/components/ui/button'
import { BookOpen } from 'lucide-react'

export function LoginPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background">
      <div className="mx-auto w-full max-w-sm space-y-6 text-center">
        <div className="flex flex-col items-center gap-2">
          <BookOpen className="h-12 w-12 text-primary" />
          <h1 className="text-2xl font-bold">Note Thing</h1>
          <p className="text-muted-foreground">Your notes, organized.</p>
        </div>
        <Button
          className="w-full"
          size="lg"
          onClick={() => { window.location.href = '/auth/google/login' }}
        >
          Sign in with Google
        </Button>
      </div>
    </div>
  )
}
