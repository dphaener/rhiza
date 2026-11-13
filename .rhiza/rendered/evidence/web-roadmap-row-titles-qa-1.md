# web-roadmap-row-titles-qa-1

evidence · `01M7AK47EVPPFTDNTHXEJG05EC` · actor [claude-code-f8bb7f6e](../actor/claude-code-f8bb7f6e.md)

## result

PASS

## statement

QA at be2335f89: built ./cmd/rhiza to a scratch path and ran it in a scratch repo under a scratch Rhiza root (env -i). Happy: roadmap --all --json gives title 'Stream exports' for a set title, 'Cache the roadmap' (intent's first clause) and 'No comma at all in this rather' (30-char cut) when none is set; roadmap --json and api call roadmap print the same bytes. A paused daemon serving the web UI rendered /p/repo/roadmap and /p/repo/fragments/roadmap with each row's link text the title and the full intent in <small> beneath; served rhiza.css clamps .rz-row-title small to 2 lines; a change set with a depends-on link keeps 'depends on 1 change set' as its own small.rz-row-note line; HTML in title and intent is escaped; row links return 200. Failure: /p/nope/roadmap 404, roadmap --milestone exits 2 with a usage error. /p/repo/ and /p/repo/runs still render 200. Scoped tests: go test ./internal/web/... ok, ./internal/api -run Roadmap ok, ./cmd/rhiza -run TestScripts/(describe-golden|roadmap) ok.

## Links

### Out

- verifies [web-roadmap-row-titles](../capability/rhiza/web-monitor/web-roadmap-page/web-roadmap-row-titles.md)

## History

- `01M79T7NW1B0KKXVVSX4AHZNHF` Fix The web roadmap shows each change set's whole intent sentence as its title, so rows run off and say little (owner, 2026-11-13: 'The titles seem to be a description and just run off. It's not informative at all.'). roadmap.go:140 sets Row.Title to c.Intent; change sets carry a short title (rhiza change new --title, at most 30 characters, defaulting to the intent's first clause). Wanted: rows show the change set's title, with the intent as the secondary line clamped to two lines, on the roadmap and wherever else the web UI lists change sets (overview, runs list), with a test.
