# conductor-land-commit-not-hand-commit-verified

evidence · `01M7RKVRFG4XVG1E2H5AG9EFAK` · actor [claude-code-5ca1b327](../actor/claude-code-5ca1b327.md)

## result

PASS

## statement

Detached checkout at e4a8b7d75: go test -trimpath -p 2 ./internal/api -run TestHandCommits (systemd-run scope MemoryMax=6G TasksMax=300) exit 0; -run TestHandCommitsConductorLandCommit -count=1 -v passes both subtests (a session of the run: land {commit:true} makes a commit with Actor: claude-code-s1, scanHand finds none; a session of another run: one hand-commit naming claude-code-s2); go vet ./internal/api, rhiza ship generate, rhiza render exit 0. Diff touches only run_handcommit_test.go (new test appended) plus the builder's actor record; run_handcommit.go and existing tests unchanged; no stubs or TODOs.

## Links

### Out

- verifies [conductor-land-commit-not-hand-commit](../capability/rhiza/workflows/workflow-runs/conductor-land-commit-not-hand-commit.md)

## History

- `01M7RJM81GCCYXK77G7HMKA0FR` Fix False hand-commit intervention: run 01M7RHR8NEJTECVZ5FYFDJ0H18 (change set 01M7RHC6VMMBXVWS3BKKFF6F8N) journals 'Interventions: 1 (hand-commit 1)' for commit 351cf62f6 'Land and render the run-status flake fix', which its own conductor (session 9239a99d, actor claude-code-9239a99d) made after rhiza land and rhiza render. The commit has no Change:/Actor:/Leaf:/Review: trailer, so handCommits (internal/api/run_handcommit.go) counts it as a hand commit and the daemon's clean streak breaks with no human involved. Earlier conductors happened to add trailers; nothing makes them. Make the conductor's own land/render commit always recognised (e.g. rhiza land/render commit with the run's trailers, or the conduct guide and a check requiring Change: and Actor: trailers on every conductor commit), and add a test where a conductor's land commit is not journalled as an intervention. This blocks milestone 40's measure (10 consecutive clean runs).
