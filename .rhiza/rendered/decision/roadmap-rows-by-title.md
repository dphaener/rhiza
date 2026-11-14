# roadmap-rows-by-title

decision · `01M79TEDXTJHD0DFEP70R36S0M`

## quote

The titles seem to be a description and just run off. It's not informative at all.

## statement

Web roadmap rows are named by the change set's short title (its display title: the title, else the intent's first clause), with the intent as a secondary line clamped to two lines; the roadmap operation carries that title for each change set

## Links

### Out

- derived-from [plan-roadmap-row-titles](../source/plan-roadmap-row-titles.md)

### In

- decided-by [web-roadmap-row-titles](../capability/rhiza/web-monitor/web-roadmap-page/web-roadmap-row-titles.md)

## History

- `01M79T7NW1B0KKXVVSX4AHZNHF` Fix The web roadmap shows each change set's whole intent sentence as its title, so rows run off and say little (owner, 2026-11-13: 'The titles seem to be a description and just run off. It's not informative at all.'). roadmap.go:140 sets Row.Title to c.Intent; change sets carry a short title (rhiza change new --title, at most 30 characters, defaulting to the intent's first clause). Wanted: rows show the change set's title, with the intent as the secondary line clamped to two lines, on the roadmap and wherever else the web UI lists change sets (overview, runs list), with a test.
