import { useEffect } from 'react'
import { useAppStore } from '../stores/appStore'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { 
  Search, 
  FileText, 
  Twitter, 
  Image, 
  BookOpen, 
  Mail, 
  Send,
  BarChart3,
  Calendar,
  Scissors,
  Play,
  Settings
} from 'lucide-react'
import { cn } from '../lib/utils'

const agentIcons: Record<string, React.ElementType> = {
  intel_collector: Search,
  script_writer: FileText,
  x_post_generator: Twitter,
  thumbnail_forge: Image,
  blog_writer: BookOpen,
  outreach_engine: Send,
  newsletter: Mail,
  clip_factory: Scissors,
  performance_eval: BarChart3,
  weekly_digest: Calendar,
}

export function Agents() {
  const { agents, loadAgents, runAgent } = useAppStore()

  useEffect(() => {
    loadAgents()
  }, [loadAgents])

  return (
    <div className="space-y-6 animate-slide-in">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold">Agents</h1>
        <p className="text-muted-foreground mt-1">
          Manage your content creation agents
        </p>
      </div>

      {/* Agents Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
        {agents.map((agent) => {
          const Icon = agentIcons[agent.type] || FileText
          
          return (
            <Card key={agent.id} className="group hover:border-primary/50 transition-colors">
              <CardHeader className="pb-3">
                <div className="flex items-start justify-between">
                  <div className="flex items-center gap-3">
                    <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10">
                      <Icon className="h-5 w-5 text-primary" />
                    </div>
                    <div>
                      <CardTitle className="text-base">{agent.name}</CardTitle>
                      <CardDescription className="text-xs">
                        {agent.description}
                      </CardDescription>
                    </div>
                  </div>
                  <div
                    className={cn(
                      'h-2 w-2 rounded-full',
                      agent.status === 'running'
                        ? 'bg-green-500 animate-pulse'
                        : agent.status === 'error'
                        ? 'bg-red-500'
                        : 'bg-gray-500'
                    )}
                  />
                </div>
              </CardHeader>
              <CardContent>
                <div className="flex items-center justify-between">
                  <div className="text-xs text-muted-foreground">
                    {agent.lastRun ? (
                      <>Last run: {new Date(agent.lastRun).toLocaleDateString()}</>
                    ) : (
                      'Never run'
                    )}
                  </div>
                  <div className="flex gap-1">
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8"
                      onClick={() => runAgent(agent.type)}
                    >
                      <Play className="h-4 w-4" />
                    </Button>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8"
                    >
                      <Settings className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          )
        })}
      </div>
    </div>
  )
}
