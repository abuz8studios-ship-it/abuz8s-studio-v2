# FINISH.md — ABUZ8s Studio v2 (abuz8ify run 2026-09-26)

```
FINISH LINE:   Windows exe on GitHub Release v2.0.0; launch → setup wizard saves
               provider → start pipeline → generated content appears in library
TARGET:        Windows x64 exe (unsigned, "Run anyway" note in README) + GitHub Release
CORE FLOW:     1. Launch exe → setup wizard
               2. Configure niche/voice, enable provider (Ollama local or cloud key) → saved
               3. Start pipeline or Run agent → agent executes via provider
               4. Content appears in Content Library; Dashboard stats update
               5. Schedules persist; app restarts clean with no panic and no provider
MONEY PATH:    n/a (free local-first tool, no payments in v1)
OUT OF SCOPE:  mobile builds (Capacitor dirs stay, unbuilt) · macOS/Linux installers ·
               OAuth sign-in (code present, unwired) · auto-updater · code-signing cert ·
               full 13-provider Settings UI (Ollama/LM Studio/OpenAI/OpenRouter wired) ·
               cloud sync · Notifications settings tab (no backend — removed)
```

## Detect (Phase 0)

- What: desktop app — AI content studio (10 agents, 13 LLM providers, cron pipeline, SQLite)
- Does: generates scheduled content (scripts, X posts, blogs, newsletters, clips…) via local or cloud LLMs
- Who: solo creators; operator-run (ABUZ8)
- Where: desktop exe, Windows x64 first; local-first, offline-capable with Ollama
- Data: stateful — `%APPDATA%/ABUZ8sStudio/{config.json,data/abuz8s.db}`
- Money: n/a v1
- Route: `/ship-desktop` (Wails v2, existing stack — no migration) + `/ship-github`

## Gate notes

- Live LLM generation proof needs Ollama or a cloud API key on the proving machine
  (neither present here). Pipeline-through-fake-provider is proven by committed Go
  tests; exe launch + wizard-save + scheduling are proven on the real artifact.
- Deploy gate items requiring a clean VM, signing cert, or updater infra are
  recorded as BLOCKERs with exact user actions, not faked.
