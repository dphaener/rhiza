# actor-record-in-checkout-verified-2

evidence · `01M6M9MCP3KDDY6BZQBPFRDMS3` · actor [claude-code-9a3620ea](../actor/claude-code-9a3620ea.md)

## result

PASS

## statement

At e7e52b13 ran go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/build-verifier-checkout: PASS. The fake verifier runs in its own detached checkout; diff -r of the session's actor record change set between RHIZA_ROOT and the checkout shows no difference, rhiza show <actor> resolves there and rhiza add actor <actor> --record is refused (captured in a scratch copy of the script: 'handle "claude-code-a6fe1895-2" is taken by node ...'). The checkout commits its copy on branch checkout-copy, the script commits the build worktree's log, merges both in a fresh worktree with a clean status, and rhiza render + rhiza check print ok, with rhiza ls --kind actor --json counting exactly one node for the handle; fold-ok passes. Read the diffs of 920738fb (session start copies the actor record change set into r.dir via copyChangeSetLog, same files and IDs) and e7e52b13 (the merge step that closes review finding f1). make check: fmt, vet and all go tests pass; self-check fails only on stale rendered actor files from the run's own chore record commits, not from this leaf.

## Links

### Out

- verifies [actor-record-in-checkout](../capability/rhiza/runtime/actor-stamping/actor-record-in-checkout.md)

## History

- `01M6M2R7RJCMNMMKH0KS2364BK` Fix An agent session's actor record is written once per checkout, not once per run, so a ship run's branch ends up with two record change sets adding the same actor handle and its review/code step fails rhiza check with 'duplicate handle'. Seen twice on 2026-11-04: usage governor (claude-code-a463ece8) and text references (claude-code-83858632: recorded uncommitted in the run worktree as change set 01M6KE23, then again in the verifier checkout as 01M6KFMJ, both committed). Make actor stamping idempotent across a run's worktree and checkouts (look up the handle in the branch, the run worktree and main before adding), or have rhiza check/land merge identical actor records
