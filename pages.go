package main

import (
	"encoding/json"
	"fmt"

	"github.com/rohanthewiz/element"
)

const darkCSS = `
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

body {
	background-color: #0d1117;
	color: #e6edf3;
	font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif;
	font-size: 15px;
	line-height: 1.6;
	height: 100vh;
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

a { color: #58a6ff; text-decoration: none; }
a:hover { color: #79c0ff; text-decoration: underline; }

/* ── Top bar ── */
header {
	background: #161b22;
	border-bottom: 1px solid #30363d;
	padding: 0.7rem 1.25rem;
	display: flex;
	align-items: center;
	flex-shrink: 0;
}

.logo {
	color: #f0f6fc;
	font-size: 1rem;
	font-weight: 600;
	text-decoration: none;
}
.logo:hover { color: #58a6ff; text-decoration: none; }

.project-path {
	margin-left: 1rem;
	color: #8b949e;
	font-size: 0.8rem;
	font-family: 'SFMono-Regular', Consolas, monospace;
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}

/* ── Two-panel layout ── */
.layout {
	display: flex;
	flex: 1;
	overflow: hidden;
}

/* ── Sidebar ── */
.sidebar {
	width: 270px;
	min-width: 200px;
	background: #161b22;
	border-right: 1px solid #30363d;
	overflow-y: auto;
	flex-shrink: 0;
	display: flex;
	flex-direction: column;
}

.sidebar-header {
	padding: 0.65rem 1rem;
	border-bottom: 1px solid #30363d;
	color: #8b949e;
	font-size: 0.7rem;
	font-weight: 600;
	text-transform: uppercase;
	letter-spacing: 0.08em;
	flex-shrink: 0;
}

.sidebar-item {
	display: block;
	padding: 0.65rem 1rem;
	border-bottom: 1px solid #21262d;
	text-decoration: none;
	border-left: 3px solid transparent;
}
.sidebar-item:hover {
	background: #1c2128;
	text-decoration: none;
}
.sidebar-item.active {
	background: #1c2128;
	border-left-color: #58a6ff;
}

.sidebar-num {
	color: #8b949e;
	font-size: 0.7rem;
	margin-bottom: 0.1rem;
}
.sidebar-id {
	color: #e6edf3;
	font-family: 'SFMono-Regular', Consolas, monospace;
	font-size: 0.78rem;
	word-break: break-all;
	margin-bottom: 0.15rem;
}
.sidebar-item.active .sidebar-id { color: #79c0ff; }
.sidebar-ts {
	color: #8b949e;
	font-size: 0.72rem;
}

/* ── Main content ── */
.main {
	flex: 1;
	overflow-y: auto;
	padding: 2rem 2.25rem;
}

h1 { color: #f0f6fc; font-size: 1.5rem; margin-bottom: 0.4rem; font-weight: 600; }

.welcome {
	max-width: 520px;
	padding-top: 3rem;
}
.welcome p { color: #8b949e; margin-top: 0.5rem; }

/* ── Response cards ── */
.response {
	background: #161b22;
	border: 1px solid #30363d;
	border-radius: 6px;
	margin-bottom: 1.25rem;
	overflow: hidden;
	max-width: 900px;
}

.response-header {
	background: #1c2128;
	border-bottom: 1px solid #30363d;
	padding: 0.5rem 1rem;
	font-size: 0.8rem;
	color: #8b949e;
	display: flex;
	align-items: center;
	gap: 0.75rem;
}
.resp-num { font-weight: 600; color: #58a6ff; }

.response-body {
	padding: 1.25rem 1.5rem;
	color: #e6edf3;
	overflow-x: auto;
}

/* ── Markdown rendered content ── */
.response-body h1, .response-body h2, .response-body h3,
.response-body h4, .response-body h5, .response-body h6 {
	color: #f0f6fc;
	margin: 1.2rem 0 0.6rem;
	font-weight: 600;
	line-height: 1.3;
}
.response-body h1 { font-size: 1.4rem; border-bottom: 1px solid #21262d; padding-bottom: 0.4rem; }
.response-body h2 { font-size: 1.2rem; border-bottom: 1px solid #21262d; padding-bottom: 0.3rem; }
.response-body h3 { font-size: 1.05rem; }
.response-body p  { margin: 0.6rem 0; }

.response-body code {
	background: #1f2937;
	color: #79c0ff;
	padding: 0.15em 0.4em;
	border-radius: 4px;
	font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', monospace;
	font-size: 0.85em;
}
.response-body pre {
	background: #1f2937;
	border: 1px solid #30363d;
	border-radius: 6px;
	padding: 1rem;
	overflow-x: auto;
	margin: 0.75rem 0;
}
.response-body pre code { background: transparent; padding: 0; color: #e6edf3; font-size: 0.875em; }

.response-body ul, .response-body ol { margin: 0.6rem 0 0.6rem 1.5rem; }
.response-body li { margin: 0.2rem 0; }
.response-body blockquote { border-left: 3px solid #30363d; padding-left: 1rem; color: #8b949e; margin: 0.75rem 0; }
.response-body table { border-collapse: collapse; width: 100%; margin: 0.75rem 0; font-size: 0.9rem; }
.response-body table th, .response-body table td { border: 1px solid #30363d; padding: 0.4rem 0.8rem; }
.response-body table th { background: #1c2128; color: #8b949e; font-weight: 600; }
.response-body table tr:hover td { background: #1c2128; }
.response-body hr { border: none; border-top: 1px solid #30363d; margin: 1rem 0; }
.response-body a { color: #58a6ff; }
.response-body a:hover { color: #79c0ff; }
.response-body strong { color: #f0f6fc; }
.response-body em { color: #c9d1d9; }

.no-responses { color: #8b949e; padding: 2rem 0; }

.error-box {
	background: #1c1317;
	border: 1px solid #6e2030;
	border-radius: 6px;
	padding: 1.5rem;
	color: #f85149;
	max-width: 600px;
}
.error-box h1 { color: #f85149; margin-bottom: 0.5rem; }
`

