export interface Config {
  id: string
  version: string
  setupComplete: boolean
  niche: NicheConfig
  voice: VoiceConfig
  content: ContentConfig
  providers: ProvidersConfig
  schedule: ScheduleConfig
}

export interface NicheConfig {
  name: string
  description: string
  targetAudience: string
  keywords: string[]
  competitors: string[]
}

export interface VoiceConfig {
  id: string
  name: string
  style: string
  tone: string
  patterns: string[]
  phrases: string[]
  avoid: string[]
  sampleText: string
}

export interface ContentConfig {
  youtube: YouTubeConfig
  xPosts: XPostsConfig
  blog: BlogConfig
  newsletter: NewsletterConfig
  thumbnails: ThumbnailsConfig
  clips: ClipsConfig
}

export interface YouTubeConfig {
  enabled: boolean
  scriptsPerDay: number
  scriptLength: number
  topics: string[]
}

export interface XPostsConfig {
  enabled: boolean
  postsPerDay: number
  style: string
}

export interface BlogConfig {
  enabled: boolean
  postsPerWeek: number
  seoEnabled: boolean
}

export interface NewsletterConfig {
  enabled: boolean
  frequency: string
  segments: string[]
}

export interface ThumbnailsConfig {
  enabled: boolean
  variantsPerVideo: number
  style: string
}

export interface ClipsConfig {
  enabled: boolean
  clipsPerDay: number
  platforms: string[]
}

export interface ProvidersConfig {
  priority: string[]
  defaultModel: string
  openrouter: ProviderConfig
  opencode: ProviderConfig
  openai: ProviderConfig
  gemini: ProviderConfig
  ollama: ProviderConfig
  lmstudio: ProviderConfig
}

export interface ProviderConfig {
  enabled: boolean
  apiKey: string
  baseUrl: string
  model: string
  timeout: number
}

export interface ScheduleConfig {
  timezone: string
  times: {
    intel: string
    scripts: string
    posts: string
    thumbnails: string
    blog: string
    outreach: string
    newsletter: string
    clips: string
  }
}

export interface Content {
  id: string
  type: ContentType
  title: string
  body: string
  status: ContentStatus
  provider: string
  model: string
  metadata: Record<string, unknown>
  tags: string[]
  scheduledAt?: string
  publishedAt?: string
  createdAt: string
  updatedAt: string
}

export type ContentType = 
  | 'script' 
  | 'thumbnail' 
  | 'x_post' 
  | 'blog_post' 
  | 'newsletter' 
  | 'outreach' 
  | 'clip_script' 
  | 'trend_report' 
  | 'idea'

export type ContentStatus = 
  | 'draft' 
  | 'scheduled' 
  | 'published' 
  | 'failed' 
  | 'processing' 
  | 'archived'

export interface Agent {
  id: string
  type: string
  name: string
  description: string
  status: AgentStatus
  prompt: string
  model: string
  provider: string
  config: Record<string, unknown>
  lastRun?: string
  createdAt: string
}

export type AgentStatus = 'idle' | 'running' | 'disabled' | 'error'

export interface Task {
  id: string
  name: string
  schedule: string
  enabled: boolean
  nextRun?: string
}

export interface Stats {
  total: number
  types: Record<string, number>
  status: Record<string, number>
}
