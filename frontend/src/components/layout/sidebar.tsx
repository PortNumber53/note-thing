import { useEffect } from 'react'
import { useNavigate, useLocation } from 'react-router'
import { BookOpen, FileText, Search, Trash2, Plus, Tag, LogOut, User, Settings, CreditCard, ChevronUp } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Separator } from '@/components/ui/separator'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { useNotebooksStore } from '@/stores/notebooks-store'
import { useTagsStore } from '@/stores/tags-store'
import { useAuthStore } from '@/stores/auth-store'
import { NotebookDialog } from '@/components/notebooks/notebook-dialog'
import { cn } from '@/lib/utils'

export function Sidebar({ className }: { className?: string }) {
  const navigate = useNavigate()
  const location = useLocation()
  const { notebooks, fetchNotebooks } = useNotebooksStore()
  const { tags, fetchTags } = useTagsStore()
  const { user, logout } = useAuthStore()

  useEffect(() => {
    fetchNotebooks()
    fetchTags()
  }, [fetchNotebooks, fetchTags])

  return (
    <div className={cn("flex h-full w-60 flex-col border-r bg-sidebar", className)}>
      <div className="flex items-center gap-2 p-4">
        <BookOpen className="h-5 w-5" />
        <span className="font-semibold">Note Thing</span>
      </div>

      <ScrollArea className="flex-1 px-2">
        <nav className="space-y-1">
          <Button
            variant={location.pathname === '/notes' ? 'secondary' : 'ghost'}
            className="w-full justify-start"
            size="sm"
            onClick={() => navigate('/notes')}
          >
            <FileText className="mr-2 h-4 w-4" />
            All Notes
          </Button>
          <Button
            variant={location.pathname === '/trash' ? 'secondary' : 'ghost'}
            className="w-full justify-start"
            size="sm"
            onClick={() => navigate('/trash')}
          >
            <Trash2 className="mr-2 h-4 w-4" />
            Trash
          </Button>
          <Button
            variant={location.pathname === '/search' ? 'secondary' : 'ghost'}
            className="w-full justify-start"
            size="sm"
            onClick={() => navigate('/search')}
          >
            <Search className="mr-2 h-4 w-4" />
            Search
          </Button>
        </nav>

        <Separator className="my-3" />

        <div className="mb-1 flex items-center justify-between px-2">
          <span className="text-xs font-medium text-muted-foreground uppercase">Notebooks</span>
          <NotebookDialog>
            <Button variant="ghost" size="icon" className="h-5 w-5">
              <Plus className="h-3 w-3" />
            </Button>
          </NotebookDialog>
        </div>
        <nav className="space-y-1">
          {notebooks.map((nb) => (
            <Button
              key={nb.id}
              variant={location.pathname === `/notebooks/${nb.id}` ? 'secondary' : 'ghost'}
              className="w-full justify-start text-sm"
              size="sm"
              onClick={() => navigate(`/notebooks/${nb.id}`)}
            >
              <BookOpen className="mr-2 h-3.5 w-3.5" />
              <span className="truncate">{nb.name}</span>
              <span className="ml-auto text-xs text-muted-foreground">{nb.noteCount}</span>
            </Button>
          ))}
        </nav>

        <Separator className="my-3" />

        <div className="mb-1 px-2">
          <span className="text-xs font-medium text-muted-foreground uppercase">Tags</span>
        </div>
        <nav className="space-y-1">
          {tags.map((tag) => (
            <Button
              key={tag.id}
              variant={location.pathname === `/tags/${tag.id}` ? 'secondary' : 'ghost'}
              className="w-full justify-start text-sm"
              size="sm"
              onClick={() => navigate(`/tags/${tag.id}`)}
            >
              <Tag className="mr-2 h-3.5 w-3.5" />
              <span className="truncate">{tag.name}</span>
            </Button>
          ))}
        </nav>
      </ScrollArea>

      <Separator />
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <button className="flex w-full items-center gap-2 p-3 hover:bg-sidebar-accent transition-colors text-left">
            {user && (
              <>
                {user.avatarUrl ? (
                  <img
                    src={user.avatarUrl}
                    alt={user.name}
                    className="h-7 w-7 rounded-full"
                  />
                ) : (
                  <div className="h-7 w-7 rounded-full bg-primary text-primary-foreground flex items-center justify-center text-xs font-medium">
                    {user.name?.charAt(0)?.toUpperCase() || '?'}
                  </div>
                )}
                <span className="flex-1 truncate text-sm">{user.name}</span>
              </>
            )}
            <ChevronUp className="h-4 w-4 text-muted-foreground" />
          </button>
        </DropdownMenuTrigger>
        <DropdownMenuContent side="top" align="start" className="w-56">
          <DropdownMenuLabel>My Account</DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem onClick={() => navigate('/account/profile')}>
            <User className="mr-2 h-4 w-4" />
            Profile
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => navigate('/account/settings')}>
            <Settings className="mr-2 h-4 w-4" />
            Settings
          </DropdownMenuItem>
          <DropdownMenuItem onClick={() => navigate('/account/billing')}>
            <CreditCard className="mr-2 h-4 w-4" />
            Billing
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem onClick={logout}>
            <LogOut className="mr-2 h-4 w-4" />
            Sign out
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  )
}
