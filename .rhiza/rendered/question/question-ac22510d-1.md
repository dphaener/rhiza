# question-ac22510d-1

question · `01M6J4H0EM577MQSQJ2N0X1QDB`

## session

ac22510d-4d6f-4609-a2c4-01a1a444927a

## statement

Change set 'rhiza feed --follow streams' (01M6J48) serves no intent yet. Its leaf change-feed-follow (01M6B3PP) sits in the accepted Area views umbrella (01M5SWSZED0AS0), which serves anyone-builds and small-accurate-context; after the carve the umbrella keeps only its serves links. Which intent, and how does the leaf come over? / Where does `rhiza feed --follow` start when no --cursor or --since is given, and how does it pick up new lands? / What does each streamed line look like, and how do --limit and `rhiza api call feed` behave with follow?

## Links

### In

- answers [feed-follow-shape](../decision/feed-follow-shape.md)

## History

- `01M6J48HBGRJ04GMXM1CFWF3B3` rhiza feed --follow streams each change set as it lands (NDJSON with --json), reusing the feed cursor (internal/api/feed.go). Carve leaf change-feed-follow out of umbrella 01M5SWSZED0AS0 (Area views)
