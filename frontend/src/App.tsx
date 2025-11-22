import { useEffect, useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import cloudflareLogo from './assets/Cloudflare_Logo.svg'
import './App.css'

type Note = {
  id: string
  title: string
  body: string
  createdAt: string
}

function App() {
  const [count, setCount] = useState(0)
  const [name, setName] = useState('unknown')
  const [notes, setNotes] = useState<Note[]>([])
  const [isLoadingNotes, setIsLoadingNotes] = useState(true)
  const [notesError, setNotesError] = useState<string | null>(null)

  useEffect(() => {
    setIsLoadingNotes(true)
    setNotesError(null)

    fetch('/api/notes')
      .then((res) => {
        if (!res.ok) {
          throw new Error(`failed to load notes (${res.status})`)
        }
        return res.json() as Promise<Note[]>
      })
      .then((data) => setNotes(data))
      .catch((err: unknown) => {
        const message = err instanceof Error ? err.message : 'unknown error'
        setNotesError(message)
      })
      .finally(() => setIsLoadingNotes(false))
  }, [])

  return (
    <>
      <div>
        <a href='https://vite.dev' target='_blank'>
          <img src={viteLogo} className='logo' alt='Vite logo' />
        </a>
        <a href='https://react.dev' target='_blank'>
          <img src={reactLogo} className='logo react' alt='React logo' />
        </a>
        <a href='https://workers.cloudflare.com/' target='_blank'>
          <img src={cloudflareLogo} className='logo cloudflare' alt='Cloudflare logo' />
        </a>
      </div>
      <h1>Vite + React + Cloudflare</h1>
      <div className='card'>
        <button
          onClick={() => setCount((count) => count + 1)}
          aria-label='increment'
        >
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <div className='card'>
        <button
          onClick={() => {
            fetch('/api/')
              .then((res) => res.json() as Promise<{ name: string }>)
              .then((data) => setName(data.name))
          }}
          aria-label='get name'
        >
          Name from API is: {name}
        </button>
        <p>
          Edit <code>worker/index.ts</code> to change the name
        </p>
      </div>
      <p className='read-the-docs'>
        Click on the Vite and React logos to learn more
      </p>

      <div className='card'>
        <h2>Notes</h2>
        {isLoadingNotes && <p>Loading notes…</p>}
        {notesError && <p>Error: {notesError}</p>}
        {!isLoadingNotes && !notesError && notes.length === 0 && <p>No notes yet.</p>}
        {!isLoadingNotes && !notesError && notes.length > 0 && (
          <ul style={{ textAlign: 'left' }}>
            {notes.map((note) => (
              <li key={note.id} style={{ marginBottom: 12 }}>
                <strong>{note.title}</strong>
                <div>{note.body}</div>
                <small>{new Date(note.createdAt).toLocaleString()}</small>
              </li>
            ))}
          </ul>
        )}
      </div>
    </>
  )
}

export default App
