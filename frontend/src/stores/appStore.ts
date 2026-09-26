import { create } from 'zustand'
import { Config, Agent, Task, Stats, Content } from '../types'

// Window interface extension for Wails runtime
declare global {
  interface Window {
    go: {
      main: {
        App: {
          GetConfig: () => Promise<Config>
          GetFullConfig: () => Promise<Config>
          UpdateConfig: (updates: Record<string, unknown>) => Promise<void>
          GetSetupComplete: () => Promise<boolean>
          CompleteSetup: () => Promise<void>
          GetStatus: () => Promise<{ isRunning: boolean; queueSize: number; setupComplete: boolean; agentCount: number }>
          GetAgents: () => Promise<Agent[]>
          RunAgent: (agentType: string) => Promise<void>
          StartAgents: () => Promise<void>
          StopAgents: () => Promise<void>
          GetSchedulerTasks: () => Promise<Task[]>
          EnableTask: (taskId: string) => Promise<void>
          DisableTask: (taskId: string) => Promise<void>
          UpdateTaskSchedule: (taskId: string, schedule: string) => Promise<void>
          GetContentList: (type: string, status: string, limit: number, offset: number) => Promise<unknown[]>
          GetContent: (id: string) => Promise<unknown>
          DeleteContent: (id: string) => Promise<void>
          SearchContent: (query: string, limit: number) => Promise<unknown[]>
          GetContentStats: () => Promise<Stats>
          TestProvider: (provider: string) => Promise<boolean>
          GetProviders: () => Promise<Record<string, boolean>>
        }
      }
    }
  }
}

interface AppState {
  config: Config | null
  setupComplete: boolean
  isLoading: boolean
  isRunning: boolean
  agents: Agent[]
  tasks: Task[]
  stats: Stats | null
  contents: Content[]

  // Actions
  checkSetup: () => Promise<void>
  loadConfig: () => Promise<void>
  updateConfig: (updates: Record<string, unknown>) => Promise<void>
  completeSetup: () => Promise<void>
  loadAgents: () => Promise<void>
  loadTasks: () => Promise<void>
  runAgent: (agentType: string) => Promise<void>
  startAgents: () => Promise<void>
  stopAgents: () => Promise<void>
  enableTask: (taskId: string) => Promise<void>
  disableTask: (taskId: string) => Promise<void>
  updateTaskSchedule: (taskId: string, schedule: string) => Promise<void>
  loadStats: () => Promise<void>
  loadContent: (contentType: string, status: string) => Promise<void>
  searchContent: (query: string) => Promise<void>
  deleteContent: (id: string) => Promise<void>
  testProvider: (provider: string) => Promise<boolean>
  fetchAllContent: () => Promise<Content[]>
}

// Check if we're in development mode (no Wails)
const isDev = typeof window === 'undefined' || !window.go?.main?.App

// Mock data for development
const mockConfig: Config = {
  id: 'test-id',
  version: '2.0.0',
  setupComplete: false,
  niche: {
    name: '',
    description: '',
    targetAudience: '',
    keywords: [],
    competitors: [],
  },
  voice: {
    id: 'voice-1',
    name: 'Default',
    style: 'conversational',
    tone: 'professional',
    patterns: [],
    phrases: [],
    avoid: [],
    sampleText: '',
  },
  content: {
    youtube: { enabled: true, scriptsPerDay: 5, scriptLength: 1800, topics: [] },
    xPosts: { enabled: true, postsPerDay: 5, style: 'hook-first' },
    blog: { enabled: true, postsPerWeek: 7, seoEnabled: true },
    newsletter: { enabled: true, frequency: 'daily', segments: ['main'] },
    thumbnails: { enabled: true, variantsPerVideo: 3, style: 'high-contrast' },
    clips: { enabled: true, clipsPerDay: 3, platforms: ['tiktok', 'instagram', 'youtube-shorts'] },
  },
  providers: {
    priority: ['ollama', 'openrouter', 'openai'],
    defaultModel: 'gpt-4-turbo',
    openrouter: { enabled: false, apiKey: '', baseUrl: 'https://openrouter.ai/api/v1', model: '', timeout: 120 },
    opencode: { enabled: false, apiKey: '', baseUrl: 'https://api.opencode.ai/v1', model: '', timeout: 120 },
    openai: { enabled: false, apiKey: '', baseUrl: 'https://api.openai.com/v1', model: '', timeout: 120 },
    gemini: { enabled: false, apiKey: '', baseUrl: '', model: '', timeout: 120 },
    ollama: { enabled: false, apiKey: '', baseUrl: 'http://localhost:11434', model: 'llama3', timeout: 300 },
    lmstudio: { enabled: false, apiKey: '', baseUrl: 'http://localhost:1234/v1', model: '', timeout: 300 },
  },
  schedule: {
    timezone: 'America/New_York',
    times: {
      intel: '0 5 * * *',
      scripts: '0 6 * * *',
      posts: '0 6 * * *',
      thumbnails: '0 7 * * *',
      blog: '0 7 * * *',
      outreach: '0 7 * * *',
      newsletter: '0 8 * * *',
      clips: '0 */4 * * *',
    },
  },
}

