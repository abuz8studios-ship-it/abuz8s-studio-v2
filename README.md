# ABUZ8s Studio v2

AI-Powered Content Creation Studio - Native Desktop & Mobile Application

## Features

- **Universal LLM Support**: 10+ providers (OpenRouter, OpenAI, Ollama, Gemini, Groq, etc.)
- **10 Content Agents**: Intel Collector, Script Writer, X Posts, Thumbnails, Blog, Newsletter, Outreach, Clips, Analytics
- **Cross-Platform**: Windows, macOS, Linux (desktop) + iOS, Android (mobile)
- **Local-First**: Run models locally with Ollama/LM Studio or cloud APIs
- **OAuth Authentication**: Google, GitHub, Discord sign-in
- **Scheduled Content**: Cron-based content pipeline

## Tech Stack

- **Backend**: Go 1.21 + Wails v2
- **Frontend**: React 18 + TypeScript + Vite + Tailwind CSS
- **Mobile**: Capacitor 8
- **Database**: SQLite

## Quick Start

### Run Compiled Binary
```powershell
cd backend
.\ABUZ8sStudio.exe
```

> Unsigned build: on first launch Windows SmartScreen may show
> "Windows protected your PC" → click **More info → Run anyway**.
> (Signing cert is tracked in NEXT.md.)

### Development Mode
```powershell
# Terminal 1: Frontend
cd frontend
npm run dev

# Terminal 2: Backend (with live reload)
cd backend
wails dev
```

## Mobile Builds

### Android
```powershell
cd frontend
npm run build
npx cap sync android
npx cap open android
# Build APK in Android Studio
```

### iOS
```powershell
cd frontend
npm run build
npx cap sync ios
npx cap open ios
# Build in Xcode
```

## LLM Providers

| Provider | Type | Status |
|----------|------|--------|
| OpenRouter | Cloud | ✅ |
| OpenAI | Cloud | ✅ |
| Gemini | Cloud | ✅ |
| Groq | Cloud | ✅ |
| Together AI | Cloud | ✅ |
| Nous Research | Cloud | ✅ |
| Kimi (Moonshot) | Cloud | ✅ |
| MiniMax | Cloud | ✅ |
| GLM (Zhipu) | Cloud | ✅ |
| Ollama | Local | ✅ |
| LM Studio | Local | ✅ |

## Content Agents

| Agent | Schedule | Description |
|-------|----------|-------------|
| Intel Collector | 5AM daily | Trending topics & competitor analysis |
| Script Writer | 6AM daily | YouTube video scripts |
| X Post Generator | 6AM daily | Viral X/Twitter posts |
| Thumbnail Forge | 7AM daily | Thumbnail concepts |
| Blog Writer | 7AM daily | SEO blog posts |
| Outreach Engine | 7AM daily | Collaboration messages |
| Newsletter | 8AM daily | Email newsletters |
| Clip Factory | Every 4 hours | Short-form content scripts |
| Performance Eval | Weekly | Analytics & optimization |
| Weekly Digest | Monday | Weekly summary |

## Project Structure

```
abuz8s-studio-v2/
├── backend/
│   ├── app.go              # Wails app bindings
│   ├── main.go             # Entry point
│   ├── internal/
│   │   ├── agents/         # Content agents
│   │   ├── auth/           # OAuth implementation
│   │   ├── config/         # Settings management
│   │   ├── models/         # Database models
│   │   └── scheduler/      # Cron scheduler
│   └── pkg/
│       └── llm/            # LLM provider adapters
├── frontend/
│   ├── src/
│   │   ├── pages/          # Setup, Dashboard, etc.
│   │   ├── components/     # UI components
│   │   └── stores/         # Zustand state
│   ├── android/            # Android project
│   └── ios/                # iOS project
└── README.md
```

## Configuration

Config stored in:
- Windows: `%APPDATA%/ABUZ8sStudio/config.json`
- macOS: `~/Library/Application Support/ABUZ8sStudio/config.json`
- Linux: `~/.config/abuz8s-studio/config.json`

## Building from Source

### Prerequisites
- Go 1.21+
- Node.js 20+
- Wails CLI: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

### Build Commands
```powershell
# One command: frontend → embed → vet + tests → exe (backend\ABUZ8sStudio.exe)
powershell -ExecutionPolicy Bypass -File scripts\build.ps1
```

Manual steps (what the script does):
```powershell
# Frontend
cd frontend
npm install
npm run build

# Sync embedded assets (go:embed reads backend\frontend\dist)
Remove-Item ..\backend\frontend\dist -Recurse -Force -ErrorAction SilentlyContinue
Copy-Item dist ..\backend\frontend\dist -Recurse -Force

# Backend (Wails needs the desktop,production tags on go build)
cd ..\backend
go vet ./... ; go test ./...
go build -tags desktop,production -trimpath -ldflags "-s -w" -o ABUZ8sStudio.exe .
```

## License

MIT
