import { argon2id } from 'hash-wasm'

const VERIFY_PLAINTEXT = 'note-thing-verify'
const IV_LENGTH = 12

// --- Key Derivation ---

export async function deriveKEK(password: string, salt: Uint8Array): Promise<CryptoKey> {
  const hash = await argon2id({
    password,
    salt,
    iterations: 3,
    memorySize: 65536, // 64 MB
    parallelism: 4,
    hashLength: 32,
    outputType: 'binary',
  })

  return crypto.subtle.importKey(
    'raw',
    hash,
    { name: 'AES-GCM' },
    false,
    ['encrypt', 'decrypt', 'wrapKey', 'unwrapKey'],
  )
}

// --- Note Key Management ---

export async function generateNoteKey(): Promise<CryptoKey> {
  return crypto.subtle.generateKey(
    { name: 'AES-GCM', length: 256 },
    true, // extractable so we can wrap it
    ['encrypt', 'decrypt'],
  )
}

export async function wrapKey(dek: CryptoKey, kek: CryptoKey): Promise<Uint8Array> {
  const iv = crypto.getRandomValues(new Uint8Array(IV_LENGTH))
  const wrapped = await crypto.subtle.wrapKey('raw', dek, kek, { name: 'AES-GCM', iv })
  return concat(iv, new Uint8Array(wrapped))
}

export async function unwrapKey(data: Uint8Array, kek: CryptoKey): Promise<CryptoKey> {
  const iv = data.slice(0, IV_LENGTH)
  const wrapped = data.slice(IV_LENGTH)
  return crypto.subtle.unwrapKey(
    'raw',
    wrapped,
    kek,
    { name: 'AES-GCM', iv },
    { name: 'AES-GCM' },
    true,
    ['encrypt', 'decrypt'],
  )
}

// --- Field Encryption ---

export async function encryptField(key: CryptoKey, plaintext: string): Promise<Uint8Array> {
  const iv = crypto.getRandomValues(new Uint8Array(IV_LENGTH))
  const encoded = new TextEncoder().encode(plaintext)
  const ciphertext = await crypto.subtle.encrypt({ name: 'AES-GCM', iv }, key, encoded)
  return concat(iv, new Uint8Array(ciphertext))
}

export async function decryptField(key: CryptoKey, data: Uint8Array): Promise<string> {
  const iv = data.slice(0, IV_LENGTH)
  const ciphertext = data.slice(IV_LENGTH)
  const plaintext = await crypto.subtle.decrypt({ name: 'AES-GCM', iv }, key, ciphertext)
  return new TextDecoder().decode(plaintext)
}

// --- Note-Level Operations ---

export type EncryptedNotePayload = {
  encryptedTitle: string  // base64
  encryptedBody: string   // base64
  noteKeyWrapped: string  // base64
}

export type DecryptedNoteContent = {
  title: string
  body: string
  dek: CryptoKey
}

export async function encryptNote(
  kek: CryptoKey,
  title: string,
  body: string,
  existingDEK?: CryptoKey,
): Promise<EncryptedNotePayload & { dek: CryptoKey }> {
  const dek = existingDEK || await generateNoteKey()
  const [encTitle, encBody, wrappedKey] = await Promise.all([
    encryptField(dek, title),
    encryptField(dek, body),
    wrapKey(dek, kek),
  ])
  return {
    encryptedTitle: toBase64(encTitle),
    encryptedBody: toBase64(encBody),
    noteKeyWrapped: toBase64(wrappedKey),
    dek,
  }
}

export async function decryptNote(
  kek: CryptoKey,
  encryptedTitle: string,
  encryptedBody: string,
  noteKeyWrapped: string,
): Promise<DecryptedNoteContent> {
  const dek = await unwrapKey(fromBase64(noteKeyWrapped), kek)
  const [title, body] = await Promise.all([
    decryptField(dek, fromBase64(encryptedTitle)),
    decryptField(dek, fromBase64(encryptedBody)),
  ])
  return { title, body, dek }
}

// --- Password Verification ---

export async function createVerifyToken(kek: CryptoKey): Promise<string> {
  const encrypted = await encryptField(kek, VERIFY_PLAINTEXT)
  return toBase64(encrypted)
}

export async function verifyPassword(kek: CryptoKey, token: string): Promise<boolean> {
  try {
    const decrypted = await decryptField(kek, fromBase64(token))
    return decrypted === VERIFY_PLAINTEXT
  } catch {
    return false
  }
}

// --- Helpers ---

export function generateSalt(): Uint8Array {
  return crypto.getRandomValues(new Uint8Array(16))
}

export function toBase64(data: Uint8Array): string {
  let binary = ''
  for (let i = 0; i < data.length; i++) {
    binary += String.fromCharCode(data[i])
  }
  return btoa(binary)
}

export function fromBase64(b64: string): Uint8Array {
  const binary = atob(b64)
  const bytes = new Uint8Array(binary.length)
  for (let i = 0; i < binary.length; i++) {
    bytes[i] = binary.charCodeAt(i)
  }
  return bytes
}

function concat(a: Uint8Array, b: Uint8Array): Uint8Array {
  const result = new Uint8Array(a.length + b.length)
  result.set(a, 0)
  result.set(b, a.length)
  return result
}
