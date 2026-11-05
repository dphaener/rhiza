# change-feed-follow-verified

evidence · `01M6MC1QEGZ27VRE618NY486ZM` · actor [claude-code-a6ba7896](../actor/claude-code-a6ba7896.md)

## result

PASS

## statement

At 3937eb3a I read the diff (feed.go runFeedFollow, Ctx.Out/JSON sink, Streamer in cmd/rhiza call) and ran: go test -trimpath -count=1 ./internal/api (ok, incl. TestFeedFollowTicks: replaced ticker, a land between ticks emitted once as one NDJSON line, stops on interrupt; TestFeedFollowLimit), -race -count=5 on TestFeedFollow (ok), and go test ./cmd/rhiza -run TestScripts/(feed-follow|feed|describe-golden) (pass; feed-follow -count=3 pass). feed-follow.txtar covers each acceptance point: background --follow --json with no cursor emits only delta, the relanded Store things at its 2026-11-01 land, and epsilon, never alpha/beta/gamma nor the actor-only record; kill -INT then wait exits 0; --since --limit 3 prints three '- ' lines oldest first and nothing else; --cursor from a limit-2 page picks up gamma; --all adds the actor record; api call feed {follow,since,limit} cmp-equal to feed --follow --json. Limit is unlimited by default in follow (q.Limit used, not the page default). go vet and gofmt clean.

## Links

### Out

- verifies [change-feed-follow](../capability/rhiza/views/change-feed-follow.md)

## History

- `01M6J48HBGRJ04GMXM1CFWF3B3` rhiza feed --follow streams each change set as it lands (NDJSON with --json), reusing the feed cursor (internal/api/feed.go). Carve leaf change-feed-follow out of umbrella 01M5SWSZED0AS0 (Area views)
