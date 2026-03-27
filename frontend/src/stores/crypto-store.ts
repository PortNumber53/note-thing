import { create } from 'zustand'
import { api } from '@/lib/api'
import {
  deriveKEK,
  verifyPassword,
  createVerifyToken,
  generateSalt,
  toBase64,
  fromBase64,
} from '@/lib/crypto'

type CryptoState = {
  hasCheckedEncryption: boolean
  isEncryptionEnabled: boolean
  isUnlocked: boolean
  kek: CryptoKey | null
  keyVersion: number | null
  kdfSalt: Uint8Array | null
  kekVerifyToken: string | null
  noteKeyCache: Map<string, CryptoKey>
  autoLockTimeout: number // minutes
  _lockTimer: ReturnType<typeof setTimeout> | null

  fetchEncryptionStatus: () => Promise<void>
  unlock: (password: string) => Promise<boolean>
  lock: () => void
  setupEncryption: (password: string) => Promise<void>
  cacheNoteKey: (noteId: string, dek: CryptoKey) => void
  getCachedNoteKey: (noteId: string) => CryptoKey | undefined
  resetActivity: () => void
}

export const useCryptoStore = create<CryptoState>()((set, get) => ({
  hasCheckedEncryption: false,
  isEncryptionEnabled: false,
  isUnlocked: false,
  kek: null,
  keyVersion: null,
  kdfSalt: null,
  kekVerifyToken: null,
  noteKeyCache: new Map(),
  autoLockTimeout: 15,
  _lockTimer: null,

  fetchEncryptionStatus: async () => {
    try {
      const data = await api.get<{
        enabled: boolean
        kdfSalt?: string
        keyVersion?: number
        kekVerify?: string
      }>('/api/encryption')

      if (data.enabled && data.kdfSalt && data.kekVerify) {
        set({
          hasCheckedEncryption: true,
          isEncryptionEnabled: true,
          kdfSalt: fromBase64(data.kdfSalt),
          keyVersion: data.keyVersion ?? 1,
          kekVerifyToken: data.kekVerify,
        })
      } else {
        set({ hasCheckedEncryption: true, isEncryptionEnabled: false })
      }
    } catch {
      set({ hasCheckedEncryption: true, isEncryptionEnabled: false })
    }
  },

  unlock: async (password: string) => {
    const { kdfSalt, kekVerifyToken } = get()
    if (!kdfSalt || !kekVerifyToken) return false

    try {
      const kek = await deriveKEK(password, kdfSalt)
      const valid = await verifyPassword(kek, kekVerifyToken)
      if (!valid) return false

      set({ kek, isUnlocked: true })
      get().resetActivity()
      return true
    } catch {
      return false
    }
  },

  lock: () => {
    const timer = get()._lockTimer
    if (timer) clearTimeout(timer)
    set({
      kek: null,
      isUnlocked: false,
      noteKeyCache: new Map(),
      _lockTimer: null,
    })
  },

  setupEncryption: async (password: string) => {
    const salt = generateSalt()
    const kek = await deriveKEK(password, salt)
    const verifyToken = await createVerifyToken(kek)

    await api.post('/api/encryption/setup', {
      kdfSalt: toBase64(salt),
      kekVerify: verifyToken,
    })

    set({
      isEncryptionEnabled: true,
      isUnlocked: true,
      kek,
      keyVersion: 1,
      kdfSalt: salt,
      kekVerifyToken: verifyToken,
    })
    get().resetActivity()
  },

  cacheNoteKey: (noteId: string, dek: CryptoKey) => {
    get().noteKeyCache.set(noteId, dek)
  },

  getCachedNoteKey: (noteId: string) => {
    return get().noteKeyCache.get(noteId)
  },

  resetActivity: () => {
    const timer = get()._lockTimer
    if (timer) clearTimeout(timer)
    const timeout = get().autoLockTimeout
    if (timeout > 0) {
      const newTimer = setTimeout(() => get().lock(), timeout * 60 * 1000)
      set({ _lockTimer: newTimer })
    }
  },
}))
