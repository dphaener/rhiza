# feed-follow-shape

decision · `01M6J4QEQC82PV59ATK5WART09`

## reasoning

Intent: the umbrella 01M5SWSZED0AS0 serves anyone-builds, small-accurate-context and the umbrella intent areas; milestone-25 says every change set planned from the backlog serves it, and its served-by list already holds the sibling carves from Area umbrellas (usage-governor, anchors, workflow-authoring, text-references), so this one follows the same pattern. The carve guide says carve keeps IDs and links and needs both change sets unlanded, which holds (umbrella accepted, this one draft). Behaviour: internal/api/feed.go already orders by land entry and encodes the cursor as the land entry ID, so following is that same filter (land entry > last emitted) repeated on a poll; the daemon polls on a ticker that tests replace (daemonTicker), which keeps the test rule of no sockets and deterministic runs. CLAUDE.md requires `rhiza <cmd> --json` and `rhiza api call <op>` to reach the same handler and print the same bytes, and the CLI to be a generic frame over request tags, so follow belongs on FeedRequest rather than as a CLI loop. Starting at the end is what follow means in tail -f and rex events --follow, and --cursor/--since already give a way to replay. NDJSON of plain FeedItems was named in the change set's intent. A default page limit of 100 would cut a follow off, so follow drops the default while an explicit limit becomes a stop count, which gives tests a way to end.

## statement

Intent and carve: this change set serves milestone-25 (the backlog is empty), as the other change sets carved out of Area umbrellas do; leaf change-feed-follow comes over by `rhiza carve change-feed-follow --from 01M5SWSZED0AS0 --to 01M6J48`, its acceptance replaced with a real check, and the umbrella keeps its serves links. Start and pick-up: follow is a field of FeedRequest (`follow`, flag --follow), not a CLI-only flag. With neither --cursor nor --since it starts at the current end of the feed and emits only change sets that land after it starts (like tail -f); with --cursor or --since it first emits the backlog after that point, oldest land first, then keeps going. It picks up new lands by re-folding the log on a poll ticker (about 2 seconds, replaceable in tests the way daemonTicker is), emitting each change set whose latest land entry is past the last one emitted, so a reopened and relanded change set appears again at its new land. It runs until interrupted (exit 0). Lines: with --json each line is one FeedItem object exactly as it appears in change_sets (NDJSON), with no wrapping page and no next_cursor; without --json each line is the same '- <landed_at> <id> <kind> <actor>: <title> (<n> leaves) #<pr>' line as Human, with no 'Nothing landed.' and no next line. --all filters as it does now. Limit: with --follow the default is unlimited; an explicit --limit N stops after N emitted change sets (backlog included) and exits 0, so scripts and tests end. `rhiza api call feed` with {"follow":true} streams the same NDJSON bytes as `rhiza feed --follow --json`, through a stream sink the operation writes to, so the two still print the same bytes.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-ac22510d-1](../question/question-ac22510d-1.md)
- derived-from [feed-follow-plan](../source/feed-follow-plan.md)

### In

- decided-by [change-feed-follow](../capability/rhiza/views/change-feed-follow.md)

## History

- `01M6J48HBGRJ04GMXM1CFWF3B3` rhiza feed --follow streams each change set as it lands (NDJSON with --json), reusing the feed cursor (internal/api/feed.go). Carve leaf change-feed-follow out of umbrella 01M5SWSZED0AS0 (Area views)