const mockAgents: Agent[] = [
  { id: 'intel', type: 'intel_collector', name: 'Intel Collector', description: 'Collects trending topics', status: 'idle', prompt: '', model: '', provider: '', config: {}, createdAt: new Date().toISOString() },
  { id: 'script', type: 'script_writer', name: 'Script Writer', description: 'Writes video scripts', status: 'idle', prompt: '', model: '', provider: '', config: {}, createdAt: new Date().toISOString() },
  { id: 'xposts', type: 'x_post_generator', name: 'X Post Generator', description: 'Creates viral posts', status: 'idle', prompt: '', model: '', provider: '', config: {}, createdAt: new Date().toISOString() },
  { id: 'thumbnail', type: 'thumbnail_forge', name: 'Thumbnail Forge', description: 'Designs thumbnails', status: 'idle', prompt: '', model: '', provider: '', config: {}, createdAt: new Date().toISOString() },
  { id: 'blog', type: 'blog_writer', name: 'Blog Writer', description: 'Writes SEO articles', status: 'idle', prompt: '', model: '', provider: '', config: {}, createdAt: new Date().toISOString() },
  { id: 'outreach', type: 'outreach_engine', name: 'Outreach Engine', description: 'Composes outreach messages', status: 'idle', prompt: '', model: '', provider: '', config: {}, createdAt: new Date().toISOString() },
  { id: 'newsletter', type: 'newsletter', name: 'Newsletter', description: 'Composes newsletters', status: 'idle', prompt: '', model: '', provider: '', config: {}, createdAt: new Date().toISOString() },
  { id: 'clips', type: 'clip_factory', name: 'Clip Factory', description: 'Creates short-form clips', status: 'idle', prompt: '', model: '', provider: '', config: {}, createdAt: new Date().toISOString() },
]

const mockTasks: Task[] = [
  { id: 'intel_collector', name: 'Intel Collector', schedule: '0 5 * * *', enabled: true },
  { id: 'script_writer', name: 'Script Writer', schedule: '0 6 * * *', enabled: true },
  { id: 'x_post_generator', name: 'X Post Generator', schedule: '0 6 * * *', enabled: true },
  { id: 'thumbnail_forge', name: 'Thumbnail Forge', schedule: '0 7 * * *', enabled: true },
  { id: 'blog_writer', name: 'Blog Writer', schedule: '0 7 * * *', enabled: true },
  { id: 'outreach_engine', name: 'Outreach Engine', schedule: '0 7 * * *', enabled: true },
  { id: 'newsletter', name: 'Newsletter', schedule: '0 8 * * *', enabled: true },
  { id: 'clip_factory', name: 'Clip Factory', schedule: '0 */4 * * *', enabled: true },
]

