// Wails runtime bindings
// This file provides type-safe access to the Wails backend

export interface WailsConfig {
  GetConfig: () => Promise<unknown>
  GetFullConfig: () => Promise<unknown>
  UpdateConfig: (updates: unknown) => Promise<void>
  GetSetupComplete: () => Promise<boolean>
  CompleteSetup: () => Promise<void>
  GetStatus: () => Promise<{
    isRunning: boolean
    queueSize: number
    setupComplete: boolean
    agentCount: number
  }>
  GetAgents: () => Promise<unknown[]>
  RunAgent: (agentType: string) => Promise<void>
  StartAgents: () => Promise<void>
  StopAgents: () => Promise<void>
  GetSchedulerTasks: () => Promise<unknown[]>
  EnableTask: (taskId: string) => Promise<void>
  DisableTask: (taskId: string) => Promise<void>
  UpdateTaskSchedule: (taskId: string, schedule: string) => Promise<void>
  GetContentList: (type: string, status: string, limit: number, offset: number) => Promise<unknown[]>
  GetContent: (id: string) => Promise<unknown>
  DeleteContent: (id: string) => Promise<void>
  SearchContent: (query: string, limit: number) => Promise<unknown[]>
  GetContentStats: () => Promise<unknown>
  TestProvider: (provider: string) => Promise<boolean>
  GetProviders: () => Promise<Record<string, boolean>>
}

export function getWailsRuntime(): WailsConfig | null {
  if (typeof window === 'undefined') return null
  
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const wails = (window as any).go?.main?.App
  if (!wails) return null
  
  return wails as WailsConfig
}

export function isWailsEnvironment(): boolean {
  return getWailsRuntime() !== null
}
