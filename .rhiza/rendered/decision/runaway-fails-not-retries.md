# runaway-fails-not-retries

decision · `01M7QN2RB1D66AZE1ZHE62A7YX`

## statement

A runaway stall session is nudged once with its reading; when a poll at least runawayGrace after the nudge still reads it over a threshold, the watchdog closes the stall sessions' panes (and parents waiting on their spawns), terminates the run and journals the current step stopped with the reading as its detail, never resuming or retrying it (a retry reruns the same runaway), and files a bug when fileBugs is on. A reading back under every threshold clears the nudge. An attended run under autonomousOnly gets a log line only, as stalls do

## wording

(assumed, would have escalated) planner's choice under the autonomous rules, from the intent's 'nudges the session, then fails the step with the reading as its reason'

## Links

### Out

- derived-from [runaway-plan-notes](../source/runaway-plan-notes.md)

### In

- decided-by [runaway-session-fails-step](../capability/rhiza/orchestration/daemon/daemon-watchdog/runaway-session-fails-step.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason
