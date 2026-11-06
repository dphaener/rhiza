# web-templates-embedded-verified

evidence · `01M6PTSS0V1CFZQZHMR445DRFP` · actor [claude-code-58d0f11f](../actor/claude-code-58d0f11f.md)

## result

PASS

## statement

Verifier claude-code-58d0f11f at 3615b0ae: read the diff (internal/web/web.go parses templates/layout.html + templates/overview.html via go:embed html/template; GET /{$} renders the overview page; static/js/live.js renamed to static/app.js and served from the embedded FS under /static/). Ran go test -trimpath -count=1 -v -run TestHandlerRendersPageAndScript ./internal/web/: PASS — GET / returns 200 text/html with the rendered title, stylesheet and script tags; GET /static/app.js returns 200 with a javascript content type and a non-empty body; /nope is 404. Whole ./internal/web/... and the Web/Events tests in ./internal/api pass; go vet and gofmt clean. The commit leaves Makefile and .github/workflows untouched, and grep finds no node/npm/yarn/esbuild/webpack/vite in either; no stale js/live.js references remain.

## Links

### Out

- verifies [web-templates-embedded](../capability/rhiza/web-monitor/web-templates-embedded.md)

## History

- `01M6PFEXSHWC6TR4F3GFKHV09R` Web server: the daemon serves a read-only web UI on 127.0.0.1:7787 (rhiza.web.port, or off), with read-only internal/api operations as JSON, server-sent events and embedded templates and assets
