# web-design-tokens-qa-1

evidence · `01M6Q3CJS06QHP2WC8GC6RPVK1` · actor [claude-code-37ef1809](../actor/claude-code-37ef1809.md)

## result

PASS

## statement

At de99be4e built rhiza with go build -trimpath, ran 'rhiza daemon --max-runs 0' in a scratch repo with rhiza.web.port=47913 and fetched the live web UI with curl. GET /static/css/rhiza-tokens.css served 200 text/css: dark soil tokens under :root/[data-theme=dark] (gold #e8d49a, teal #86c9b7, rose #d9a48f), light under @@media (prefers-color-scheme: light), radius-md 12px, radius-lg 16px, @@font-face for Petrona, Instrument Sans, JetBrains Mono; neither served stylesheet contains a gradient or an http(s) URL.

## Links

### Out

- verifies [web-design-tokens](../capability/rhiza/web-monitor/web-design-tokens.md)

## History

- `01M6PFFFZNADHH2WGMB5MES8QT` Web design system: the Claude Design system built from the owner's basics as one embedded stylesheet, embedded fonts and the shared layout and partials every page renders with
