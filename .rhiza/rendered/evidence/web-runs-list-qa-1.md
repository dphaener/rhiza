# web-runs-list-qa-1

evidence · `01M6Q7JXJ0HKW0200KG35D3ZQA` · actor [claude-code-1544fbc8](../actor/claude-code-1544fbc8.md)

## result

PASS

## statement

At 6d5ec8e2 built rhiza (go build -trimpath) and ran its real daemon (paused, rhiza.web.port 18777) in a scratch clone holding copies of this repo's 100 run journals. GET /runs answered 200: a Recent group of 50 rows newest first (Fix The daemon never reaps its, Web UI Roadmap pages, ...), each linking to /runs/<id>, with its change set title, workflow (ship-change/plan), state pill, current step (build/build) and start time (Nov 6 01:13), no ID text; the list section declares data-rz-events=run and /fragments/runs. With one older run's lock held by flock, /fragments/runs put it alone under Live (running pill) before Recent. Appending to a journal emitted event: run on /events.

## Links

### Out

- verifies [web-runs-list](../capability/rhiza/web-monitor/web-runs-list.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
