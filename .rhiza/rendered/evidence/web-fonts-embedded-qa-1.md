# web-fonts-embedded-qa-1

evidence · `01M6Q3CJTF8W8RJE1KGKHVF1NA` · actor [claude-code-37ef1809](../actor/claude-code-37ef1809.md)

## result

PASS

## statement

At de99be4e built rhiza with go build -trimpath, ran 'rhiza daemon --max-runs 0' in a scratch repo with rhiza.web.port=47913 and fetched the live web UI with curl. /static/fonts/{Petrona,InstrumentSans,JetBrainsMono}-Variable.woff2 each served 200 font/woff2 with wOF2 magic (86400, 89136, 71628 bytes); each OFL-*.txt served with 'SIL OPEN FONT LICENSE'; the served / page references only /-rooted src/href paths (css, app.js, icons, manifest), all of which returned 200.

## Links

### Out

- verifies [web-fonts-embedded](../capability/rhiza/web-monitor/web-fonts-embedded.md)

## History

- `01M6PFFFZNADHH2WGMB5MES8QT` Web design system: the Claude Design system built from the owner's basics as one embedded stylesheet, embedded fonts and the shared layout and partials every page renders with
