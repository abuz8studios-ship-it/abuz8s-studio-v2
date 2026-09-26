import { useEffect } from 'react'
import { useAppStore } from '../stores/appStore'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { 
  FileText, 
  Twitter, 
  Image, 
  BookOpen, 
  Mail, 
  Send,
  Play,
  TrendingUp
} from 'lucide-react'

const statDefs = [
  { name: 'Video Scripts', type: 'script', icon: FileText, color: 'text-blue-500' },
  { name: 'X Posts', type: 'x_post', icon: Twitter, color: 'text-sky-500' },
  { name: 'Thumbnails', type: 'thumbnail', icon: Image, color: 'text-purple-500' },
  { name: 'Blog Posts', type: 'blog_post', icon: BookOpen, color: 'text-green-500' },
  { name: 'Newsletters', type: 'newsletter', icon: Mail, color: 'text-yellow-500' },
  { name: 'Outreach', type: 'outreach', icon: Send, color: 'text-pink-500' },
]

export function Dashboard() {
  const { agents, loadAgents, isRunning, startAgents, stopAgents, stats, loadStats } = useAppStore()

  useEffect(() => {
    loadAgents()
    loadStats()
  }, [loadAgents, loadStats])

  // const runningAgents = agents.filter(a => a.status === 'running')

  return (
    <div className="space-y-6 animate-slide-in">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <p className="text-muted-foreground mt-1">
            Overview of your content creation pipeline
          </p>
        </div>
        <Button
          onClick={isRunning ? stopAgents : startAgents}
          className="gap-2"
        >
          <Play className="h-4 w-4" />
          {isRunning ? 'Stop Pipeline' : 'Start Pipeline'}
        </Button>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {statDefs.map((stat) => (
          <Card key={stat.name} className="hover:border-primary/50 transition-colors">
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">{stat.name}</CardTitle>
              <stat.icon className={`h-4 w-4 ${stat.color}`} />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{stats?.types?.[stat.type] ?? 0}</div>
              <p className="text-xs text-muted-foreground">
                Items in library
              </p>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Recent Activity */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <TrendingUp className="h-5 w-5" />
            Agent Status
          </CardTitle>
          <CardDescription>
            Current status of your content agents
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {agents.slice(0, 4).map((agent) => (
              <div
                key={agent.id}
                className="flex items-center justify-between p-3 rounded-lg bg-muted/50"
              >
                <div className="flex items-center gap-3">
                  <div
                    className={`h-2 w-2 rounded-full ${
                      agent.status === 'running'
                        ? 'bg-green-500 animate-pulse'
                        : agent.status === 'error'
                        ? 'bg-red-500'
                        : 'bg-gray-500'
                    }`}
                  />
                  <div>
                    <p className="font-medium">{agent.name}</p>
                    <p className="text-sm text-muted-foreground">
                      {agent.description}
                    </p>
                  </div>
                </div>
                <span className="text-xs px-2 py-1 rounded-full bg-muted">
                  {agent.status}
                </span>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
