# abandon-keeps-unpushed-branch-verified

evidence · `01M6MFT7KCRSTWW0QMYJQZB1WR` · actor [claude-code-b30372f5](../actor/claude-code-b30372f5.md)

## result

PASS

## statement

At 3ebbaa2d, cmd/rhiza/testdata/script/run-abandon-unpushed.txtar sits beside run-abandon-unmerged.txtar. It runs a build that commits on cs/abandon-built, pushes it with -u to a bare origin.git, runs a second build whose commit is never pushed, fails a third run, abandons it, and asserts stdout 'delete-branch kept: it has commits merged into neither its upstream nor HEAD', that git log cs/abandon-built still shows both 'built' commits, and that origin/cs/abandon-built..cs/abandon-built counts 1. go test -trimpath ./cmd/rhiza -run 'TestScripts/run-abandon-(unpushed|unmerged)$' -count=1 -v: both PASS. The message comes from internal/api/run.go's git branch -d guard, which refuses a branch not merged into its upstream or HEAD.

## Links

### Out

- verifies [abandon-keeps-unpushed-branch](../capability/rhiza/workflows/workflow-runs/run-abandon/abandon-keeps-unpushed-branch.md)

## History

- `01M6JWSCZYYQVFR6ZTVWS6JWJ0` Fix rhiza run abandon keeps a local branch holding commits that were never pushed, instead of force-deleting them (replaces draft change set 01M6B9FCNTTH65, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
