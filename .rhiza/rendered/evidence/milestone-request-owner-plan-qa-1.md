# milestone-request-owner-plan-qa-1

evidence · `01M6N3YHYB83SSK8KKJ45EF262` · actor [claude-code-79d053d4](../actor/claude-code-79d053d4.md)

## result

PASS

## statement

At ac061b1b, with the built binary in a scratch repo (RHIZA_* unset, RHIZA_TEST_DETACH_LOG so nothing is spawned): queued a milestone request 'Ship v3', a feature 'Add search' and a bug 'Crash on save', ran rhiza daemon set rhiza.autonomy.plan auto and rhiza daemon --permission-mode manual --max-runs 4 --once. It started three plan runs; the detach log shows 'run start --foreground --workflow plan --new Ship v3 --request <id> --permission-mode manual --milestone' with no --autonomous, and '--new Add search ... --autonomous' and '--new Fix Crash on save ... --autonomous'. request ls shows all three planning.

## Links

### Out

- verifies [milestone-request-owner-plan](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/milestone-request-owner-plan.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
