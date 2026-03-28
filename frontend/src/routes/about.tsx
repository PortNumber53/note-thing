import { useNavigate } from 'react-router'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { BookOpen, Shield, Heart, Globe } from 'lucide-react'

export function AboutPage() {
  const navigate = useNavigate()

  return (
    <div className="min-h-screen bg-background">
      <nav className="border-b">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
          <button onClick={() => navigate('/')} className="flex items-center gap-2">
            <BookOpen className="h-6 w-6 text-primary" />
            <span className="text-lg font-bold">Note Thing</span>
          </button>
          <div className="flex items-center gap-4">
            <button onClick={() => navigate('/features')} className="text-sm text-muted-foreground hover:text-foreground">Features</button>
            <button onClick={() => navigate('/plans')} className="text-sm text-muted-foreground hover:text-foreground">Plans</button>
            <button onClick={() => navigate('/about')} className="text-sm font-medium">About</button>
            <Button variant="outline" size="sm" onClick={() => navigate('/login')}>Sign In</Button>
            <Button size="sm" onClick={() => navigate('/login')}>Get Started</Button>
          </div>
        </div>
      </nav>

      <section className="mx-auto max-w-3xl px-6 py-20">
        <h1 className="text-4xl font-bold">About Note Thing</h1>

        <div className="mt-12 space-y-12">
          <div>
            <div className="flex items-center gap-3">
              <Shield className="h-6 w-6 text-primary" />
              <h2 className="text-2xl font-semibold">Privacy first</h2>
            </div>
            <p className="mt-4 text-muted-foreground leading-relaxed">
              We built Note Thing because we believe your thoughts are yours alone. Every note is encrypted
              on your device before it reaches our servers using AES-256-GCM with Argon2id key derivation.
              We have zero knowledge of your content — not even we can read your notes.
            </p>
          </div>

          <Separator />

          <div>
            <div className="flex items-center gap-3">
              <Heart className="h-6 w-6 text-primary" />
              <h2 className="text-2xl font-semibold">Built with care</h2>
            </div>
            <p className="mt-4 text-muted-foreground leading-relaxed">
              Note Thing is designed to be simple, fast, and reliable. We use battle-tested cryptographic
              libraries and follow security best practices. Our encryption protocol is inspired by
              Notesnook and Standard Notes — proven approaches that have been audited and trusted by
              millions of users.
            </p>
          </div>

          <Separator />

          <div>
            <div className="flex items-center gap-3">
              <Globe className="h-6 w-6 text-primary" />
              <h2 className="text-2xl font-semibold">Available everywhere</h2>
            </div>
            <p className="mt-4 text-muted-foreground leading-relaxed">
              Access your notes from any device. Our web app works on any modern browser, and our
              native mobile apps for Android and iOS provide a seamless experience. All your notes
              sync automatically across devices, encrypted end-to-end.
            </p>
          </div>

          <Separator />

          <div>
            <h2 className="text-2xl font-semibold">Our technology</h2>
            <ul className="mt-4 space-y-2 text-muted-foreground">
              <li><strong className="text-foreground">Encryption:</strong> AES-256-GCM with per-note keys, wrapped by Argon2id-derived master key</li>
              <li><strong className="text-foreground">Backend:</strong> Go with PostgreSQL, deployed on dedicated infrastructure</li>
              <li><strong className="text-foreground">Frontend:</strong> React with TypeScript, deployed on Cloudflare Workers</li>
              <li><strong className="text-foreground">Mobile:</strong> Flutter for iOS and Android with platform-native crypto</li>
              <li><strong className="text-foreground">Search:</strong> Client-side full-text search — your queries stay on your device</li>
            </ul>
          </div>
        </div>
      </section>

      <Separator />

      <section className="mx-auto max-w-4xl px-6 py-20 text-center">
        <h2 className="text-3xl font-bold">Try Note Thing today</h2>
        <p className="mt-4 text-muted-foreground">Start with the free plan. No credit card required.</p>
        <Button size="lg" className="mt-8" onClick={() => navigate('/login')}>Get Started for Free</Button>
      </section>

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