const markedCDN = "https://cdn.jsdelivr.net/npm/marked/marked.min.js"

func renderPage(b *element.Builder, sessions []sessionInfo, activeID, displayDir, extraHead string, mainFn func()) {
	b.Html().R(
		pageHead(b, "Claude Code History", extraHead),
		b.Body().R(
			// Top bar
			b.Header().R(
				b.A("href", "/", "class", "logo").T("Claude Code History"),
				b.Wrap(func() {
					if displayDir != "" {
						b.Span("class", "project-path").T(displayDir)
					}
				}),
			),
			// Two-panel layout
			b.DivClass("layout").R(
				// Sidebar
				b.Nav("class", "sidebar").R(
					b.DivClass("sidebar-header").T("Sessions"),
					b.Wrap(func() {
						for i, s := range sessions {
							cls := "sidebar-item"
							if s.id == activeID {
								cls += " active"
							}
							b.A("href", "/session/"+s.id, "class", cls).R(
								b.DivClass("sidebar-num").F("#%d", i+1),
								b.DivClass("sidebar-id").T(s.id),
								b.DivClass("sidebar-ts").T(s.timestamp),
							)
						}
					}),
				),
				// Main area
				b.Main("class", "main").R(
					b.Wrap(mainFn),
				),
			),
		),
	)
}

func pageHead(b *element.Builder, title string, extraHead string) (x any) {
	b.Head().R(
		b.Meta("charset", "utf-8").R(),
		b.Meta("name", "viewport", "content", "width=device-width, initial-scale=1").R(),
		b.Title().T(title),
		b.Style().R(b.Wrap(func() { b.WriteString(darkCSS) })),
		b.Wrap(func() { b.WriteString(extraHead) }),
	)
	return
}

func renderSessionsPage(sessions []sessionInfo, displayDir string) string {
	b := element.AcquireBuilder()
	defer element.ReleaseBuilder(b)

	renderPage(b, sessions, "", displayDir, "", func() {
		b.DivClass("welcome").R(
			b.H1().T("Claude Code History"),
			b.P().T("Select a session from the sidebar to view its responses."),
		)
	})
	return b.String()
}

func renderSessionPage(id, timestamp string, responses []Response, sessions []sessionInfo) string {
	b := element.AcquireBuilder()
	defer element.ReleaseBuilder(b)

	type jsResp struct {
		ID    string   `json:"id"`
		Num   int      `json:"num"`
		Ts    string   `json:"ts"`
		Texts []string `json:"texts"`
	}

	jsData := make([]jsResp, len(responses))
	for i, r := range responses {
		jsData[i] = jsResp{
			ID:    fmt.Sprintf("resp-%d", r.Number),
			Num:   r.Number,
			Ts:    r.Timestamp,
			Texts: r.Texts,
		}
	}
	jsonBytes, _ := json.Marshal(jsData)

	renderScript := fmt.Sprintf(`
const _responses = %s;
document.addEventListener('DOMContentLoaded', function() {
	_responses.forEach(function(r) {
		var el = document.getElementById(r.id);
		if (el) {
			el.innerHTML = r.texts.map(function(t) { return marked.parse(t); }).join('');
		}
	});
});
`, string(jsonBytes))

	extraHead := fmt.Sprintf(`<script src="%s"></script>`, markedCDN)

	renderPage(b, sessions, id, "", extraHead, func() {
		b.H1().T("Session")
		b.Wrap(func() {
			if len(responses) == 0 {
				b.PClass("no-responses").T("No assistant responses found in this session.")
				return
			}
			for _, r := range responses {
				b.DivClass("response").R(
					b.DivClass("response-header").R(
						b.Span("class", "resp-num").F("Response %d", r.Number),
						b.Span().T(r.Timestamp),
					),
					b.Div("class", "response-body", "id", fmt.Sprintf("resp-%d", r.Number)).R(),
				)
			}
		})
		b.Script("type", "text/javascript").R(
			b.Wrap(func() { b.WriteString(renderScript) }),
		)
	})
	return b.String()
}

func renderErrorPage(msg string) string {
	b := element.AcquireBuilder()
	defer element.ReleaseBuilder(b)

	b.Html().R(
		pageHead(b, "Error — Claude Code History", ""),
		b.Body().R(
			b.Header().R(
				b.A("href", "/", "class", "logo").T("Claude Code History"),
			),
			b.DivClass("layout").R(
				b.Main("class", "main").R(
					b.DivClass("error-box").R(
						b.H1().T("Error"),
						b.P().T(msg),
					),
				),
			),
		),
	)
	return b.String()
}
