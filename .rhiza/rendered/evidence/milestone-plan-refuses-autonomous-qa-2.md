# milestone-plan-refuses-autonomous-qa-2

evidence · `01M6N5M6FMB4SGYXED9NFB2RTS` · actor [claude-code-40877905](../actor/claude-code-40877905.md)

## result

PASS

## statement

At 54c67b81, built go build -trimpath ./cmd/rhiza and used it by hand in a scratch git repo (RHIZA_* unset, fake rex on PATH, a no-agent plan workflow). With a milestone request whose run the daemon had drawn, run start --workflow plan --new 'Ship v4' --request <it> --autonomous without --milestone exited 2 with '--autonomous is refused on a plan run that plans a milestone: milestone planning needs the owner, even in autonomous mode' (usage), queued, --foreground, --dry-run --json and via api call run.start; with --milestone too. No draft and no .git/rhiza/runs were written. The same request without --autonomous ran to succeeded with a journal not autonomous; a feature request's --request --autonomous run succeeded with autonomous true. run start <record change set> --autonomous and --new --milestone --autonomous each exited 2 the same way with change ls unchanged; after editing the milestone run's journal to autonomous:true, run resume exited 2 with the rule. The cycle-1 gap (milestone request kind not checked) is fixed. Scoped testscripts run-start-milestone-autonomous, plan-requests, plan-workflow-autonomous(-left) pass.

## Links

### Out

- verifies [milestone-plan-refuses-autonomous](../capability/rhiza/workflows/autonomous-runs/milestone-plan-refuses-autonomous.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
