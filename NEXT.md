# NEXT.md — v1.1 backlog (one line each)

- Wire full 13-provider list into Settings providers tab (backend factory already supports all).
- Wire OAuth sign-in (internal/auth exists, unused) or delete it.
- Capacitor mobile: sync + EAS/Studio builds for Android/iOS (dirs present, unbuilt).
- macOS + Linux installers via GitHub Actions matrix (Wails per-OS package).
- Code-sign Windows exe (Azure Trusted Signing) to drop the SmartScreen note.
- Auto-updater (wails-update or Tauri-style) + update channel.
- Notifications backend (Settings Alerts tab was removed — no backend existed).
- Cloud sync for config/content (StorageConfig.CloudSync is a stub flag).
- Streaming generation UI (provider Stream exists; no frontend surface).
- E2E generation proof with live Ollama or cloud key (blocked: neither on prova machine).
- Export scheduler tasks to visible next-run times in Agents page (GetNextRun exists).
