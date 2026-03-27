import { create } from 'zustand'
import type { BillingStatus } from '@/types'
import { api } from '@/lib/api'

type BillingState = {
  status: BillingStatus | null
  isLoading: boolean
  fetchStatus: () => Promise<void>
  createCheckout: (successUrl: string, cancelUrl: string) => Promise<string>
  createPortal: (returnUrl: string) => Promise<string>
  cancelSubscription: () => Promise<void>
  reactivateSubscription: () => Promise<void>
}

export const useBillingStore = create<BillingState>()((set) => ({
  status: null,
  isLoading: false,

  fetchStatus: async () => {
    set({ isLoading: true })
    try {
      const status = await api.get<BillingStatus>('/api/billing/status')
      set({ status })
    } finally {
      set({ isLoading: false })
    }
  },

  createCheckout: async (successUrl: string, cancelUrl: string) => {
    const result = await api.post<{ checkoutUrl: string }>('/api/billing/checkout', { successUrl, cancelUrl })
    return result.checkoutUrl
  },

  createPortal: async (returnUrl: string) => {
    const result = await api.post<{ portalUrl: string }>('/api/billing/portal', { returnUrl })
    return result.portalUrl
  },

  cancelSubscription: async () => {
    await api.post('/api/billing/cancel')
  },

  reactivateSubscription: async () => {
    await api.post('/api/billing/reactivate')
  },
}))
