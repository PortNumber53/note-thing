import { createBrowserRouter, Navigate } from 'react-router'
import { AppLayout } from '@/App'
import { LoginPage } from '@/components/auth/login-page'
import { AuthCallback } from '@/components/auth/auth-callback'
import { LandingPage } from './landing'
import { FeaturesPage } from './features'
import { PlansPage } from './plans'
import { AboutPage } from './about'
import { AllNotesView } from './all-notes'
import { NotebookNotesView } from './notebook-notes'
import { TagNotesView } from './tag-notes'
import { SearchView } from './search-results'
import { TrashView } from './trash'
import { ProfileView } from './profile'
import { SettingsView } from './settings'
import { BillingView } from './billing'
import { SetupEncryptionPage } from './setup-encryption'

export const router = createBrowserRouter([
  // Public marketing pages
  { path: '/', element: <LandingPage /> },
  { path: '/features', element: <FeaturesPage /> },
  { path: '/plans', element: <PlansPage /> },
  { path: '/about', element: <AboutPage /> },

  // Auth
  { path: '/login', element: <LoginPage /> },
  { path: '/auth/callback', element: <AuthCallback /> },
  { path: '/setup-encryption', element: <SetupEncryptionPage /> },

  // Authenticated app
  {
    path: '/app',
    element: <AppLayout />,
    children: [
      { index: true, element: <Navigate to="/app/notes" replace /> },
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
