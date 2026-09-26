import { useEffect, useState } from 'react'
import { useAppStore } from '../stores/appStore'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/Card'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Textarea } from '../components/ui/Textarea'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '../components/ui/Tabs'
import { Key, User, Database } from 'lucide-react'

type TestState = 'testing' | 'ok' | 'fail' | undefined

export function Settings() {
  const { config, loadConfig, updateConfig, testProvider, fetchAllContent, deleteContent } = useAppStore()
  const [activeTab, setActiveTab] = useState('general')
  const [saving, setSaving] = useState(false)
  const [savedTick, setSavedTick] = useState(0)

  const [niche, setNiche] = useState({ name: '', description: '', targetAudience: '' })
  const [voice, setVoice] = useState({ style: 'conversational', tone: 'professional' })
  const [ollama, setOllama] = useState({ enabled: false, baseUrl: 'http://localhost:11434', model: 'llama3' })
  const [lmstudio, setLmstudio] = useState({ enabled: false, baseUrl: 'http://localhost:1234/v1' })
  const [openai, setOpenai] = useState({ enabled: false, apiKey: '' })
  const [openrouter, setOpenrouter] = useState({ enabled: false, apiKey: '' })
  const [tests, setTests] = useState<Record<string, TestState>>({})

  useEffect(() => {
    loadConfig()
  }, [loadConfig])

  useEffect(() => {
    if (!config) return
    setNiche({
      name: config.niche?.name ?? '',
      description: config.niche?.description ?? '',
      targetAudience: config.niche?.targetAudience ?? '',
    })
    setVoice({
      style: config.voice?.style ?? 'conversational',
      tone: config.voice?.tone ?? 'professional',
    })
    const p = config.providers
    if (p) {
      setOllama({ enabled: p.ollama?.enabled ?? false, baseUrl: p.ollama?.baseUrl ?? 'http://localhost:11434', model: p.ollama?.model ?? 'llama3' })
      setLmstudio({ enabled: p.lmstudio?.enabled ?? false, baseUrl: p.lmstudio?.baseUrl ?? 'http://localhost:1234/v1' })
      setOpenai({ enabled: p.openai?.enabled ?? false, apiKey: p.openai?.apiKey ?? '' })
      setOpenrouter({ enabled: p.openrouter?.enabled ?? false, apiKey: p.openrouter?.apiKey ?? '' })
    }
  }, [config])

  const saveGeneral = async () => {
    setSaving(true)
    try {
      await updateConfig({ niche: { ...niche }, voice: { ...voice } })
      setSavedTick((n) => n + 1)
    } finally {
      setSaving(false)
    }
  }

  const saveProviders = async () => {
    setSaving(true)
    try {
      await updateConfig({
        providers: {
          ollama: { ...ollama },
          lmstudio: { ...lmstudio },
          openai: { ...openai },
          openrouter: { ...openrouter },
        },
      })
      setSavedTick((n) => n + 1)
    } finally {
      setSaving(false)
    }
  }

  const runTest = async (name: string) => {
    setTests((t) => ({ ...t, [name]: 'testing' }))
    const ok = await testProvider(name)
    setTests((t) => ({ ...t, [name]: ok ? 'ok' : 'fail' }))
  }

  const testLabel = (name: string) => {
    const s = tests[name]
    if (s === 'testing') return 'Testing...'
    if (s === 'ok') return 'Connected'
    if (s === 'fail') return 'Failed'
    return 'Test'
  }

  const exportData = async () => {
    const all = await fetchAllContent()
    const blob = new Blob([JSON.stringify(all, null, 2)], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `abuz8s-studio-export-${new Date().toISOString().slice(0, 10)}.json`
    a.click()
    URL.revokeObjectURL(url)
  }

  const clearData = async () => {
    if (!window.confirm('Delete ALL generated content? This cannot be undone.')) return
    const all = await fetchAllContent()
    for (const item of all) {
      await deleteContent(item.id)
    }
  }

  return (
    <div className="space-y-6 animate-slide-in">
      {/* Header */}
      <div>
        <h1 className="text-3xl font-bold">Settings</h1>
        <p className="text-muted-foreground mt-1">
          Configure your studio preferences{savedTick > 0 ? ' · Saved' : ''}
        </p>
      </div>

      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList className="grid w-full grid-cols-3 lg:w-[300px]">
          <TabsTrigger value="general" className="gap-2">
            <User className="h-4 w-4" />
            General
          </TabsTrigger>
          <TabsTrigger value="providers" className="gap-2">
            <Key className="h-4 w-4" />
            Providers
          </TabsTrigger>
          <TabsTrigger value="data" className="gap-2">
            <Database className="h-4 w-4" />
            Data
          </TabsTrigger>
        </TabsList>

        <TabsContent value="general" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Niche & Brand</CardTitle>
              <CardDescription>
                Update your content niche settings
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Niche Name</label>
                <Input value={niche.name} onChange={(e) => setNiche({ ...niche, name: e.target.value })} />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Description</label>
                <Textarea value={niche.description} onChange={(e) => setNiche({ ...niche, description: e.target.value })} />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Target Audience</label>
                <Input value={niche.targetAudience} onChange={(e) => setNiche({ ...niche, targetAudience: e.target.value })} />
              </div>
              <Button onClick={saveGeneral} disabled={saving}>{saving ? 'Saving...' : 'Save Changes'}</Button>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Voice Profile</CardTitle>
              <CardDescription>
                Customize your content voice
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Voice Style</label>
                <select
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={voice.style}
                  onChange={(e) => setVoice({ ...voice, style: e.target.value })}
                >
                  <option value="conversational">Conversational</option>
                  <option value="professional">Professional</option>
                  <option value="casual">Casual</option>
                  <option value="authoritative">Authoritative</option>
                  <option value="humorous">Humorous</option>
                </select>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Tone</label>
                <select
                  className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  value={voice.tone}
                  onChange={(e) => setVoice({ ...voice, tone: e.target.value })}
                >
                  <option value="professional">Professional</option>
                  <option value="friendly">Friendly</option>
                  <option value="inspirational">Inspirational</option>
                  <option value="educational">Educational</option>
                  <option value="controversial">Controversial</option>
                </select>
              </div>
              <Button onClick={saveGeneral} disabled={saving}>{saving ? 'Saving...' : 'Save Changes'}</Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="providers" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Local Models</CardTitle>
              <CardDescription>
                Configure local LLM providers
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between p-4 rounded-lg border">
                <div>
                  <p className="font-medium">Ollama</p>
                  <p className="text-sm text-muted-foreground">{ollama.baseUrl}</p>
                </div>
                <div className="flex items-center gap-2">
                  <Button variant="outline" size="sm" onClick={() => runTest('ollama')}>{testLabel('ollama')}</Button>
                  <input
                    type="checkbox"
                    className="h-5 w-5"
                    checked={ollama.enabled}
                    onChange={(e) => setOllama({ ...ollama, enabled: e.target.checked })}
                  />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Ollama Model</label>
                <Input value={ollama.model} onChange={(e) => setOllama({ ...ollama, model: e.target.value })} />
              </div>
              <div className="flex items-center justify-between p-4 rounded-lg border">
                <div>
                  <p className="font-medium">LM Studio</p>
                  <p className="text-sm text-muted-foreground">{lmstudio.baseUrl}</p>
                </div>
                <div className="flex items-center gap-2">
                  <Button variant="outline" size="sm" onClick={() => runTest('lmstudio')}>{testLabel('lmstudio')}</Button>
                  <input
                    type="checkbox"
                    className="h-5 w-5"
                    checked={lmstudio.enabled}
                    onChange={(e) => setLmstudio({ ...lmstudio, enabled: e.target.checked })}
                  />
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Cloud Providers</CardTitle>
              <CardDescription>
                Configure cloud LLM providers
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between p-4 rounded-lg border">
                <div>
                  <p className="font-medium">OpenAI</p>
                  <p className="text-sm text-muted-foreground">GPT-4, GPT-3.5 Turbo</p>
                </div>
                <div className="flex items-center gap-2">
                  <Button variant="outline" size="sm" onClick={() => runTest('openai')}>{testLabel('openai')}</Button>
                  <input
                    type="checkbox"
                    className="h-5 w-5"
                    checked={openai.enabled}
                    onChange={(e) => setOpenai({ ...openai, enabled: e.target.checked })}
                  />
                </div>
              </div>
              {openai.enabled && (
                <div className="space-y-2">
                  <label className="text-sm font-medium">OpenAI API Key</label>
                  <Input
                    type="password"
                    placeholder="sk-..."
                    value={openai.apiKey}
                    onChange={(e) => setOpenai({ ...openai, apiKey: e.target.value })}
                  />
                </div>
              )}
              <div className="flex items-center justify-between p-4 rounded-lg border">
                <div>
                  <p className="font-medium">OpenRouter</p>
                  <p className="text-sm text-muted-foreground">100+ models, one key</p>
                </div>
                <div className="flex items-center gap-2">
                  <Button variant="outline" size="sm" onClick={() => runTest('openrouter')}>{testLabel('openrouter')}</Button>
                  <input
                    type="checkbox"
                    className="h-5 w-5"
                    checked={openrouter.enabled}
                    onChange={(e) => setOpenrouter({ ...openrouter, enabled: e.target.checked })}
                  />
                </div>
              </div>
              {openrouter.enabled && (
                <div className="space-y-2">
                  <label className="text-sm font-medium">OpenRouter API Key</label>
                  <Input
                    type="password"
                    placeholder="sk-..."
                    value={openrouter.apiKey}
                    onChange={(e) => setOpenrouter({ ...openrouter, apiKey: e.target.value })}
                  />
                </div>
              )}
              <Button onClick={saveProviders} disabled={saving}>{saving ? 'Saving...' : 'Save Providers'}</Button>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="data" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Data Management</CardTitle>
              <CardDescription>
                Manage your stored data
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex items-center justify-between p-4 rounded-lg border border-destructive/50">
                <div>
                  <p className="font-medium text-destructive">Clear All Data</p>
                  <p className="text-sm text-muted-foreground">
                    Delete all generated content (settings are kept)
                  </p>
                </div>
                <Button variant="destructive" onClick={clearData}>Clear Data</Button>
              </div>
              <div className="flex items-center justify-between p-4 rounded-lg border">
                <div>
                  <p className="font-medium">Export Data</p>
                  <p className="text-sm text-muted-foreground">
                    Download all your content as JSON
                  </p>
                </div>
                <Button variant="outline" onClick={exportData}>Export</Button>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
