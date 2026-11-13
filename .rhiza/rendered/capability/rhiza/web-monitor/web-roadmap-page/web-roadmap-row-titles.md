# web-roadmap-row-titles

capability · `01M79TEDZ8NQS65M1CQQ1DSE52` · parent [web-roadmap-page](../web-roadmap-page.md)

## acceptance

An internal/api test asserts roadmap --json gives each change set its display title (a set title, and the intent's first clause when none is set); an internal/web roadmap test renders /roadmap over a fixture whose change set has a long intent and a short title and asserts the row's link text is the title and the intent sits in the secondary line, and the stylesheet clamps that line to two lines; describe-golden.txtar is regenerated

## statement

Each change set row on the web roadmap shows the change set's title (as change ls shows it), with its intent beneath as a secondary line clamped to two lines, and any depends-on note kept on the row; the roadmap operation's change sets carry that title as a title field

## Links

### Out

- decided-by [roadmap-rows-by-title](../../../../decision/roadmap-rows-by-title.md)

### In

- verifies [web-roadmap-row-titles-verified](../../../../evidence/web-roadmap-row-titles-verified.md)
- verifies [web-roadmap-row-titles-qa-1](../../../../evidence/web-roadmap-row-titles-qa-1.md)
- verifies [pr-546-fix-the-web-roadmap-shows-each](../../../../evidence/pr-546-fix-the-web-roadmap-shows-each.md)

## History

- `01M79T7NW1B0KKXVVSX4AHZNHF` Fix The web roadmap shows each change set's whole intent sentence as its title, so rows run off and say little (owner, 2026-11-13: 'The titles seem to be a description and just run off. It's not informative at all.'). roadmap.go:140 sets Row.Title to c.Intent; change sets carry a short title (rhiza change new --title, at most 30 characters, defaulting to the intent's first clause). Wanted: rows show the change set's title, with the intent as the secondary line clamped to two lines, on the roadmap and wherever else the web UI lists change sets (overview, runs list), with a test.
