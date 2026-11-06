# web-design-components-verified-2

evidence · `01M6Q3MG1D78WNAC10X4550DDZ` · actor [claude-code-f1ae5389](../actor/claude-code-f1ae5389.md)

## result

PASS

## statement

Verifier, attempt 2, at 1ef23973 in a detached checkout: go build ./... and go vet ./internal/web pass; go test -trimpath -count=1 -v ./internal/web/ passes all 16 tests, including TestLayoutGolden (layout-overview, layout-runs: header with daemon pill and live indicator, nav to Overview, Runs and Roadmap) and TestPartialsGolden (panel with kicker, pills per state, glyphs, rows, steps, term, empty, stale). Running the test with -update rewrote the golden files with no diff (git status was clean), so the goldens match the -update output. Read the fix diff: Row and Session lose their ID fields; the row partial makes the title the link (the whole row follows it through ::after) and drops the ID column; the steps partial links each session by who ran it and the model; ID strings remain only in href attributes of rows.golden and steps.golden. This answers the quality and model findings under owner-surfaces-no-ids. The CSS adds titled grid columns for desktop and mobile. internal/web/pages/*.html still show IDs, but they are the delivered design reference: they are not embedded or served (go:embed covers only templates, static and the manifest). Note: the term fixture's path text reads 'session 01M6PG2C'. That is caller-supplied fixture data, not something the partial adds.

## Links

### Out

- verifies [web-design-components](../capability/rhiza/web-monitor/web-design-components.md)

## History

- `01M6PFFFZNADHH2WGMB5MES8QT` Web design system: the Claude Design system built from the owner's basics as one embedded stylesheet, embedded fonts and the shared layout and partials every page renders with
