# byo-watchdog-stall-runaway

capability · `01M7RB624W6QN05GB3VGZPNE90` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath -p 2 ./internal/api -run 'TestStallLimit|TestWatchdog|TestRunaway' passes, with new cases: StallLimit("gate", 0, 0) returns CommandLimit and an explicit stall still wins; a myteam-flow-shaped run silent in step work is nudged then stopped and resumed with --retry; one silent in step gate past commandLimit is stopped and resumed; a runaway work session (process count past the threshold from a fake process table) fails the step with the runaway reason; (failure case) a gate step still inside commandLimit is neither stopped nor nudged; rhiza docs regenerates docs/reference/config.md with commandLimit described as covering command and gate steps

## files

internal/api/watchdog_policy.go, internal/api/watchdog_policy_test.go, internal/api/daemon_watchdog_test.go, internal/api/watchdog_runaway_test.go, docs/reference/config.md

## must_not

Change the built-in workflows' stall= values; special-case a workflow or step name; read the real process table in a test; run go test ./internal/api unscoped more than once

## statement

The watchdog's stall and runaway detection apply to a BYO workflow's steps: a gate step with no stall= gets the command limit, and Go tests over a BYO-shaped workflow show its agent step nudged then stopped and retried, its gate step stopped and retried past the command limit, and a runaway session of its agent step treated as a stall with the reason

## Links

### Out

- decided-by [byo-gate-command-limit](../../../../decision/byo-gate-command-limit.md)
- decided-by [byo-runaway-go-test](../../../../decision/byo-runaway-go-test.md)

### In

- verifies [byo-watchdog-stall-runaway-verified](../../../../evidence/byo-watchdog-stall-runaway-verified.md)
- verifies [byo-watchdog-stall-runaway-qa-1](../../../../evidence/byo-watchdog-stall-runaway-qa-1.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
