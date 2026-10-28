# ship-needs-leaves-verified-3

evidence · `01M60EBWX66QB0AJSGSZWQTPR9` · actor [claude-code-9c5f8a51](../actor/claude-code-9c5f8a51.md)

## result

PASS

## statement

Verifier claude-code-9c5f8a51, on the code at 736c045: `go test ./cmd/rhiza -run 'TestScripts/run-ship' -count=1` passed. In it, run-ship-needs-leaves.txtar runs `rhiza run start --workflow ship-change` on a change set with no leaves; the command exits non-zero, stderr says 'has no leaves' and names `rhiza run start --workflow plan`, and the --json form carries "next": "rhiza run start --workflow plan". Afterwards .git/rhiza/runs does not exist and no worktree is listed. The same command on a change set with a capability prints 'Run ... succeeded'. In the diff, internal/api/run.go:404-408 calls checkShipLeaves before newRunID, canRun, t.open and run.Start, so the refusal writes nothing. The acceptance quotes `rhiza run start plan` as next. The QA review's confirmed finding f2 ruled that form unrunnable (the positional is the change set) and called for `--workflow plan`, which c2886a0 and 736c045 put in place.

## Links

### Out

- verifies [ship-needs-leaves](../capability/ship-through-rhiza/ship-needs-leaves.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
