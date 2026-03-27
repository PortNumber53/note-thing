import { createBrowserRouter, Navigate } from 'react-router'
import { AppLayout } from '@/app'
import { LoginPage } from '@/components/auth/login-page'
import { AuthCallback } from '@/components/auth/auth-callback'
import { AllNotesView } from './all-notes'
import { NotebookNotesView } from './notebook-notes'
import { TagNotesView } from './tag-notes'
import { SearchView } from './search-results'
import { TrashView } from './trash'
import { ProfileView } from './profile'
import { SettingsView } from './settings'
import { BillingView } from './billing'

export const router = createBrowserRouter([
  { path: '/login', element: <LoginPage /> },
  { path: '/auth/callback', element: <AuthCallback /> },
  {
    path: '/',
    element: <AppLayout />,
    children: [
      { index: true, element: <Navigate to="/notes" replace /> },
      { path: 'notes', element: <AllNotesView /> },
      { path: 'notebooks/:notebookId', element: <NotebookNotesView /> },
      { path: 'tags/:tagId', element: <TagNotesView /> },
      { path: 'search', element: <SearchView /> },
      { path: 'trash', element: <TrashView /> },
      { path: 'account/profile', element: <ProfileView /> },
      { path: 'account/settings', element: <SettingsView /> },
      { path: 'account/billing', element: <BillingView /> },
    ],
  },
])
