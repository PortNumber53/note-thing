import { useNavigate } from 'react-router'
import { Button } from '@/components/ui/button'
import { Check } from 'lucide-react'
import logoImg from '@/assets/logo.png'

const freePlan = {
  name: 'Free',
  price: '$0',
  period: 'forever',
  description: 'For personal use and getting started.',
  features: [
    'Up to 50 notes',
    '1 notebook',
    'End-to-end encryption',
    'Markdown editor',
    'Tags & search',
    'Web & mobile apps',
    '1MB max note size',
  ],
  cta: 'Get Started',
}

const proPlan = {
  name: 'Pro',
  price: '$10.99',
  period: '/month',
  description: 'For power users who need more.',
  features: [
    'Unlimited notes',
    'Unlimited notebooks',
    'End-to-end encryption',
    'Markdown editor',
    'Tags & search',
    'Web & mobile apps',
    'Unlimited note size',
    '14-day free trial',
    'Priority support',
  ],
  cta: 'Start Free Trial',
}

export function PlansPage() {
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
            <button onClick={() => navigate('/features')} className="text-sm text-muted-foreground hover:text-foreground">Features</button>
            <button onClick={() => navigate('/plans')} className="text-sm font-medium">Plans</button>
            <button onClick={() => navigate('/about')} className="text-sm text-muted-foreground hover:text-foreground">About</button>
            <Button variant="outline" size="sm" onClick={() => navigate('/login')}>Sign In</Button>
            <Button size="sm" onClick={() => navigate('/login')}>Get Started</Button>
          </div>
        </div>
      </nav>

      <section className="mx-auto max-w-4xl px-6 py-20 text-center">
        <h1 className="text-4xl font-bold">Simple, transparent pricing</h1>
        <p className="mt-4 text-lg text-muted-foreground">
          Start free, upgrade when you need more.
        </p>
      </section>

      <section className="mx-auto max-w-4xl px-6 pb-20">
        <div className="grid gap-8 md:grid-cols-2">
          {/* Free */}
          <div className="rounded-xl border p-8">
            <h3 className="text-xl font-semibold">{freePlan.name}</h3>
            <div className="mt-4 flex items-baseline gap-1">
              <span className="text-4xl font-bold">{freePlan.price}</span>
              <span className="text-muted-foreground">/{freePlan.period}</span>
            </div>
            <p className="mt-2 text-sm text-muted-foreground">{freePlan.description}</p>
            <Button variant="outline" className="mt-6 w-full" size="lg" onClick={() => navigate('/login')}>
              {freePlan.cta}
            </Button>
            <ul className="mt-6 space-y-3">
              {freePlan.features.map((f) => (
                <li key={f} className="flex items-center gap-2 text-sm">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  {f}
                </li>
              ))}
            </ul>
          </div>

          {/* Pro */}
          <div className="rounded-xl border-2 border-primary p-8 relative">
            <div className="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-primary px-3 py-0.5 text-xs font-medium text-primary-foreground">
              Most popular
            </div>
            <h3 className="text-xl font-semibold">{proPlan.name}</h3>
            <div className="mt-4 flex items-baseline gap-1">
              <span className="text-4xl font-bold">{proPlan.price}</span>
              <span className="text-muted-foreground">{proPlan.period}</span>
            </div>
            <p className="mt-2 text-sm text-muted-foreground">{proPlan.description}</p>
            <Button className="mt-6 w-full" size="lg" onClick={() => navigate('/login')}>
              {proPlan.cta}
            </Button>
            <ul className="mt-6 space-y-3">
              {proPlan.features.map((f) => (
                <li key={f} className="flex items-center gap-2 text-sm">
                  <Check className="h-4 w-4 text-primary shrink-0" />
                  {f}
                </li>
              ))}
            </ul>
          </div>
        </div>
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
