import { useNavigate } from 'react-router'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Shield, Lock, Search, Tag, Trash2, Moon, Smartphone, FileText, FolderOpen, Zap, Key, BookOpen } from 'lucide-react'
import logoImg from '@/assets/logo.png'

const features = [
  {
    icon: Shield,
    title: 'End-to-End Encryption',
    description: 'Notes are encrypted on your device using AES-256-GCM before being sent to our servers. We use Argon2id for key derivation. We literally cannot read your notes.',
  },
  {
    icon: Key,
    title: 'Per-Note Encryption Keys',
    description: 'Each note has its own encryption key, wrapped by your master key. Compromising one note does not compromise others.',
  },
  {
    icon: FileText,
    title: 'Markdown Editor',
    description: 'Write in Markdown with a live preview editor. Bold, italic, headings, code blocks, links, and more.',
  },
  {
    icon: FolderOpen,
    title: 'Notebooks',
    description: 'Organize your notes into notebooks. Create as many as you need with the Pro plan.',
  },
  {
    icon: Tag,
    title: 'Tags',
    description: 'Add tags to your notes for flexible cross-notebook organization. Filter notes by tag instantly.',
  },
  {
    icon: Search,
    title: 'Client-Side Search',
    description: 'Full-text search runs entirely on your device after decryption. Your search queries never reach our servers.',
  },
  {
    icon: Trash2,
    title: 'Trash & Recovery',
    description: 'Deleted notes go to trash first. Restore them anytime or permanently delete when you are ready.',
  },
  {
    icon: Smartphone,
    title: 'Mobile App',
    description: 'Native Android and iOS apps with the same encryption. Your notes sync seamlessly across all your devices.',
  },
  {
    icon: Moon,
    title: 'Dark Mode',
    description: 'System, light, or dark theme. The editor follows your preference for comfortable writing at any time.',
  },
  {
    icon: Lock,
    title: 'Auto-Lock',
    description: 'Your vault auto-locks after 15 minutes of inactivity. Keys are cleared from memory for maximum security.',
  },
  {
    icon: Zap,
    title: 'Fast & Lightweight',
    description: 'Built with performance in mind. AES-256-GCM runs with hardware acceleration via Web Crypto API.',
  },
  {
    icon: BookOpen,
    title: 'Open API',
    description: 'RESTful API with JWT authentication. Build your own integrations and workflows.',
  },
]

export function FeaturesPage() {
  const navigate = useNavigate()

  return (
    <div className="min-h-screen bg-background">
      <nav className="border-b">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
          <button onClick={() => navigate('/')} className="flex items-center gap-2">
            <img src={logoImg} alt="Note Thing" className="h-6 w-6 rounded" />
            <span className="text-lg font-bold">Note Thing</span>
          </button>
          <div className="flex items-center gap-4">
            <button onClick={() => navigate('/features')} className="text-sm font-medium">Features</button>
            <button onClick={() => navigate('/plans')} className="text-sm text-muted-foreground hover:text-foreground">Plans</button>
            <button onClick={() => navigate('/about')} className="text-sm text-muted-foreground hover:text-foreground">About</button>
            <Button variant="outline" size="sm" onClick={() => navigate('/login')}>Sign In</Button>
            <Button size="sm" onClick={() => navigate('/login')}>Get Started</Button>
          </div>
        </div>
      </nav>

      <section className="mx-auto max-w-4xl px-6 py-20 text-center">
        <h1 className="text-4xl font-bold">Features</h1>
        <p className="mt-4 text-lg text-muted-foreground">
          Everything you need for secure, organized note-taking.
        </p>
      </section>

      <section className="mx-auto max-w-6xl px-6 pb-20">
        <div className="grid gap-8 md:grid-cols-2 lg:grid-cols-3">
          {features.map((f) => (
            <div key={f.title} className="rounded-lg border p-6">
              <f.icon className="h-8 w-8 text-primary" />
              <h3 className="mt-4 text-lg font-semibold">{f.title}</h3>
              <p className="mt-2 text-sm text-muted-foreground">{f.description}</p>
            </div>
          ))}
        </div>
      </section>

      <Separator />

      <section className="mx-auto max-w-4xl px-6 py-20 text-center">
        <h2 className="text-3xl font-bold">Ready to try it?</h2>
        <p className="mt-4 text-muted-foreground">Start with the free plan. No credit card required.</p>
        <Button size="lg" className="mt-8" onClick={() => navigate('/login')}>Get Started for Free</Button>
      </section>

      <footer className="border-t">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-6">
          <div className="flex items-center gap-2 text-sm text-muted-foreground">
            <img src={logoImg} alt="Note Thing" className="h-4 w-4 rounded" />
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
