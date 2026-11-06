# web-roadmap-page-qa-2

evidence · `01M6Q91VK2M9QV1MEMZ75AGTK2` · actor [claude-code-0ebbd179](../actor/claude-code-0ebbd179.md)

## result

PASS

## statement

At 1394f316 I built rhiza (go build -trimpath), cloned this repo to a scratch dir and ran rhiza daemon --max-runs 0, paused, with rhiza.web.port 18931 and no current change set, then fetched pages with curl. /api/roadmap?all=true returns milestone, milestone, backlog, umbrella. GET /roadmap shows them in that order: both milestones titled by their intents' statements, then the Backlog, then the umbrella labelled Umbrella with its 10 Area change sets after the backlog. The kicker reads '2 milestones · 16 change sets · 1 in backlog', so the umbrella is no longer counted as a milestone. Every change set row shows its intent and a state pill, links to /changes/<full id>, and shows no handle or ID. The cycle-1 finding is fixed. One nit remains: the backlog aside still says '1 change sets'.

## Links

### Out

- verifies [web-roadmap-page](../capability/rhiza/web-monitor/web-roadmap-page.md)

## History

- `01M6PFH625P3BYQAB8P6J38BBT` Web Roadmap pages: milestone groups and the backlog, and a change set's page with its intent, state, leaves, questions, runs and history
