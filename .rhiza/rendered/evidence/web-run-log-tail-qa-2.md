# web-run-log-tail-qa-2

evidence · `01M6QAB278FMQ2RW3AF8KP9XDE` · actor [claude-code-ec8909b8](../actor/claude-code-ec8909b8.md)

## result

PASS

## statement

At 73d2ed3f with the same built daemon web UI: /runs/<plan run> (200) shows a Log section with the driver log's last lines in the rz-term block (log · tail, live dot) wired to /fragments/run-log?run=<id> and data-rz-events=log; appending two lines to that run's driver log emitted event: log with the run's id on /events and the fragment then ended with both new lines; /runs/<run without a log> shows No log / Only a run the daemon started headless has a driver log.

## Links

### Out

- verifies [web-run-log-tail](../capability/rhiza/web-monitor/web-run-log-tail.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
