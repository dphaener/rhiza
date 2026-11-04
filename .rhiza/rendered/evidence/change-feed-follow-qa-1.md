# change-feed-follow-qa-1

evidence · `01M6MCDBN5AQVSYC8ENYZCY28H` · actor [claude-code-e16c5f61](../actor/claude-code-e16c5f61.md)

## result

PASS

## statement

QA at 741ce893 with a binary built from the branch (go build ./cmd/rhiza) in a scratch git repo with real wall-clock times. rhiza feed --follow --since <t> --limit 2 printed the two oldest '- ' lines and exited 0 with no 'Nothing landed.' or next line; --follow --json --all --since --limit 3 printed three bare FeedItem lines (actor record included under --all, no change_sets or next_cursor) and rhiza api call feed '{"follow":true,"all":true,"since":...,"limit":3}' was cmp-identical; --follow --cursor <next_cursor of a limit-1 page> --limit 2 picked up beta and gamma, exit 0. Live: feed --follow and feed --follow --json started with no cursor emitted nothing already landed, then delta and eps as each landed with --record (~2s poll), and SIGINT exited 0. --follow --since then streamed the backlog oldest first. After reopen and reland of an intention change set, feed --follow --json and api call feed '{"follow":true}' each emitted it once at its new land entry (leaves 2), byte-identical. A bad --cursor exits 2 with a clear message. Scoped go test ./internal/api -run Feed and TestScripts/(feed-follow|describe-golden) pass. Note: with no new lands, a follow whose reader closed (| head -1) keeps running until the next write, as tail -f does.

## Links

### Out

- verifies [change-feed-follow](../capability/rhiza/views/change-feed-follow.md)

## History

- `01M6J48HBGRJ04GMXM1CFWF3B3` rhiza feed --follow streams each change set as it lands (NDJSON with --json), reusing the feed cursor (internal/api/feed.go). Carve leaf change-feed-follow out of umbrella 01M5SWSZED0AS0 (Area views)
