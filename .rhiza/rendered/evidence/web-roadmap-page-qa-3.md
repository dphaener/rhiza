# web-roadmap-page-qa-3

evidence · `01M6Q9SE74H84SDA3D7XMMYCGE` · actor [claude-code-4b0924b9](../actor/claude-code-4b0924b9.md)

## result

PASS

## statement

At e54ae787, on the same scratch daemon (port 18943, no current change set), GET /roadmap returned 200 and links each of its 16 change sets to /changes/<full id>; following those links for an accepted, a draft, a withdrawn and two landed change sets gave their pages with 200, so the roadmap's links now land on working change set pages. The group order and umbrella labelling were checked in cycle 2 and this cycle's fix does not touch internal/web/roadmap.go.

## Links

### Out

- verifies [web-roadmap-page](../capability/rhiza/web-monitor/web-roadmap-page.md)

## History

- `01M6PFH625P3BYQAB8P6J38BBT` Web Roadmap pages: milestone groups and the backlog, and a change set's page with its intent, state, leaves, questions, runs and history
