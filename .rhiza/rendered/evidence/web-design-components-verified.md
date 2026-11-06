# web-design-components-verified

evidence · `01M6Q37C01S5V2CDM48G31WN6V` · actor [claude-code-f480fe06](../actor/claude-code-f480fe06.md)

## result

PASS

## statement

At af1e9e67, go test -trimpath -count=1 ./internal/web/... passes; TestLayoutGolden renders the layout (overview with daemon pill, runs without) and TestPartialsGolden renders stale, pill (all 5 run and 5 change set states), panel with kicker, glyph (5 kinds, both sizes), rows with teal rz-id short IDs, steps timeline with sessions, always-dark term block and empty state with fixture data, comparing each against internal/web/testdata/*.golden.html. Re-running with -update left git status clean, so the goldens are what -update regenerates. Read the diff: header carries the daemon pill, live indicator and nav to Overview, Runs and Roadmap; partial markup matches the delivered design pages (pages/run.html, change-set.html) and uses classes defined in static/css/rhiza.css; go vet and gofmt clean.

## Links

### Out

- verifies [web-design-components](../capability/rhiza/web-monitor/web-design-components.md)

## History

- `01M6PFFFZNADHH2WGMB5MES8QT` Web design system: the Claude Design system built from the owner's basics as one embedded stylesheet, embedded fonts and the shared layout and partials every page renders with
