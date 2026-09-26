import { useEffect, useState } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { useAppStore } from './stores/appStore'
import { Layout } from './components/Layout'
import { SetupWizard } from './pages/SetupWizard'
import { Dashboard } from './pages/Dashboard'
import { ContentLibrary } from './pages/ContentLibrary'
import { Agents } from './pages/Agents'
import { Settings } from './pages/Settings'

function App() {
  const { setupComplete, checkSetup, isLoading } = useAppStore()
  const [initialized, setInitialized] = useState(false)

  useEffect(() => {
    const init = async () => {
      await checkSetup()
      setInitialized(true)
    }
    init()
  }, [checkSetup])

  if (!initialized || isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-background">
        <div className="animate-pulse-glow">
          <div className="text-4xl font-bold text-primary">ABUZ8s Studio</div>
          <div className="text-muted-foreground text-center mt-2">Loading...</div>
        </div>
      </div>
    )
  }

  return (
    <Routes>
      {!setupComplete ? (
        <Route path="*" element={<SetupWizard />} />
      ) : (
        <Route element={<Layout />}>
          <Route path="/" element={<Navigate to="/dashboard" />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/content" element={<ContentLibrary />} />
          <Route path="/agents" element={<Agents />} />
          <Route path="/settings" element={<Settings />} />
        </Route>
      )}
    </Routes>
  )
}

export default App
