# web-roadmap-row-titles-verified

evidence · `01M7AHVKZN5TN9KMYZ0B7SHCFH` · actor [claude-code-f93b5e91](../actor/claude-code-f93b5e91.md)

## result

PASS

## statement

At b8ca2bcb7 re-ran the Checks block: go test -trimpath ./internal/web/... ok; go test -trimpath ./internal/api/ -run 'Roadmap|Plan|Ready|Describe' ok; go test -trimpath ./cmd/rhiza -run 'TestScripts/(roadmap|plan|describe)' ok (describe-golden matches the regenerated title field); rhiza ship generate exit 0 (0 changed); rhiza render exit 0. TestRoadmapChangeSetsCarryTheirDisplayTitle passes, asserting roadmap --json titles are the set title ('Stream exports') or the intent's first clause ('Cache the roadmap'). TestRoadmapRowShowsTitleThenIntent passes, asserting /roadmap renders the short title as link text with the long intent in <small>, and rhiza.css clamps .rz-row-title small with -webkit-line-clamp: 2. Diff read: RoadmapChangeSet gains Title, web rows use Title/Sub=intent and keep depends-on as a separate Note line; no stubs or TODOs, nothing outside the leaf.

## Links

### Out

- verifies [web-roadmap-row-titles](../capability/rhiza/web-monitor/web-roadmap-page/web-roadmap-row-titles.md)

## History

- `01M79T7NW1B0KKXVVSX4AHZNHF` Fix The web roadmap shows each change set's whole intent sentence as its title, so rows run off and say little (owner, 2026-11-13: 'The titles seem to be a description and just run off. It's not informative at all.'). roadmap.go:140 sets Row.Title to c.Intent; change sets carry a short title (rhiza change new --title, at most 30 characters, defaulting to the intent's first clause). Wanted: rows show the change set's title, with the intent as the secondary line clamped to two lines, on the roadmap and wherever else the web UI lists change sets (overview, runs list), with a test.
