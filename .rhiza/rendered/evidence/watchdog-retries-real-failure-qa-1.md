# watchdog-retries-real-failure-qa-1

evidence · `01M7QTPW6BX2PYCDC9KB0V6YDY` · actor [claude-code-b8e7291f](../actor/claude-code-b8e7291f.md)

## result

PASS

## statement

QA with a scratch build: rhiza daemon settings lists rhiza.daemon.watchdog.failedRetries and rhiza.project.watchdog.failedRetries (default 1); daemon set accepts 0 and 3 and refuses -1, abc, 1.5 and empty; inbox keeps a failed run while no daemon runs; go test ./internal/api -run 'Watchdog|Inbox' and TestScripts/watchdog-policy pass

## Links

### Out

- verifies [watchdog-retries-real-failure](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-retries-real-failure.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox
