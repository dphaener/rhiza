# web-design-components-qa-1

evidence · `01M6Q3CJVQRE9NE73P2HFEEKG0` · actor [claude-code-37ef1809](../actor/claude-code-37ef1809.md)

## result

PASS

## statement

At de99be4e built rhiza with go build -trimpath, ran 'rhiza daemon --max-runs 0' in a scratch repo with rhiza.web.port=47913 and fetched the live web UI with curl. GET / served the layout: rz-header with nav Overview (aria-current=page), Runs, Roadmap, the rz-header-end live indicator and the stale partial with its warning glyph; the daemon pill is omitted on the live page because the overview handler passes no Daemon pill (the golden test covers it). go test -trimpath -count=1 ./internal/web/... passes; every rz-* class in layout.html and partials.html is defined in static/css/rhiza.css.

## Links

### Out

- verifies [web-design-components](../capability/rhiza/web-monitor/web-design-components.md)

## History

- `01M6PFFFZNADHH2WGMB5MES8QT` Web design system: the Claude Design system built from the owner's basics as one embedded stylesheet, embedded fonts and the shared layout and partials every page renders with
