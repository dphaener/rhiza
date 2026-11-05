# web-templates-embedded-qa-1

evidence · `01M6PWYR13H5SB9VSNVN7WWQWG` · actor [claude-code-e4071629](../actor/claude-code-e4071629.md)

## result

PASS

## statement

QA at dfd42a2b: a scratch build (go build ./cmd/rhiza) ran a real long-running rhiza daemon (paused, --max-runs 0) in a fresh scratch repo with the run env unset, and curl hit it. GET / gave 200 text/html with the rendered Overview layout (header, nav, stale banner, /static/app.js). /static/app.js (text/javascript), /static/css/rhiza-tokens.css and rhiza.css (text/css), favicon.svg and .ico, a woff2 font and /site.webmanifest were all served from the binary with no files on disk beside it.

## Links

### Out

- verifies [web-templates-embedded](../capability/rhiza/web-monitor/web-templates-embedded.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
