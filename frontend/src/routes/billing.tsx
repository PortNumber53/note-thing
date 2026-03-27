import { useEffect, useState } from 'react'
import { CreditCard, ExternalLink, AlertCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { useBillingStore } from '@/stores/billing-store'
import { api } from '@/lib/api'
import type { BillingStatus } from '@/types'

export function BillingView() {
  const { status, isLoading, fetchStatus, createCheckout, createPortal, cancelSubscription, reactivateSubscription } = useBillingStore()

  useEffect(() => {
    fetchStatus()
  }, [fetchStatus])

  if (isLoading || !status) {
    return (
      <div className="flex-1 p-8">
        <h1 className="text-2xl font-bold mb-6">Billing</h1>
        <p className="text-muted-foreground">Loading...</p>
      </div>
    )
  }

  return (
    <div className="flex-1 overflow-auto p-8">
      <h1 className="text-2xl font-bold mb-6">Billing</h1>
      <div className="max-w-lg space-y-8">
        {status.subscription ? (
          <SubscriptionInfo
            status={status}
            onCancel={async () => { await cancelSubscription(); fetchStatus() }}
            onReactivate={async () => { await reactivateSubscription(); fetchStatus() }}
            onManage={async () => {
              const url = await createPortal(window.location.href)
              window.location.href = url
            }}
          />
        ) : (
          <NoSubscription
            onSubscribe={async () => {
              const url = await createCheckout(
                `${window.location.origin}/account/billing?success=true`,
                `${window.location.origin}/account/billing?canceled=true`,
              )
              window.location.href = url
            }}
          />
        )}

        <Separator />
        <AdminSection />
      </div>
    </div>
  )
}

function NoSubscription({ onSubscribe }: { onSubscribe: () => void }) {
  return (
    <div className="rounded-lg border p-6 space-y-4">
      <div className="flex items-center gap-3">
        <CreditCard className="h-8 w-8 text-primary" />
        <div>
          <h2 className="text-lg font-semibold">Note Thing Pro</h2>
          <p className="text-sm text-muted-foreground">$10.99/month &middot; 14-day free trial</p>
        </div>
      </div>
      <ul className="text-sm space-y-1 text-muted-foreground">
        <li>Unlimited notes and notebooks</li>
        <li>Full-text search</li>
        <li>Tags and organization</li>
      </ul>
      <Button className="w-full" size="lg" onClick={onSubscribe}>
        Start Free Trial
      </Button>
    </div>
  )
}

function SubscriptionInfo({
  status,
  onCancel,
  onReactivate,
  onManage,
}: {
  status: BillingStatus
  onCancel: () => void
  onReactivate: () => void
  onManage: () => void
}) {
  const sub = status.subscription!
  const isTrialing = sub.status === 'trialing'
  const isCanceling = sub.cancelAtPeriodEnd

  const statusLabel = isCanceling ? 'Canceling' : isTrialing ? 'Trial' : sub.status === 'active' ? 'Active' : sub.status
  const statusColor = isCanceling ? 'text-yellow-600' : (sub.status === 'active' || isTrialing) ? 'text-green-600' : 'text-red-600'

  const periodEnd = sub.currentPeriodEnd || sub.trialEnd
  const endDate = periodEnd ? new Date(periodEnd).toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' }) : null

  let trialDaysLeft: number | null = null
  if (isTrialing && sub.trialEnd) {
    trialDaysLeft = Math.max(0, Math.ceil((new Date(sub.trialEnd).getTime() - Date.now()) / (1000 * 60 * 60 * 24)))
  }

  return (
    <div className="rounded-lg border p-6 space-y-4">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-lg font-semibold">Note Thing Pro</h2>
          <p className="text-sm text-muted-foreground">${(sub.amountCents / 100).toFixed(2)}/{sub.currency?.toUpperCase() === 'USD' ? 'mo' : sub.currency}</p>
        </div>
        <span className={`text-sm font-medium ${statusColor}`}>{statusLabel}</span>
      </div>

      {isTrialing && trialDaysLeft !== null && (
        <div className="flex items-center gap-2 rounded-md bg-secondary p-3">
          <AlertCircle className="h-4 w-4 text-muted-foreground" />
          <span className="text-sm">{trialDaysLeft} days left in trial</span>
        </div>
      )}

      {isCanceling && endDate && (
        <p className="text-sm text-muted-foreground">Your subscription will end on {endDate}</p>
      )}

      <div className="flex gap-2">
        <Button variant="outline" size="sm" onClick={onManage}>
          <ExternalLink className="mr-1.5 h-3.5 w-3.5" />
          Manage in Stripe
        </Button>
        {isCanceling ? (
          <Button variant="outline" size="sm" onClick={onReactivate}>
            Reactivate
          </Button>
        ) : (
          <Button variant="ghost" size="sm" className="text-destructive" onClick={onCancel}>
            Cancel
          </Button>
        )}
      </div>
    </div>
  )
}

function AdminSection() {
  const [isAdmin, setIsAdmin] = useState(false)
  const [newPrice, setNewPrice] = useState('')
  const [migration, setMigration] = useState<Record<string, unknown> | null>(null)
  const [changing, setChanging] = useState(false)

  useEffect(() => {
    // Check if admin by trying to fetch migration status
    api.get<{ migration: Record<string, unknown> | null }>('/api/admin/billing/migration')
      .then((data) => {
        setIsAdmin(true)
        setMigration(data.migration)
      })
      .catch(() => setIsAdmin(false))
  }, [])

  if (!isAdmin) return null

  const handleChangePrice = async () => {
    const cents = Math.round(parseFloat(newPrice) * 100)
    if (isNaN(cents) || cents <= 0) return
    setChanging(true)
    try {
      const result = await api.post<{ migration: Record<string, unknown> }>('/api/admin/billing/price', { amountCents: cents })
      setMigration(result.migration)
      setNewPrice('')
    } finally {
      setChanging(false)
    }
  }

  return (
    <div className="space-y-4">
      <h2 className="text-sm font-medium uppercase text-muted-foreground">Admin: Price Management</h2>

      {migration && (
        <div className="rounded-md bg-secondary p-3 text-sm space-y-1">
          <p>Migration: <strong>{migration.status as string}</strong></p>
          {(migration.totalSubs as number) > 0 && (
            <p>{migration.migratedSubs as number}/{migration.totalSubs as number} migrated, {migration.failedSubs as number} failed</p>
          )}
        </div>
      )}

      <div className="flex items-end gap-2">
        <div className="space-y-1">
          <Label>New monthly price ($)</Label>
          <Input
            value={newPrice}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => setNewPrice(e.target.value)}
            placeholder="10.99"
            type="number"
            step="0.01"
            className="w-32"
          />
        </div>
        <Button onClick={handleChangePrice} disabled={changing || !newPrice}>
          {changing ? 'Changing...' : 'Change Price'}
        </Button>
      </div>
    </div>
  )
}
