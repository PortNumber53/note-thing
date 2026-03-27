import MiniSearch from 'minisearch'

const index = new MiniSearch<{ id: string; title: string; body: string }>({
  fields: ['title', 'body'],
  storeFields: ['title'],
  searchOptions: {
    boost: { title: 2 },
    fuzzy: 0.2,
    prefix: true,
  },
})

export function rebuildIndex(notes: Array<{ id: string; title: string; body: string }>) {
  index.removeAll()
  index.addAll(notes)
}

export function addToIndex(note: { id: string; title: string; body: string }) {
  if (index.has(note.id)) {
    index.discard(note.id)
  }
  index.add(note)
}

export function removeFromIndex(id: string) {
  if (index.has(id)) {
    index.discard(id)
  }
}

export function searchNotes(query: string): string[] {
  if (!query.trim()) return []
  return index.search(query).map((r) => r.id)
}
