import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date
  return d.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

export function formatDateTime(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date
  return d.toLocaleString('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function formatDuration(ms: number): string {
  const seconds = Math.floor(ms / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  
  if (hours > 0) {
    return `${hours}h ${minutes % 60}m`
  }
  if (minutes > 0) {
    return `${minutes}m ${seconds % 60}s`
  }
  return `${seconds}s`
}

export function truncate(str: string, length: number): string {
  if (str.length <= length) return str
  return str.slice(0, length) + '...'
}

export function debounce<T extends (...args: unknown[]) => unknown>(
  fn: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: ReturnType<typeof setTimeout>
  return (...args) => {
    clearTimeout(timeoutId)
    timeoutId = setTimeout(() => fn(...args), delay)
  }
}

export function getAgentIcon(type: string): string {
  const icons: Record<string, string> = {
    intel_collector: 'Search',
    script_writer: 'FileText',
    x_post_generator: 'Twitter',
    thumbnail_forge: 'Image',
    blog_writer: 'BookOpen',
    outreach_engine: 'Send',
    newsletter: 'Mail',
    clip_factory: 'Scissors',
    performance_eval: 'BarChart3',
    weekly_digest: 'Calendar',
  }
  return icons[type] || 'Bot'
}

export function getContentTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    script: 'Video Script',
    thumbnail: 'Thumbnail',
    x_post: 'X Post',
    blog_post: 'Blog Post',
    newsletter: 'Newsletter',
    outreach: 'Outreach',
    clip_script: 'Clip Script',
    trend_report: 'Trend Report',
    idea: 'Idea',
  }
  return labels[type] || type
}

export function getStatusColor(status: string): string {
  const colors: Record<string, string> = {
    draft: 'bg-gray-500',
    scheduled: 'bg-blue-500',
    published: 'bg-green-500',
    failed: 'bg-red-500',
    processing: 'bg-yellow-500',
    archived: 'bg-gray-400',
  }
  return colors[status] || 'bg-gray-500'
}
