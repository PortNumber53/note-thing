export type User = {
  id: string
  email: string
  name: string
  avatarUrl: string
}

export type Note = {
  id: string
  title: string
  body: string
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
