# land-commit-verified

evidence · `01M7RKRE82F8JPPN9YDW8REA82` · actor [claude-code-053f011d](../actor/claude-code-053f011d.md)

## result

PASS

## statement

In a detached checkout at 06d9a1997 (under systemd-run MemoryMax=6G TasksMax=300, -p 2): go test -trimpath ./cmd/rhiza -run 'TestScripts/(land-commit|describe-golden)$' exit 0 (land-commit.txtar covers Change:/Actor: trailers, clean .rhiza, staged code.go kept out, draft land --commit conflict with HEAD unchanged, --dry-run usage error with state still accepted); go test -trimpath ./internal/api -run TestLandCommitViaHead exit 0 (both subtests: head argv land then render without --commit, one new .rhiza-only commit with trailers; failed head land leaves HEAD unchanged). Re-ran every Checks: line: internal/api 'Lifecycle|Land|Ship|Reference|Docs|Schema|Describe' exit 0, cmd/rhiza land-commit|describe-golden|lifecycle|docs.*|help.* exit 0, rhiza ship generate exit 0 leaving docs/reference/cli.md unchanged (so it is generated output), rhiza render exit 0. Diff read: only land gains LandRequest; runViaHead args, ship_build.go, ship_close.go and handCommits untouched; commit uses pathspec -- .rhiza; no stubs or TODOs.

## Links

### Out

- verifies [land-commit](../capability/rhiza/change-sets/change-lifecycle/land-commit.md)

## History

- `01M7RJM81GCCYXK77G7HMKA0FR` Fix False hand-commit intervention: run 01M7RHR8NEJTECVZ5FYFDJ0H18 (change set 01M7RHC6VMMBXVWS3BKKFF6F8N) journals 'Interventions: 1 (hand-commit 1)' for commit 351cf62f6 'Land and render the run-status flake fix', which its own conductor (session 9239a99d, actor claude-code-9239a99d) made after rhiza land and rhiza render. The commit has no Change:/Actor:/Leaf:/Review: trailer, so handCommits (internal/api/run_handcommit.go) counts it as a hand commit and the daemon's clean streak breaks with no human involved. Earlier conductors happened to add trailers; nothing makes them. Make the conductor's own land/render commit always recognised (e.g. rhiza land/render commit with the run's trailers, or the conduct guide and a check requiring Change: and Actor: trailers on every conductor commit), and add a test where a conductor's land commit is not journalled as an intervention. This blocks milestone 40's measure (10 consecutive clean runs).
