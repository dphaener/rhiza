# milestone-request-owner-plan-qa-2

evidence · `01M6N5MBMJNEKECDWHXYEXB988` · actor [claude-code-40877905](../actor/claude-code-40877905.md)

## result

PASS

## statement

At 54c67b81, with the built binary in a scratch repo, rhiza.autonomy.plan=auto, queued a milestone request 'Ship v6' and a bug 'Crash on save' and ran rhiza daemon --permission-mode manual --max-runs 2 --once (RHIZA_TEST_NOW and RHIZA_TEST_DETACH_LOG set so nothing spawned). The detach log shows 'run start --foreground --workflow plan --new Ship v6 --request <id> --permission-mode manual --milestone' with no --autonomous, and 'run start ... --new Fix Crash on save ... --autonomous'. Starting a milestone request's run as the daemon does leaves its journal with no autonomous:true.

## Links

### Out

- verifies [milestone-request-owner-plan](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/milestone-request-owner-plan.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
