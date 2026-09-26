import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAppStore } from '../stores/appStore'
import { Button } from '../components/ui/Button'
import { Input } from '../components/ui/Input'
import { Textarea } from '../components/ui/Textarea'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../components/ui/Card'
import { Zap, ChevronRight, ChevronLeft, Check } from 'lucide-react'
import { cn } from '../lib/utils'

const steps = [
  { id: 'niche', title: 'Niche & Brand', description: 'Define your content niche' },
  { id: 'voice', title: 'Voice Profile', description: 'Set your content voice' },
  { id: 'providers', title: 'AI Providers', description: 'Configure LLM providers' },
  { id: 'review', title: 'Review', description: 'Confirm your settings' },
]

export function SetupWizard() {
  const navigate = useNavigate()
  const { completeSetup, updateConfig, loadConfig } = useAppStore()
  const [currentStep, setCurrentStep] = useState(0)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    loadConfig()
  }, [loadConfig])
  const [formData, setFormData] = useState({
    niche: {
      name: '',
      description: '',
      targetAudience: '',
    },
    voice: {
      style: 'conversational',
      tone: 'professional',
    },
    providers: {
      openai: { enabled: false, apiKey: '' },
      openrouter: { enabled: false, apiKey: '' },
      ollama: { enabled: false, baseUrl: 'http://localhost:11434' },
    },
  })

  const handleNext = async () => {
    if (currentStep < steps.length - 1) {
      setCurrentStep(currentStep + 1)
      return
    }
    // Final step: persist everything the user entered, then mark setup done.
    setSaving(true)
    try {
      await updateConfig({
        niche: { ...formData.niche },
        voice: { ...formData.voice },
        providers: {
          openai: { ...formData.providers.openai },
          openrouter: { ...formData.providers.openrouter },
          ollama: { ...formData.providers.ollama },
        },
      })
      await completeSetup()
      navigate('/dashboard')
    } finally {
      setSaving(false)
    }
  }

  const handleBack = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1)
    }
  }

  const renderStep = () => {
    switch (currentStep) {
      case 0:
        return (
          <div className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Niche Name</label>
              <Input
                placeholder="e.g., AI Technology, Personal Finance, Fitness"
                value={formData.niche.name}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    niche: { ...formData.niche, name: e.target.value },
                  })
                }
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Description</label>
              <Textarea
                placeholder="Describe your content niche in detail..."
                value={formData.niche.description}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    niche: { ...formData.niche, description: e.target.value },
                  })
                }
              />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium">Target Audience</label>
              <Input
                placeholder="e.g., Tech enthusiasts, Small business owners"
                value={formData.niche.targetAudience}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    niche: { ...formData.niche, targetAudience: e.target.value },
                  })
                }
              />
            </div>
          </div>
        )
      case 1:
        return (
          <div className="space-y-4">
            <div className="space-y-2">
              <label className="text-sm font-medium">Voice Style</label>
              <select
                className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                value={formData.voice.style}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    voice: { ...formData.voice, style: e.target.value },
                  })
                }
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
                value={formData.voice.tone}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    voice: { ...formData.voice, tone: e.target.value },
                  })
                }
              >
                <option value="professional">Professional</option>
                <option value="friendly">Friendly</option>
                <option value="inspirational">Inspirational</option>
                <option value="educational">Educational</option>
                <option value="controversial">Controversial</option>
              </select>
            </div>
          </div>
        )
      case 2:
        return (
          <div className="space-y-6">
            <div className="space-y-4">
              <h3 className="font-medium">Local Models (Recommended)</h3>
              <div className="flex items-center justify-between p-4 rounded-lg border">
                <div>
                  <p className="font-medium">Ollama</p>
                  <p className="text-sm text-muted-foreground">
                    Run models locally on your machine
                  </p>
                </div>
                <input
                  type="checkbox"
                  checked={formData.providers.ollama.enabled}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      providers: {
                        ...formData.providers,
                        ollama: {
                          ...formData.providers.ollama,
                          enabled: e.target.checked,
                        },
                      },
                    })
                  }
                  className="h-5 w-5"
                />
              </div>
            </div>
            <div className="space-y-4">
              <h3 className="font-medium">Cloud Providers</h3>
              <div className="space-y-3">
                <div className="flex items-center justify-between p-4 rounded-lg border">
                  <div>
                    <p className="font-medium">OpenAI</p>
                    <p className="text-sm text-muted-foreground">
                      GPT-4, GPT-3.5 Turbo
                    </p>
                  </div>
                  <input
                    type="checkbox"
                    checked={formData.providers.openai.enabled}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        providers: {
                          ...formData.providers,
                          openai: {
                            ...formData.providers.openai,
                            enabled: e.target.checked,
                          },
                        },
                      })
                    }
                    className="h-5 w-5"
                  />
                </div>
                {formData.providers.openai.enabled && (
                  <Input
                    type="password"
                    placeholder="OpenAI API Key"
                    value={formData.providers.openai.apiKey}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        providers: {
                          ...formData.providers,
                          openai: {
                            ...formData.providers.openai,
                            apiKey: e.target.value,
                          },
                        },
                      })
                    }
                  />
                )}
                <div className="flex items-center justify-between p-4 rounded-lg border">
                  <div>
                    <p className="font-medium">OpenRouter</p>
                    <p className="text-sm text-muted-foreground">
                      100+ models, one API key
                    </p>
                  </div>
                  <input
                    type="checkbox"
                    checked={formData.providers.openrouter.enabled}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        providers: {
                          ...formData.providers,
                          openrouter: {
                            ...formData.providers.openrouter,
                            enabled: e.target.checked,
                          },
                        },
                      })
                    }
                    className="h-5 w-5"
                  />
                </div>
                {formData.providers.openrouter.enabled && (
                  <Input
                    type="password"
                    placeholder="OpenRouter API Key"
                    value={formData.providers.openrouter.apiKey}
                    onChange={(e) =>
                      setFormData({
                        ...formData,
                        providers: {
                          ...formData.providers,
                          openrouter: {
                            ...formData.providers.openrouter,
                            apiKey: e.target.value,
                          },
                        },
                      })
                    }
                  />
                )}
              </div>
            </div>
          </div>
        )
      case 3:
        return (
          <div className="space-y-6">
            <div className="p-4 rounded-lg bg-muted">
              <h3 className="font-medium mb-2">Niche Configuration</h3>
              <dl className="space-y-1 text-sm">
                <div className="flex justify-between">
                  <dt className="text-muted-foreground">Name:</dt>
                  <dd>{formData.niche.name || 'Not set'}</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-muted-foreground">Audience:</dt>
                  <dd>{formData.niche.targetAudience || 'Not set'}</dd>
                </div>
              </dl>
            </div>
            <div className="p-4 rounded-lg bg-muted">
              <h3 className="font-medium mb-2">Voice Profile</h3>
              <dl className="space-y-1 text-sm">
                <div className="flex justify-between">
                  <dt className="text-muted-foreground">Style:</dt>
                  <dd>{formData.voice.style}</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-muted-foreground">Tone:</dt>
                  <dd>{formData.voice.tone}</dd>
                </div>
              </dl>
            </div>
            <div className="p-4 rounded-lg bg-muted">
              <h3 className="font-medium mb-2">Providers</h3>
              <dl className="space-y-1 text-sm">
                <div className="flex justify-between">
                  <dt className="text-muted-foreground">Ollama:</dt>
                  <dd>{formData.providers.ollama.enabled ? 'Enabled' : 'Disabled'}</dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-muted-foreground">OpenAI:</dt>
                  <dd>
                    {formData.providers.openai.enabled
                      ? 'Enabled (API key set)'
                      : 'Disabled'}
                  </dd>
                </div>
                <div className="flex justify-between">
                  <dt className="text-muted-foreground">OpenRouter:</dt>
                  <dd>
                    {formData.providers.openrouter.enabled
                      ? 'Enabled (API key set)'
                      : 'Disabled'}
                  </dd>
                </div>
              </dl>
            </div>
          </div>
        )
      default:
        return null
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <Card className="w-full max-w-2xl">
        <CardHeader className="text-center">
          <div className="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-lg bg-primary">
            <Zap className="h-7 w-7 text-primary-foreground" />
          </div>
          <CardTitle className="text-2xl">Welcome to ABUZ8s Studio</CardTitle>
          <CardDescription>
            Let's get your content creation pipeline set up
          </CardDescription>
        </CardHeader>
        <CardContent>
          {/* Progress */}
          <div className="mb-8">
            <div className="flex items-center justify-between">
              {steps.map((step, index) => (
                <div key={step.id} className="flex items-center">
                  <div
                    className={cn(
                      'flex h-8 w-8 items-center justify-center rounded-full text-sm font-medium',
                      index <= currentStep
                        ? 'bg-primary text-primary-foreground'
                        : 'bg-muted text-muted-foreground'
                    )}
                  >
                    {index < currentStep ? (
                      <Check className="h-4 w-4" />
                    ) : (
                      index + 1
                    )}
                  </div>
                  {index < steps.length - 1 && (
                    <div
                      className={cn(
                        'h-1 w-12',
                        index < currentStep ? 'bg-primary' : 'bg-muted'
                      )}
                    />
                  )}
                </div>
              ))}
            </div>
            <div className="mt-4">
              <h3 className="font-medium">{steps[currentStep].title}</h3>
              <p className="text-sm text-muted-foreground">
                {steps[currentStep].description}
              </p>
            </div>
          </div>

          {/* Step Content */}
          <div className="min-h-[300px]">{renderStep()}</div>

          {/* Navigation */}
          <div className="flex justify-between mt-8">
            <Button
              variant="outline"
              onClick={handleBack}
              disabled={currentStep === 0}
            >
              <ChevronLeft className="mr-2 h-4 w-4" />
              Back
            </Button>
            <Button onClick={handleNext} disabled={saving}>
              {currentStep === steps.length - 1 ? (
                saving ? 'Saving...' : 'Complete Setup'
              ) : (
                <>
                  Next
                  <ChevronRight className="ml-2 h-4 w-4" />
                </>
              )}
            </Button>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
