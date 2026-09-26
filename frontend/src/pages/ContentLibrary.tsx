import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAppStore } from '../stores/appStore'
import { Card, CardContent } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Search, FileText, Image, Twitter, BookOpen, Mail, Send, Trash2 } from 'lucide-react'
import { getContentTypeLabel } from '../lib/utils'

const contentTypes = [
  { id: 'all', name: 'All Content', icon: null },
  { id: 'script', name: 'Video Scripts', icon: FileText },
  { id: 'x_post', name: 'X Posts', icon: Twitter },
  { id: 'thumbnail', name: 'Thumbnails', icon: Image },
  { id: 'blog_post', name: 'Blog Posts', icon: BookOpen },
  { id: 'newsletter', name: 'Newsletters', icon: Mail },
  { id: 'outreach', name: 'Outreach', icon: Send },
]

export function ContentLibrary() {
  const navigate = useNavigate()
  const { contents, loadContent, searchContent, deleteContent } = useAppStore()
  const [selectedType, setSelectedType] = useState('all')
  const [searchQuery, setSearchQuery] = useState('')

  useEffect(() => {
    loadContent('', '')
  }, [loadContent])

  const handleTypeChange = (type: string) => {
    setSelectedType(type)
    setSearchQuery('')
    loadContent(type === 'all' ? '' : type, '')
  }

  const handleSearch = (query: string) => {
    setSearchQuery(query)
    if (!query.trim()) {
      loadContent(selectedType === 'all' ? '' : selectedType, '')
    } else {
      searchContent(query)
    }
  }

  const handleDelete = async (id: string) => {
    if (window.confirm('Delete this content item?')) {
      await deleteContent(id)
    }
  }

  return (
    <div className="space-y-6 animate-slide-in">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold">Content Library</h1>
        <p className="text-muted-foreground mt-1">
          Browse and manage all your generated content
        </p>
      </div>

      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder="Search content..."
            className="pl-10"
            value={searchQuery}
            onChange={(e) => handleSearch(e.target.value)}
          />
        </div>
      </div>

      {/* Content Type Tabs */}
      <div className="flex flex-wrap gap-2">
        {contentTypes.map((type) => (
          <Button
            key={type.id}
            variant={selectedType === type.id ? 'default' : 'outline'}
            size="sm"
            onClick={() => handleTypeChange(type.id)}
            className="gap-2"
          >
            {type.icon && <type.icon className="h-4 w-4" />}
            {type.name}
          </Button>
        ))}
      </div>

      {/* Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-4">
        {contents.length === 0 ? (
          <Card className="col-span-full py-12">
            <CardContent className="text-center">
              <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-muted">
                <FileText className="h-6 w-6 text-muted-foreground" />
              </div>
              <h3 className="font-medium">No content yet</h3>
              <p className="text-sm text-muted-foreground mt-1">
                Start the content pipeline to generate your first pieces
              </p>
              <Button className="mt-4" onClick={() => navigate('/agents')}>
                Generate Content
              </Button>
            </CardContent>
          </Card>
        ) : (
          contents.map((item) => (
            <Card key={item.id} className="hover:border-primary/50 transition-colors">
              <CardContent className="pt-6">
                <div className="flex items-start justify-between gap-2">
                  <div>
                    <span className="text-xs px-2 py-1 rounded-full bg-muted">
                      {getContentTypeLabel(item.type)}
                    </span>
                    <h3 className="font-medium mt-2 line-clamp-2">{item.title || 'Untitled'}</h3>
                  </div>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => handleDelete(item.id)}
                    title="Delete"
                  >
                    <Trash2 className="h-4 w-4 text-destructive" />
                  </Button>
                </div>
                <p className="text-sm text-muted-foreground mt-2 line-clamp-4">{item.body}</p>
                <div className="flex items-center justify-between mt-4 text-xs text-muted-foreground">
                  <span>{item.provider}{item.model ? ` · ${item.model}` : ''}</span>
                  <span>{item.status}</span>
                </div>
              </CardContent>
            </Card>
          ))
        )}
      </div>
    </div>
  )
}
