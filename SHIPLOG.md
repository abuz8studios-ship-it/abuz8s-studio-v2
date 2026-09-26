# SHIPLOG.md

```
✅ ABUZ8s Studio v2.0.0 — DONE-EXCEPT-BLOCKER
WHAT:     AI content creation studio for Windows (10 agents, 13 LLM providers, scheduled pipeline, SQLite library)
USE IT:   https://github.com/abuz8studios-ship-it/abuz8s-studio-v2/releases/tag/v2.0.0
PROOF:
  - clean clone + scripts/build.ps1 -> 14.3MB exe, vet clean, all Go tests pass [live-probed]
  - exe launches on fresh profile: window "ABUZ8s Studio", 10 tasks scheduled, config+db created, no panic [live-probed]
  - wizard payload applies + switches provider without restart (committed test) [build-verified]
  - no secrets in repo, history (fresh), bundle, or build output [verified]
  - no TODO/FIXME on core path [verified]
  - npm audit: 0 vulnerabilities [live-probed]
UPDATE:   binary: download new exe from Releases (settings/db kept) · source: git pull + scripts/build.ps1
BLOCKERS: live end-to-end generation proof needs Ollama running locally or a cloud API key
          (Settings -> Providers -> Test, then Run agent) — none on prova machine
NEXT:     full 13-provider Settings UI (see NEXT.md)
```

## 2026-09-26 — v2.0.0 (abuz8ify probe-to-launch)

What was broken → fixed (each with proof):

1. Fresh-startup nil panic (P0) — `initAgents` called `.Name()` on nil provider when
   nothing enabled. Now: multi-provider `llm.Client`, nil-safe startup, clean
   "no LLM provider configured" errors. Proof: `TestFreshStartupNoProvider` +
   real-exe fresh launch (`[App] Startup complete`, no panic).
2. SetupWizard discarded all input (P0) — never called `UpdateConfig`. Now persists
   niche/voice/providers, OpenRouter added. Proof: wizard-payload test in
   `app_smoke_test.go` + `tsc` clean.
3. Scheduler silently dead (P0) — `cron.WithSeconds()` vs 5-field specs, errors
   ignored. Now `cron.New()` + error logging. Proof: `scheduler_test.go`
   (failed before, passes after) + 10 "Enabled task" lines in launch log.
4. `go vet` copylocks failure (P0) — `Config.Get()` copied mutex. Now field copy.
   Proof: `go vet ./...` exit 0.
5. `TestProvider` stubbed `true` (P1) — now constructs the named provider from
   stored config and really tests it. Proof: `factory_test.go` (13/13 names).
6. Only 3/13 providers wired (P1) — now all 13 via `llm.NewProvider` factory +
   priority selection; `GetProviders`/`GetPublic` report all. Proof: factory test.
7. Settings page fully dead (P1) — all buttons wired (save niche/voice/providers,
   per-provider Test, Export JSON, Clear content); Alerts tab removed (no backend).
   Proof: `tsc && vite build` clean.
8. Dashboard stats hardcoded 0 (P1) — now live from `GetContentStats`.
9. ContentLibrary static (P1) — now live list/search/delete; Generate → Agents.
10. Build was broken ritual (P1) — exe booted to Wails "Error" dialog (missing
    `desktop,production` tags); embed copy was manual; nested `dist/dist` junk.
    Now `scripts/build.ps1` does all four steps; junk removed.
11. No git/GitHub (P1) — repo created, pushed, tagged `v2.0.0`, Release with exe.

Deliberately out of v1 (see FINISH.md + NEXT.md): mobile builds, macOS/Linux
installers, OAuth, auto-updater, signing cert, full 13-provider Settings UI,
cloud sync, Notifications tab.
