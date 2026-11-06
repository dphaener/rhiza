# web-change-set-page-qa-3

evidence · `01M6Q9SDTW8TAFDPWTJF5RK9HF` · actor [claude-code-4b0924b9](../actor/claude-code-4b0924b9.md)

## result

PASS

## statement

At e54ae787 I built rhiza (go build -trimpath), cloned this repo to a scratch dir and ran rhiza daemon --max-runs 0, paused, with rhiza.web.port 18943 and no current change set, then fetched pages with curl. GET /changes/01M6PFH625P3BYQAB8P6J38BBT (accepted) shows its title, intent, kind, state, both leaves with statement, acceptance and 'has evidence' ('2 leaves, 2 with evidence'), what it serves, Depends on 'Web UI server' and 'Web UI design system' (landed, linking to /changes/<id>), and its Proposed and Accepted history; after copying this run's journal into the clone, Runs lists the ship-change run linking to /runs/<id> with time 2026-11-06 00:57. Cycle 2's finding is fixed: the landed change sets 01M6J965 and 01M6FENG now show 3 leaves, 3 with evidence and 1 leaf, 1 with evidence, with statements, acceptance and Serves, matching change.show, and their Landed history. In the clone I drafted a change set with a capability and a question blocking it: its page shows the leaf with 'no evidence yet' and the question under Open questions; after rhiza withdraw it shows 0 leaves, as change.show does, and a Withdrawn entry. Unknown ULID and non-ID refs return 404, POST returns 405. go test -trimpath -count=1 ./internal/web/... is ok.

## Links

### Out

- verifies [web-change-set-page](../capability/rhiza/web-monitor/web-change-set-page.md)

## History

- `01M6PFH625P3BYQAB8P6J38BBT` Web Roadmap pages: milestone groups and the backlog, and a change set's page with its intent, state, leaves, questions, runs and history
