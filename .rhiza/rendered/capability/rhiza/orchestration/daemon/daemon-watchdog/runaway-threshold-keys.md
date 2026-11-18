# runaway-threshold-keys

capability · `01M7QN2RE9HAZ65J0A2DABCHG1` · parent [daemon-watchdog](../daemon-watchdog.md)

## acceptance

go test -trimpath -p 2 ./internal/api -run 'TestWatchdogPolicy' passes, a test resolving the three keys' defaults (400, 8G as bytes, 5m) and a project key over a machine key; env -u RHIZA_SESSION go test -trimpath -p 2 ./cmd/rhiza -run TestScripts/watchdog-policy passes with the three keys listed at their defaults and one set through rhiza daemon set. (failure case) rhiza daemon set rhiza.daemon.watchdog.runawayMemory lots and rhiza daemon set rhiza.project.watchdog.runawayProcs -1 are refused naming the key; docs/reference/config.md is regenerated (rhiza ship generate) so the docs check passes

## files

internal/api/watchdog_policy.go, internal/api/watchdog_policy_test.go, cmd/rhiza/testdata/script/watchdog-policy.txtar, docs/reference/config.md

## must_not

Change any existing watchdog key, its default or meaning; read a process table; act on the keys (the runaway pass's leaf)

## statement

The watchdog policy gains runawayProcs, runawayMemory and runawayGrace, machine (rhiza.daemon.watchdog.*) and project (rhiza.project.watchdog.*) keys with their defaults, listed by rhiza daemon settings and documented in the generated docs/reference/config.md

## Links

### Out

- decided-by [runaway-thresholds-defaults](../../../../../decision/runaway-thresholds-defaults.md)

### In

- depends-on [runaway-session-fails-step](runaway-session-fails-step.md)
- verifies [runaway-threshold-keys-verified](../../../../../evidence/runaway-threshold-keys-verified.md)
- verifies [runaway-threshold-keys-qa-1](../../../../../evidence/runaway-threshold-keys-qa-1.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason
