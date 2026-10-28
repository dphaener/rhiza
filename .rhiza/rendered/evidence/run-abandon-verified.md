# run-abandon-verified

evidence · `01M624J5ME5N668FV7YYCNZPC7` · actor [claude-code-8989bf11](../actor/claude-code-8989bf11.md)

## result

PASS

## statement

At e2a7780 ran go test ./cmd/rhiza -run 'TestScripts/(run-abandon|describe-golden)$' -count=1: both PASS. run-abandon.txtar stops a fake hold-workflow run at its wait, abandons it, and checks that run status shows it failed with detail 'abandoned by test' (the actor), that git worktree list still has 2 entries and the one cs/* branch remains, that rhiza run start on the change set then succeeds, that run resume and resume --retry exit 2 with '<id> was abandoned', and that abandoning the finished run (and the abandoned one) exits 2 with 'is finished'. Re-running describe-golden with -update left the tree clean, so the golden is current. make check (RHIZA_ROOT/CHANGE/RUN unset) exited 0. Read the diff: run.abandon op in internal/api/run.go appends a Failed finished entry under the run lock, re-reads to refuse an ended run, and touches no worktree or branch; State.Abandoned in internal/run/journal.go gates runRunResume before the retry branch.

## Links

### Out

- verifies [run-abandon](../capability/rhiza/workflows/workflow-runs/run-abandon.md)

## History

- `01M5ZNV4E4AJR52EX631180DMS` A change set has at most one live run: rhiza run start refuses a second run on a change set whose run is still going, naming that run.
