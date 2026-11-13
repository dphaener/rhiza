# pr-546-fix-the-web-roadmap-shows-each

evidence · `01M7AM3CPZAZS9M1YVNP8ASY1K` · actor [dphaener](../actor/dphaener.md)

## head

ae98d902ac93de6add0feb73e1e27666f7bc3700

## pr

https://github.com/dphaener/rhiza/pull/546

## result

9 checks passed: changes SUCCESS, go-install SUCCESS, model SUCCESS, lint SKIPPED, scripts SKIPPED, api SKIPPED, rest SKIPPED, budgets SKIPPED, check SUCCESS

## statement

Pull request #546 (https://github.com/dphaener/rhiza/pull/546) passed CI at ae98d90

## status

verified

## Links

### Out

- verifies [web-roadmap-row-titles](../capability/rhiza/web-monitor/web-roadmap-page/web-roadmap-row-titles.md)

## History

- `01M79T7NW1B0KKXVVSX4AHZNHF` Fix The web roadmap shows each change set's whole intent sentence as its title, so rows run off and say little (owner, 2026-11-13: 'The titles seem to be a description and just run off. It's not informative at all.'). roadmap.go:140 sets Row.Title to c.Intent; change sets carry a short title (rhiza change new --title, at most 30 characters, defaulting to the intent's first clause). Wanted: rows show the change set's title, with the intent as the secondary line clamped to two lines, on the roadmap and wherever else the web UI lists change sets (overview, runs list), with a test.
