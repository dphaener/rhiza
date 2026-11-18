# daemon-failed-retries-setting

decision · `01M7QQ8DMFKFPBPC5ZQ4GJRXCQ`

## reasoning

The intent's 'once by default, then the inbox'; a watchdog setting beside retryBudget follows how every other watchdog limit is set per machine and project

## statement

The watchdog retries a failed run whose failure the classifier calls real (or cannot judge) with rhiza run resume --retry, up to a new watchdog setting failedRetries (rhiza.daemon.watchdog.failedRetries and rhiza.project.watchdog.failedRetries, default 1, 0 turns it off), under the same enabled and autonomousOnly rules as its other passes; transient failures keep their own retryBudget. Each retry is journalled as a watchdog attempt; once failedRetries are spent the run stays failed and goes to the inbox, its item naming the retries tried, and the inbox lists no failed run while a retry is still due

## wording

(assumed, would have escalated) the conductor's ruling, not asked of the owner

## Links

### Out

- derived-from [m40-retry-conduct-plan](../source/m40-retry-conduct-plan.md)

### In

- decided-by [watchdog-retries-real-failure](../capability/rhiza/orchestration/daemon/daemon-watchdog/watchdog-retries-real-failure.md)

## History

- `01M7NVMMPK3A9SRVA9Q63N5SR9` A failed run is retried, once by default, from the step its workflow names for that failure (default the failed step; a plan whose review-check refuses restarts at plan), with the failed step's output and findings in the next attempt's brief; once the retries are spent it goes to the inbox
