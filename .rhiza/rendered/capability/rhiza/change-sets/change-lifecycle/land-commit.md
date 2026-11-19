# land-commit

capability · `01M7RJTNV0969V311Y3SBK5BMP` · parent [change-lifecycle](../change-lifecycle.md)

## acceptance

go test -trimpath ./cmd/rhiza -run 'TestScripts/(land-commit|describe-golden)$' passes: land-commit.txtar lands an accepted change set with RHIZA_ACTOR set and rhiza land --commit --json, then git log -1 --format=%(trailers) shows Change: <the change set> and Actor: <that actor's handle>, git status --porcelain .rhiza is empty, and a file staged outside .rhiza before the land is still staged and not in the commit; (failure case) rhiza land --commit on a draft change set exits non-zero with HEAD unchanged and nothing committed; (failure case) rhiza land --commit --dry-run is a usage error with nothing landed or committed; describe-golden.txtar is regenerated with the new land request type. And go test -trimpath ./internal/api -run TestLandCommitViaHead passes: with SHIP_HEAD_RHIZA naming a head build that is not the running binary and ctx.Exec a fake head that records its argv and lands and renders in-process, land with commit set runs land and render through the head with no --commit in either argv, and then HEAD is a new commit of only .rhiza with Change: and Actor: trailers; (failure case) when the fake head's land exits non-zero, land with commit set returns an error, HEAD is unchanged and nothing is committed

## files

internal/api/lifecycle.go, internal/api/land_commit.go, internal/api/land_commit_test.go, cmd/rhiza/testdata/script/land-commit.txtar, cmd/rhiza/testdata/script/describe-golden.txtar

## must_not

Change rhiza land without --commit, or the other lifecycle verbs' requests; forward --commit to the head build or change the arguments runViaHead forwards (internal/api/ship_build.go); commit anything outside .rhiza; push; change ship land's own commit (internal/api/ship_close.go); change handCommits' attribution rule

## statement

rhiza land --commit lands the change set, renders, stages .rhiza and commits only .rhiza with the message 'Land and render the change set' and the trailers 'Change: <change set ID>' and 'Actor: <actor handle>' (the actor the land entry was written as, resolved as every write resolves it: RHIZA_ACTOR, else git's user.email; its ID when it has no handle), all in one trailer block; without --commit land behaves as today. When SHIP_HEAD_RHIZA names a head build other than the running binary (as bugfix and ship-rhiza set in [env]), the land and the render still run through that head (runViaHead, its forwarded arguments unchanged: --commit is never forwarded), and the running rhiza makes the commit in-process after the head's land and render return, so a head built without --commit lands the same way; a failed head land or render commits nothing

## Links

### Out

- decided-by [conductor-land-commit-trailers](../../../../decision/conductor-land-commit-trailers.md)

### In

- depends-on [conductor-land-commit-not-hand-commit](../../workflows/workflow-runs/conductor-land-commit-not-hand-commit.md)
- depends-on [conduct-guide-land-commit](../../guide/agent-guide/conduct-guide-land-commit.md)
- verifies [land-commit-verified](../../../../evidence/land-commit-verified.md)
- verifies [land-commit-qa-1](../../../../evidence/land-commit-qa-1.md)

## History

- `01M7RJM81GCCYXK77G7HMKA0FR` Fix False hand-commit intervention: run 01M7RHR8NEJTECVZ5FYFDJ0H18 (change set 01M7RHC6VMMBXVWS3BKKFF6F8N) journals 'Interventions: 1 (hand-commit 1)' for commit 351cf62f6 'Land and render the run-status flake fix', which its own conductor (session 9239a99d, actor claude-code-9239a99d) made after rhiza land and rhiza render. The commit has no Change:/Actor:/Leaf:/Review: trailer, so handCommits (internal/api/run_handcommit.go) counts it as a hand commit and the daemon's clean streak breaks with no human involved. Earlier conductors happened to add trailers; nothing makes them. Make the conductor's own land/render commit always recognised (e.g. rhiza land/render commit with the run's trailers, or the conduct guide and a check requiring Change: and Actor: trailers on every conductor commit), and add a test where a conductor's land commit is not journalled as an intervention. This blocks milestone 40's measure (10 consecutive clean runs).