export const useAppStore = create<AppState>((set, get) => ({
  config: null,
  setupComplete: false,
  isLoading: false,
  isRunning: false,
  agents: [],
  tasks: [],
  stats: null,
  contents: [],

  checkSetup: async () => {
    if (isDev) {
      set({ setupComplete: false })
      return
    }
    try {
      const complete = await window.go.main.App.GetSetupComplete()
      set({ setupComplete: complete })
    } catch (error) {
      console.error('Failed to check setup:', error)
    }
  },

  loadConfig: async () => {
    if (isDev) {
      set({ config: mockConfig })
      return
    }
    try {
      const config = await window.go.main.App.GetConfig()
      set({ config })
    } catch (error) {
      console.error('Failed to load config:', error)
    }
  },

  updateConfig: async (updates) => {
    if (isDev) {
      set(state => ({ config: { ...state.config!, ...updates } }))
      return
    }
    try {
      await window.go.main.App.UpdateConfig(updates)
      await get().loadConfig()
    } catch (error) {
      console.error('Failed to update config:', error)
    }
  },

  completeSetup: async () => {
    if (isDev) {
      set({ setupComplete: true })
      return
    }
    try {
      await window.go.main.App.CompleteSetup()
      set({ setupComplete: true })
    } catch (error) {
      console.error('Failed to complete setup:', error)
    }
  },

  loadAgents: async () => {
    if (isDev) {
      set({ agents: mockAgents })
      return
    }
    try {
      const agents = await window.go.main.App.GetAgents()
      set({ agents })
    } catch (error) {
      console.error('Failed to load agents:', error)
    }
  },

  loadTasks: async () => {
    if (isDev) {
      set({ tasks: mockTasks })
      return
    }
    try {
      const tasks = await window.go.main.App.GetSchedulerTasks()
      set({ tasks })
    } catch (error) {
      console.error('Failed to load tasks:', error)
    }
  },

  runAgent: async (agentType) => {
    if (isDev) {
      console.log('Running agent:', agentType)
      return
    }
    try {
      await window.go.main.App.RunAgent(agentType)
    } catch (error) {
      console.error('Failed to run agent:', error)
    }
  },

  startAgents: async () => {
    if (isDev) {
      set({ isRunning: true })
      return
    }
    try {
      await window.go.main.App.StartAgents()
      set({ isRunning: true })
    } catch (error) {
      console.error('Failed to start agents:', error)
    }
  },

  stopAgents: async () => {
    if (isDev) {
      set({ isRunning: false })
      return
    }
    try {
      await window.go.main.App.StopAgents()
      set({ isRunning: false })
    } catch (error) {
      console.error('Failed to stop agents:', error)
    }
  },

  enableTask: async (taskId) => {
    if (isDev) {
      set(state => ({
        tasks: state.tasks.map(t => t.id === taskId ? { ...t, enabled: true } : t)
      }))
      return
    }
    try {
      await window.go.main.App.EnableTask(taskId)
      await get().loadTasks()
    } catch (error) {
      console.error('Failed to enable task:', error)
    }
  },

  disableTask: async (taskId) => {
    if (isDev) {
      set(state => ({
        tasks: state.tasks.map(t => t.id === taskId ? { ...t, enabled: false } : t)
      }))
      return
    }
    try {
      await window.go.main.App.DisableTask(taskId)
      await get().loadTasks()
    } catch (error) {
      console.error('Failed to disable task:', error)
    }
  },

  updateTaskSchedule: async (taskId, schedule) => {
    if (isDev) {
      set(state => ({
        tasks: state.tasks.map(t => t.id === taskId ? { ...t, schedule } : t)
      }))
      return
    }
    try {
      await window.go.main.App.UpdateTaskSchedule(taskId, schedule)
      await get().loadTasks()
    } catch (error) {
      console.error('Failed to update task schedule:', error)
    }
  },

  loadStats: async () => {
    if (isDev) {
      set({ stats: { total: 0, types: {}, status: {} } })
      return
    }
    try {
      const stats = await window.go.main.App.GetContentStats()
      set({ stats })
    } catch (error) {
      console.error('Failed to load stats:', error)
    }
  },

  loadContent: async (contentType, status) => {
    if (isDev) {
      set({ contents: [] })
      return
    }
    try {
      const contents = await window.go.main.App.GetContentList(contentType, status, 100, 0)
      set({ contents: contents as Content[] })
    } catch (error) {
      console.error('Failed to load content:', error)
    }
  },

  searchContent: async (query) => {
    if (isDev) {
      return
    }
    try {
      if (!query.trim()) {
        await get().loadContent('', '')
        return
      }
      const contents = await window.go.main.App.SearchContent(query, 100)
      set({ contents: contents as Content[] })
    } catch (error) {
      console.error('Failed to search content:', error)
    }
  },

  deleteContent: async (id) => {
    if (isDev) {
      set(state => ({ contents: state.contents.filter(c => c.id !== id) }))
      return
    }
    try {
      await window.go.main.App.DeleteContent(id)
      set(state => ({ contents: state.contents.filter(c => c.id !== id) }))
    } catch (error) {
      console.error('Failed to delete content:', error)
    }
  },

  testProvider: async (provider) => {
    if (isDev) {
      console.log('Testing provider:', provider)
      return false
    }
    try {
      return await window.go.main.App.TestProvider(provider)
    } catch (error) {
      console.error('Failed to test provider:', error)
      return false
    }
  },

  fetchAllContent: async () => {
    if (isDev) {
      return []
    }
    try {
      return (await window.go.main.App.GetContentList('', '', 10000, 0)) as Content[]
    } catch (error) {
      console.error('Failed to fetch content:', error)
      return []
    }
  },
}))
