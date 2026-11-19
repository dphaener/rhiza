# conductor-land-commit-not-hand-commit

capability · `01M7RJTNWYVVR0Z2YYZ1WXY2Q4` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/api -run TestHandCommitsConductorLandCommit passes: on a run branch (handCommitBranch), the land operation with commit set, run as the actor of a session recorded with the run (Run R1), makes a commit and scanHand finds no hand commit; (failure case) the same land commit made as an actor whose only session belongs to another run (R2) is journalled as one hand-commit intervention naming that actor

## files

internal/api/run_handcommit_test.go

## must_not

Change internal/api/run_handcommit.go or the existing hand-commit tests' expectations

## statement

A test shows that the commit rhiza land --commit makes as a conductor session's actor, on its run's branch, is the run's: handCommits journals no hand-commit intervention for it

## Links

### Out

- depends-on [land-commit](../../change-sets/change-lifecycle/land-commit.md)
- decided-by [conductor-land-commit-trailers](../../../../decision/conductor-land-commit-trailers.md)

### In

- verifies [conductor-land-commit-not-hand-commit-verified](../../../../evidence/conductor-land-commit-not-hand-commit-verified.md)
- verifies [conductor-land-commit-not-hand-commit-qa-1](../../../../evidence/conductor-land-commit-not-hand-commit-qa-1.md)

## History

- `01M7RJM81GCCYXK77G7HMKA0FR` Fix False hand-commit intervention: run 01M7RHR8NEJTECVZ5FYFDJ0H18 (change set 01M7RHC6VMMBXVWS3BKKFF6F8N) journals 'Interventions: 1 (hand-commit 1)' for commit 351cf62f6 'Land and render the run-status flake fix', which its own conductor (session 9239a99d, actor claude-code-9239a99d) made after rhiza land and rhiza render. The commit has no Change:/Actor:/Leaf:/Review: trailer, so handCommits (internal/api/run_handcommit.go) counts it as a hand commit and the daemon's clean streak breaks with no human involved. Earlier conductors happened to add trailers; nothing makes them. Make the conductor's own land/render commit always recognised (e.g. rhiza land/render commit with the run's trailers, or the conduct guide and a check requiring Change: and Actor: trailers on every conductor commit), and add a test where a conductor's land commit is not journalled as an intervention. This blocks milestone 40's measure (10 consecutive clean runs).
