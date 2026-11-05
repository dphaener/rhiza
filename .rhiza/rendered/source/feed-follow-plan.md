# feed-follow-plan

source · `01M6J4ZSA47CRX58E4RAEJ5C9W`

## quotes

```json
[]
```

## raw_sha256

bc8c2a646bee935a35782bf603e5596edde6c8fc714fcdfe325722b0e6cf8004

## session

ac22510d-4d6f-4609-a2c4-01a1a444927a

## summary

Autonomous plan run on change set 01M6J48 (no owner present). The planner read internal/api/feed.go (one-shot paged feed over the folded log, opaque base64url land-entry cursor, actor-only filter unless --all), found no streaming operation in the API registry, and the daemon's poll loop with a test-replaceable ticker and Ctx.Interrupts as the precedent for running until stopped. It raised one question (question-ac22510d-1) on intent and carve, start point and pick-up, and line shape with --limit and api call. The arbiter answered with Decision feed-follow-shape: serve milestone-25 and carve change-feed-follow out of the Area views umbrella (which keeps its serves links); follow is a FeedRequest field (--follow); without --cursor or --since it emits only change sets landing after it starts, with them it emits the backlog first; it refolds on a ~2s replaceable poll ticker and emits each change set whose latest land entry is past the last emitted, so a relanded change set reappears; runs until interrupted, exit 0; --json lines are bare FeedItem objects (NDJSON), human lines are feed's '- ' lines with no 'Nothing landed.' or next line; --all filters as now; --limit N stops after N emitted change sets; `rhiza api call feed` with follow:true streams the same bytes through a stream sink. The planner then carved the leaf, linked the change set serves milestone-25 and rewrote the leaf's statement and acceptance from that Decision.

## title

rhiza feed --follow streams, planned in an autonomous run

## Links

### In

- derived-from [feed-follow-shape](../decision/feed-follow-shape.md)

## History

- `01M6J48HBGRJ04GMXM1CFWF3B3` rhiza feed --follow streams each change set as it lands (NDJSON with --json), reusing the feed cursor (internal/api/feed.go). Carve leaf change-feed-follow out of umbrella 01M5SWSZED0AS0 (Area views)
