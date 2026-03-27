export type User = {
  id: string
  email: string
  name: string
  avatarUrl: string
  createdAt: string
  updatedAt: string
}

export type Note = {
  id: string
  title: string
  body: string
  encryptedTitle?: string
  encryptedBody?: string
  noteKeyWrapped?: string
  keyVersion?: number
  isEncrypted: boolean
  notebookId: string | null
  tags: Tag[]
  createdAt: string
  updatedAt: string
}

export type Notebook = {
  id: string
  name: string
  isDefault: boolean
  noteCount: number
  createdAt: string
  updatedAt: string
}

export type Tag = {
  id: string
  name: string
}

export type UserSettings = {
  defaultNotebookId: string | null
}

export type EditorPreviewMode = 'split' | 'live' | 'source'

export type BillingSubscription = {
  id: string
  status: string
  stripePriceId: string
  trialEnd: string | null
  currentPeriodEnd: string | null
  cancelAtPeriodEnd: boolean
  amountCents: number
  currency: string
}

export type BillingStatus = {
  subscription: BillingSubscription | null
  hasActiveAccess: boolean
}
