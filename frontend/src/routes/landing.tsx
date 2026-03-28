import { useNavigate } from 'react-router'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { BookOpen, Shield, Search, Smartphone } from 'lucide-react'

export function LandingPage() {
  const navigate = useNavigate()

  return (
    <div className="min-h-screen bg-background">
      {/* Nav */}
      <nav className="border-b">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
          <div className="flex items-center gap-2">
            <BookOpen className="h-6 w-6 text-primary" />
            <span className="text-lg font-bold">Note Thing</span>
          </div>
          <div className="flex items-center gap-4">
            <button onClick={() => navigate('/features')} className="text-sm text-muted-foreground hover:text-foreground">Features</button>
            <button onClick={() => navigate('/plans')} className="text-sm text-muted-foreground hover:text-foreground">Plans</button>
            <button onClick={() => navigate('/about')} className="text-sm text-muted-foreground hover:text-foreground">About</button>
            <Button variant="outline" size="sm" onClick={() => navigate('/login')}>Sign In</Button>
            <Button size="sm" onClick={() => navigate('/login')}>Get Started</Button>
          </div>
        </div>
      </nav>

      {/* Hero */}
      <section className="mx-auto max-w-4xl px-6 py-24 text-center">
        <h1 className="text-5xl font-bold tracking-tight">
          Your notes, <span className="text-primary">encrypted</span> and organized.
        </h1>
        <p className="mx-auto mt-6 max-w-2xl text-lg text-muted-foreground">
          Note Thing is a secure note-taking app with end-to-end encryption. Your notes are encrypted
          on your device before they ever reach our servers. We can never read your notes.
        </p>
        <div className="mt-8 flex justify-center gap-4">
          <Button size="lg" onClick={() => navigate('/login')}>Start for Free</Button>
          <Button size="lg" variant="outline" onClick={() => navigate('/features')}>Learn More</Button>
        </div>
        <p className="mt-4 text-sm text-muted-foreground">Free plan includes 50 notes. No credit card required.</p>
      </section>

      <Separator />

      {/* Feature highlights */}
      <section className="mx-auto max-w-6xl px-6 py-20">
        <div className="grid gap-12 md:grid-cols-3">
          <div className="text-center">
            <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
              <Shield className="h-6 w-6 text-primary" />
            </div>
            <h3 className="text-lg font-semibold">Zero-Knowledge Encryption</h3>
            <p className="mt-2 text-sm text-muted-foreground">
              AES-256-GCM encryption with Argon2id key derivation. Your master password never leaves your device.
            </p>
          </div>
          <div className="text-center">
            <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
              <Search className="h-6 w-6 text-primary" />
            </div>
            <h3 className="text-lg font-semibold">Full-Text Search</h3>
            <p className="mt-2 text-sm text-muted-foreground">
              Search across all your notes instantly. Client-side search keeps your queries private too.
            </p>
          </div>
          <div className="text-center">
            <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10">
              <Smartphone className="h-6 w-6 text-primary" />
            </div>
            <h3 className="text-lg font-semibold">Web & Mobile</h3>
            <p className="mt-2 text-sm text-muted-foreground">
              Access your notes from anywhere. Available on web, iOS, and Android with seamless sync.
            </p>
          </div>
        </div>
      </section>

      <Separator />

      {/* How it works */}
      <section className="mx-auto max-w-4xl px-6 py-20 text-center">
        <h2 className="text-3xl font-bold">How it works</h2>
        <div className="mt-12 grid gap-8 md:grid-cols-3">
          <div>
            <div className="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-full bg-primary text-primary-foreground font-bold">1</div>
            <h3 className="font-semibold">Create an account</h3>
            <p className="mt-1 text-sm text-muted-foreground">Sign up with email or Google. Free plan, no credit card needed.</p>
          </div>
          <div>
            <div className="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-full bg-primary text-primary-foreground font-bold">2</div>
            <h3 className="font-semibold">Set your master password</h3>
            <p className="mt-1 text-sm text-muted-foreground">Choose a strong password that encrypts all your notes. Only you know it.</p>
          </div>
          <div>
            <div className="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-full bg-primary text-primary-foreground font-bold">3</div>
            <h3 className="font-semibold">Start writing</h3>
            <p className="mt-1 text-sm text-muted-foreground">Create notes in Markdown, organize with notebooks and tags.</p>
          </div>
        </div>
      </section>

      <Separator />

      {/* CTA */}
      <section className="mx-auto max-w-4xl px-6 py-20 text-center">
        <h2 className="text-3xl font-bold">Ready to secure your notes?</h2>
        <p className="mt-4 text-muted-foreground">Join thousands of users who trust Note Thing with their most important ideas.</p>
        <Button size="lg" className="mt-8" onClick={() => navigate('/login')}>Get Started for Free</Button>
      </section>

      {/* Footer */}
      <footer className="border-t">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-6">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <BookOpen className="h-4 w-4" />
            <span>Note Thing</span>
          </div>
          <div className="flex gap-6 text-sm text-muted-foreground">
            <button onClick={() => navigate('/features')} className="hover:text-foreground">Features</button>
            <button onClick={() => navigate('/plans')} className="hover:text-foreground">Plans</button>
            <button onClick={() => navigate('/about')} className="hover:text-foreground">About</button>
          </div>
        </div>
      </footer>
    </div>
  )
}
