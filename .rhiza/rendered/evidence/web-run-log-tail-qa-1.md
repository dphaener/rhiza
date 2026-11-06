# web-run-log-tail-qa-1

evidence · `01M6Q7JXVX0088FZ286SVQ8YZV` · actor [claude-code-1544fbc8](../actor/claude-code-1544fbc8.md)

## result

FAIL

## statement

At 6d5ec8e2 with the built rhiza's real daemon web UI: /api/run.tail is served (200); /fragments/run-log?run=<plan run> rendered the driver log's last lines in the rz-term block (Live) matching the file's tail; a run with no log rendered No log / Only a run the daemon started headless has a driver log; appending two lines to the log emitted event: log with the run's id on /events and the fragment then showed both new lines. But the run's page itself (/runs/<id>) answers 404 for every run because run.status is not read-only, so no run page shows the tail.

## Links

### Out

- verifies [web-run-log-tail](../capability/rhiza/web-monitor/web-run-log-tail.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
