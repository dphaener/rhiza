# ship-needs-leaves-verified

evidence · `01M6097AWRRRXHTGNSX4983MTR` · actor [claude-code-1d8bb0bf](../actor/claude-code-1d8bb0bf.md)

## result

PASS

## statement

At 91e157b I ran `go test ./cmd/rhiza -run TestScripts/run-ship-needs-leaves -count=1 -v`: on a change set with no leaves, `rhiza run start --workflow ship-change` exited 3 with 'change set ... has no leaves: ship-change builds and reviews its leaves, so plan it first' and 'next: rhiza run start plan' (the --json form carried "next": "rhiza run start plan"), .git/rhiza/runs did not exist and `git worktree list` showed only the main checkout; after adding capability `thing` the same command printed 'Run <id> succeeded' and journaled the run. The diff puts checkShipLeaves in runRunStart after the change set resolves and before newRunID, the journal or the worktree are made. `go test ./...`, go vet and gofmt passed; make self-check failed only on rendered files for this verifier session's own actor, which is not in the leaf commit.

## Links

### Out

- verifies [ship-needs-leaves](../capability/ship-through-rhiza/ship-needs-leaves.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
