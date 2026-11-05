# abandon-unpushed-branch-scope

decision · `01M6JX11AVXF2645G0FEY2CMF4`

## reasoning

internal/api/run.go:1070 already deletes with `git branch -d` (PR #164), and run-abandon-unmerged.txtar proves the no-upstream case keeps the branch. git branch -d checks merge against the upstream when one is set, so an ahead-of-upstream branch should also be kept, but nothing tests that, and the model has no leaf recording the behavior. The change set's intent (keep unpushed commits instead of force-deleting) is therefore best served by recording the capability as a buildable leaf and pinning the untested upstream case with a script, not by rewriting working code. Runs build only childless leaves, so the work needs that leaf to ship under milestone-25.

## statement

Change set 01M6JWSC builds no new deletion logic: it adds one childless leaf under run-abandon saying rhiza run abandon keeps a local branch holding commits merged into neither its upstream nor HEAD, and covers it with a test for the missing case (a branch with an upstream plus local commits never pushed to it) beside the existing run-abandon-unmerged.txtar no-upstream case; the code changes only if that new test shows the branch is deleted

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-b947d57e-1](../question/question-b947d57e-1.md)

### In

- decided-by [abandon-keeps-unpushed-branch](../capability/rhiza/workflows/workflow-runs/run-abandon/abandon-keeps-unpushed-branch.md)

## History

- `01M6JWSCZYYQVFR6ZTVWS6JWJ0` Fix rhiza run abandon keeps a local branch holding commits that were never pushed, instead of force-deleting them (replaces draft change set 01M6B9FCNTTH65, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
