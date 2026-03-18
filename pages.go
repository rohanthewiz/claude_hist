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
}

header {
	background: #161b22;
	border-bottom: 1px solid #30363d;
	padding: 0.9rem 0;
}

header .inner {
	max-width: 980px;
	margin: 0 auto;
	padding: 0 1.5rem;
	display: flex;
	align-items: center;
	gap: 1rem;
}

.logo {
	color: #f0f6fc;
	font-size: 1rem;
	font-weight: 600;
	text-decoration: none;
}

.logo:hover { color: #58a6ff; text-decoration: none; }

.container {
	max-width: 980px;
	margin: 0 auto;
	padding: 2rem 1.5rem;
}

h1 { color: #f0f6fc; font-size: 1.6rem; margin-bottom: 0.4rem; font-weight: 600; }
h2 { color: #8b949e; font-size: 0.95rem; font-weight: 400; margin-bottom: 1.75rem; font-family: 'SFMono-Regular', Consolas, monospace; word-break: break-all; }

a { color: #58a6ff; text-decoration: none; }
a:hover { color: #79c0ff; text-decoration: underline; }

.sessions-table {
	width: 100%;
	border-collapse: collapse;
	border: 1px solid #30363d;
	border-radius: 6px;
	overflow: hidden;
}

.sessions-table th {
	background: #161b22;
	color: #8b949e;
	font-weight: 600;
	font-size: 0.75rem;
	text-transform: uppercase;
	letter-spacing: 0.08em;
	padding: 0.6rem 1rem;
	text-align: left;
	border-bottom: 1px solid #30363d;
}

.sessions-table td {
	padding: 0.75rem 1rem;
	border-bottom: 1px solid #21262d;
	color: #e6edf3;
	vertical-align: middle;
}

.sessions-table tr:last-child td { border-bottom: none; }
.sessions-table tbody tr:hover td { background: #161b22; }

.col-num { color: #8b949e; font-size: 0.875rem; width: 3.5rem; }
.col-ts  { color: #8b949e; font-size: 0.875rem; white-space: nowrap; }
.session-link { font-family: 'SFMono-Regular', Consolas, monospace; font-size: 0.875rem; }

.response {
	background: #161b22;
	border: 1px solid #30363d;
	border-radius: 6px;
	margin-bottom: 1.25rem;
	overflow: hidden;
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

.response-body pre code {
	background: transparent;
	padding: 0;
	color: #e6edf3;
	font-size: 0.875em;
}

.response-body ul, .response-body ol { margin: 0.6rem 0 0.6rem 1.5rem; }
.response-body li { margin: 0.2rem 0; }

.response-body blockquote {
	border-left: 3px solid #30363d;
	padding-left: 1rem;
	color: #8b949e;
	margin: 0.75rem 0;
}

.response-body table { border-collapse: collapse; width: 100%; margin: 0.75rem 0; font-size: 0.9rem; }
.response-body table th, .response-body table td { border: 1px solid #30363d; padding: 0.4rem 0.8rem; }
.response-body table th { background: #1c2128; color: #8b949e; font-weight: 600; }
.response-body table tr:hover td { background: #1c2128; }

.response-body hr { border: none; border-top: 1px solid #30363d; margin: 1rem 0; }
.response-body a { color: #58a6ff; }
.response-body a:hover { color: #79c0ff; }
.response-body strong { color: #f0f6fc; }
.response-body em { color: #c9d1d9; }

.no-sessions, .no-responses {
	color: #8b949e;
	text-align: center;
	padding: 3rem 0;
	font-size: 1rem;
}

.error-box {
	background: #1c1317;
	border: 1px solid #6e2030;
	border-radius: 6px;
	padding: 1.5rem;
	color: #f85149;
}

.error-box h1 { color: #f85149; margin-bottom: 0.5rem; }
`

const markedCDN = "https://cdn.jsdelivr.net/npm/marked/marked.min.js"

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

func pageHeader(b *element.Builder, logoText, logoHref string) (x any) {
	b.Header().R(
		b.DivClass("inner").R(
			b.A("href", logoHref, "class", "logo").T(logoText),
		),
	)
	return
}

func renderSessionsPage(sessions []sessionInfo, displayDir string) string {
	b := element.AcquireBuilder()
	defer element.ReleaseBuilder(b)

	b.Html().R(
		pageHead(b, "Claude Code History", ""),
		b.Body().R(
			pageHeader(b, "Claude Code History", "/"),
			b.DivClass("container").R(
				b.H1().T("Sessions"),
				b.H2().T(displayDir),
				b.Wrap(func() {
					if len(sessions) == 0 {
						b.PClass("no-sessions").T("No sessions found for this project.")
						return
					}
					b.Table("class", "sessions-table").R(
						b.THead().R(
							b.Tr().R(
								b.Th("class", "col-num").T("#"),
								b.Th().T("Session ID"),
								b.Th("class", "col-ts").T("Started"),
							),
						),
						b.TBody().R(
							b.Wrap(func() {
								for i, s := range sessions {
									b.Tr().R(
										b.Td("class", "col-num").F("%d", i+1),
										b.Td().R(
											b.A("href", "/session/"+s.id, "class", "session-link").T(s.id),
										),
										b.Td("class", "col-ts").T(s.timestamp),
									)
								}
							}),
						),
					)
				}),
			),
		),
	)
	return b.String()
}

func renderSessionPage(id, timestamp string, responses []Response) string {
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

	b.Html().R(
		pageHead(b, "Session: "+id, fmt.Sprintf(`<script src="%s"></script>`, markedCDN)),
		b.Body().R(
			pageHeader(b, "← Sessions", "/"),
			b.DivClass("container").R(
				b.H1().T("Session"),
				b.H2().T(id),
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
				}),
			),
			b.Script("type", "text/javascript").R(
				b.Wrap(func() { b.WriteString(renderScript) }),
			),
		),
	)
	return b.String()
}

func renderErrorPage(msg string) string {
	b := element.AcquireBuilder()
	defer element.ReleaseBuilder(b)

	b.Html().R(
		pageHead(b, "Error — Claude Code History", ""),
		b.Body().R(
			pageHeader(b, "← Sessions", "/"),
			b.DivClass("container").R(
				b.DivClass("error-box").R(
					b.H1().T("Error"),
					b.P().T(msg),
				),
			),
		),
	)
	return b.String()
}