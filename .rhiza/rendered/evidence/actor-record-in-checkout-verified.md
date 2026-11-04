# actor-record-in-checkout-verified

evidence · `01M6M7J2SF1X6X53MJDY94Z8H6` · actor [claude-code-2524aa23](../actor/claude-code-2524aa23.md)

## result

PASS

## statement

Verified at 920738fb by claude-code-2524aa23. Read the diff: runSessionStart now copies the session's actor record change set (copyChangeSetLog, the same files and IDs) from the run worktree into the session's own checkout when r.dir is set. Ran go test ./cmd/rhiza -run TestScripts/build-verifier-checkout -count=1: PASS. In it the verifier in its detached checkout finds the record change set, diff -r against the run worktree's copy shows no difference, rhiza show <actor> resolves, and rhiza add actor <actor> --record exits non-zero. A scratch copy of that scenario (outside the repo, nothing committed) showed the refusal message 'handle "claude-code-f3ff4d45-2" is taken by node ...'. It then committed the record in the checkout on a branch and the same record in the run worktree, merged them with no conflict, ran rhiza render and rhiza check (check=0, ok), and the merged log held one file with that actor handle, which rhiza show resolved. gofmt and go vet are clean on internal/api and cmd/rhiza. go test ./internal/api -short passes (ok, 433s); the full run under load average 17.8 failed only the TestReadBudget timing budget and the test timeout, both load-bound and not touched by this diff.

## Links

### Out

- verifies [actor-record-in-checkout](../capability/rhiza/runtime/actor-stamping/actor-record-in-checkout.md)

## History

- `01M6M2R7RJCMNMMKH0KS2364BK` Fix An agent session's actor record is written once per checkout, not once per run, so a ship run's branch ends up with two record change sets adding the same actor handle and its review/code step fails rhiza check with 'duplicate handle'. Seen twice on 2026-11-04: usage governor (claude-code-a463ece8) and text references (claude-code-83858632: recorded uncommitted in the run worktree as change set 01M6KE23, then again in the verifier checkout as 01M6KFMJ, both committed). Make actor stamping idempotent across a run's worktree and checkouts (look up the handle in the branch, the run worktree and main before adding), or have rhiza check/land merge identical actor records
