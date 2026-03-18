# Session: add-web-server

**Date:** 2026-03-17
**Goal:** Convert `claude_hist` CLI tool into a web app while keeping all original CLI functionality.

---

## What Was Done

### Libraries Used

Picked up versions from `/Users/RAllison3/cbre_projs/edp_dataflow/go.mod`:

| Library | Version | Purpose |
|---------|---------|---------|
| `github.com/rohanthewiz/rweb` | `v0.1.25-0.20260314122458-7e864a80be8e` | HTTP web server |
| `github.com/rohanthewiz/element` | `v0.5.6` | HTML builder |
| `github.com/rohanthewiz/logger` | `v1.3.0` | Structured logging |
| `github.com/rohanthewiz/serr` | `v1.3.0` | Structured error wrapping |

SKILL.md locations:
- `/Users/RAllison3/projs/go/rweb/ai_docs/SKILL.md`
- `/Users/RAllison3/projs/go/element/ai_docs/SKILL.md`
- `/Users/RAllison3/projs/go/logger/ai_docs/SKILL.md`
- serr has no SKILL.md — read `/Users/RAllison3/projs/go/exported/serr/serr.go` directly

### Files Created / Modified

#### `history.go` (new)
Extracted all data types and extraction logic from the original `main.go`:
- Types: `contentBlock`, `message`, `record`, `sessionInfo`, `Response`
- Functions: `resolveProjectDir`, `getSessions`, `firstTimestamp`, `getResponses`, `formatTimestamp`
- `getResponses` returns `[]Response` (data) instead of printing — new function for web use
- All errors wrapped with `serr.Wrap` / `serr.New`

#### `pages.go` (new)
HTML page builders using the `element` library:
- Dark GitHub-inspired theme via inline CSS (`darkCSS` const)
- `pageHead(b, title, extraHead)` — shared `<head>` helper, returns `(x any)` for use in `.R()`
- `pageHeader(b, logoText, logoHref)` — shared `<header>` helper, same return convention
- `renderSessionsPage(sessions, displayDir)` — table of sessions with clickable links
- `renderSessionPage(id, timestamp, responses)` — response cards with markdown rendering
- `renderErrorPage(msg)` — error display
- Markdown rendering: response texts are JSON-encoded into the page; `marked.js` (CDN) renders them client-side via `document.addEventListener('DOMContentLoaded', ...)`
- Used `b.Wrap(func() { b.WriteString(rawContent) })` inside `b.Style()` and `b.Script()` to inject raw CSS/JS without HTML encoding

#### `main.go` (rewritten)
Kept all original CLI behavior, added `--web [addr]` mode:

```
Usage: claude_hist [OPTIONS] [SESSION_ID | NUMBER]

Options:
  (none)              Show responses from the latest session
  <1-999>             Show responses for session at that index (see --list)
  --dir <path>        Use the given project path instead of the current directory
  --list              List available sessions with numeric indices
  --all               Show responses from all sessions
  --web [addr]        Start web UI (default addr: :7070)
  --help, -h          Show this help
```

`serveWeb(projectDir, displayDir, addr)` sets up rweb server with:
- `GET /` → sessions list page
- `GET /session/:id` → session detail with responses

#### `go.mod` (updated)
Added the four library dependencies; `go mod tidy` added transitive deps (`logrus`, `golang.org/x/sys`). Go version bumped to 1.24.0 by tidy.

### Key Patterns Used

**element helper functions** must return `(x any)` to be usable inside `.R()`:
```go
func pageHead(b *element.Builder, title string, extraHead string) (x any) {
    b.Head().R(...)
    return
}
```

**Raw content in `<style>`/`<script>` tags** — use `b.Wrap` + `b.WriteString`:
```go
b.Style().R(b.Wrap(func() { b.WriteString(cssContent) }))
b.Script("type", "text/javascript").R(b.Wrap(func() { b.WriteString(jsCode) }))
```

**Inline iteration** (when index needed) — use `b.Wrap` with a regular Go loop:
```go
b.TBody().R(
    b.Wrap(func() {
        for i, s := range sessions {
            b.Tr().R(...)
        }
    }),
)
```

**JSON embedding for JS data** — `json.Marshal` escapes `<`, `>`, `&` by default (safe for HTML):
```go
jsonBytes, _ := json.Marshal(jsData)
renderScript := fmt.Sprintf(`const _responses = %s; ...`, string(jsonBytes))
```

### Usage

```bash
# Web mode (current project dir)
cchist --web

# Web mode with custom dir and port
cchist --dir /path/to/project --web :8080

# CLI still works as before
cchist --list
cchist 2
cchist > responses.md
```
