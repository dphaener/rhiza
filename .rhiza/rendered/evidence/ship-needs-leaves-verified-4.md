# ship-needs-leaves-verified-4

evidence · `01M61BWHAPV2KHNK33AJ6CJSEN` · actor [claude-code-610bc6c7](../actor/claude-code-610bc6c7.md)

## result

PASS

## statement

At 1929007, go test ./cmd/rhiza -run TestScripts/run-ship-needs-leaves -count=1 passes: on a change set with no leaves `rhiza run start --workflow ship-change` exits non-zero with 'has no leaves' and 'rhiza run start --workflow plan' on stderr (and "next": "rhiza run start --workflow plan" under --json), no .git/rhiza/runs and no worktree for the change set; after adding a capability the same command prints 'Run <id> succeeded' and journals the run. Read internal/api/run.go: checkShipLeaves runs in runRunStart before newRunID, t.open (worktree) and run.Start, so nothing is written on refusal. The 1929007 fix only renames the leaf's statement/acceptance to `run start --workflow plan`, matching the code. make check exits 0.

## Links

### Out

- verifies [ship-needs-leaves](../capability/ship-through-rhiza/ship-needs-leaves.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
