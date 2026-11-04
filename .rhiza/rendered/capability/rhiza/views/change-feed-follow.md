# change-feed-follow

capability · `01M6B3PPF7X8K68YEWCN80ZXKB` · parent [views](../views.md)

## acceptance

A testscript (feed-follow.txtar) runs rhiza feed --follow --json in the background with no cursor and shows it emits nothing for change sets landed before it started, then one NDJSON FeedItem line per record change set landed after, a relanded change set again at its new land, and no actor-only record without --all; an interrupt stops it with exit 0; rhiza feed --follow --since <t> --limit 3 emits the backlog oldest first and exits 0 after three lines; rhiza feed --follow --cursor <next_cursor> --limit N picks up after that page; the human form prints '- ' lines only; and rhiza api call feed '{"follow":true,"since":...,"limit":N}' prints the same bytes as rhiza feed --follow --json --since ... --limit N. A Go test in internal/api drives follow with a replaced ticker and checks a land between ticks is emitted once.

## statement

rhiza feed --follow (FeedRequest field follow) streams each change set as it lands until interrupted: with neither --cursor nor --since it emits only change sets landing after it starts, with either it first emits the backlog after that point oldest land first; it refolds the log on a test-replaceable poll ticker (about 2 seconds) and emits each change set whose latest land entry is past the last one emitted, so a relanded change set appears again; --json prints one bare FeedItem per line (NDJSON, no page, no next_cursor), human output prints feed's '- ' lines with no 'Nothing landed.' and no next line; --all filters as feed does; --limit N stops after N emitted change sets and exits 0, unlimited by default; rhiza api call feed with {"follow":true} prints the same bytes as rhiza feed --follow --json through a stream sink

## Links

### Out

- decided-by [feed-follow-shape](../../../decision/feed-follow-shape.md)

### In

- verifies [change-feed-follow-verified](../../../evidence/change-feed-follow-verified.md)
- verifies [change-feed-follow-qa-1](../../../evidence/change-feed-follow-qa-1.md)
- verifies [pr-262-rhiza-feed-follow-streams](../../../evidence/pr-262-rhiza-feed-follow-streams.md)

## History

- `01M6J48HBGRJ04GMXM1CFWF3B3` rhiza feed --follow streams each change set as it lands (NDJSON with --json), reusing the feed cursor (internal/api/feed.go). Carve leaf change-feed-follow out of umbrella 01M5SWSZED0AS0 (Area views)
