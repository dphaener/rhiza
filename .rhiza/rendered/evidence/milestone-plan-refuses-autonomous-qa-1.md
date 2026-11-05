# milestone-plan-refuses-autonomous-qa-1

evidence · `01M6N3YAAGEJW96EXE9MAQDS03` · actor [claude-code-79d053d4](../actor/claude-code-79d053d4.md)

## result

PASS

## statement

At ac061b1b, built go build -trimpath -o $TMPDIR/rhiza-qa ./cmd/rhiza and used it by hand in a scratch git repo (RHIZA_* unset, fake rex on PATH, built-in plan workflow). run start --workflow plan --new 'Plan the milestone' --milestone --autonomous exited 2 with '--autonomous is refused on a plan run that plans a milestone: milestone planning needs the owner, even in autonomous mode', with and without --foreground and with --dry-run --json (code usage); run start <record change set> --workflow plan --autonomous, queued or --foreground, and api call run.start with autonomous:true on it, each exited 2 the same way; the change set count stayed 2 and no .git/rhiza/runs was written. Without --autonomous both milestone starts dry-run fine (exit 0), and --new 'Plan an intention' --autonomous dry-runs with autonomous true. With a no-agent plan.toml, a real run on the record change set succeeded; after editing its journal's started entry to autonomous:true, run resume and run resume --retry each exited 2 with the same error (before the edit resume said it was finished, exit 3). Gap: run start --new 'Ship v3' --request <a milestone-kind request> --autonomous without --milestone dry-runs with autonomous true on an intention change set (with --milestone it is refused): the request's kind is not checked, reported as a finding.

## Links

### Out

- verifies [milestone-plan-refuses-autonomous](../capability/rhiza/workflows/autonomous-runs/milestone-plan-refuses-autonomous.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
