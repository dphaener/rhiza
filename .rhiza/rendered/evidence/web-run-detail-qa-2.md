# web-run-detail-qa-2

evidence · `01M6QAB22YPEFF5W0J4AHGKP3E` · actor [claude-code-ec8909b8](../actor/claude-code-ec8909b8.md)

## result

PASS

## statement

At 73d2ed3f built rhiza (go build -trimpath) and ran its real daemon (paused, rhiza.web.port 18779) in a scratch clone holding copies of this repo's 101 run journals and 93 driver logs. GET /runs/<id> answered 200 for a ship-change run with a review step, a plan run and a succeeded ship-change run, and 404 for an unknown run. The ship-change page shows its steps in order (build/build, build/render, build/check, review/code, review/review) with done pills, durations, times and details (3 of 3 leaves passed, exit 0), each step's sessions (Session 1 finished · 1141274 tokens, matching input+output+cache_creation+cache_read from /api/run.status), and a Review section per check (compliance passed, quality open, qa open, scope passed) with each finding's severity, reason, place and confirmed verdict; no node, run or finding ID in the visible text. The plan run shows interview through close-out in order. The detail section declares data-rz-events=run and /fragments/run-detail?run=<id> (200); appending a line to the run's journal emitted event: run with its id on /events.

## Links

### Out

- verifies [web-run-detail](../capability/rhiza/web-monitor/web-run-detail.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
