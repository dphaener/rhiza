# web-roadmap-page-qa-1

evidence · `01M6Q7HG5P1DBQPKHE5YCYY46W` · actor [claude-code-3eb24189](../actor/claude-code-3eb24189.md)

## result

FAIL

## statement

At 75d6f69b I built rhiza (go build -trimpath) and ran rhiza daemon --max-runs 0, paused, with rhiza.web.port 18917, in a scratch clone of this repo, and fetched pages with curl. GET /roadmap returned 200 with each change set's title, state pill and href=/changes/<full id>, no handle or ID in the text, milestones titled by their intent's statement, and data-rz-fragment=/fragments/roadmap. Accepting a change set in the clone sent event: change on /events, and /fragments/roadmap then showed it, so the page is live. But /api/roadmap?all=true returns milestone-27, milestone-28, backlog, then the umbrella group 'areas'. The page counts that umbrella as a milestone ('3 milestones · 16 change sets · 1 in backlog') and moves it ahead of the backlog, so the page does not group them the way rhiza roadmap does. The backlog aside says '1 change sets'.

## Links

### Out

- verifies [web-roadmap-page](../capability/rhiza/web-monitor/web-roadmap-page.md)

## History

- `01M6PFH625P3BYQAB8P6J38BBT` Web Roadmap pages: milestone groups and the backlog, and a change set's page with its intent, state, leaves, questions, runs and history
