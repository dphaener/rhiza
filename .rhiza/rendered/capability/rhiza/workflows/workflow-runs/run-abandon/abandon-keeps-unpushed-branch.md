# abandon-keeps-unpushed-branch

capability · `01M6JX3902A2MVZQSQQJVYCGBZ` · parent [run-abandon](../run-abandon.md)

## acceptance

Beside run-abandon-unmerged.txtar (no upstream, branch kept), a testscript pushes cs/<slug> to a bare origin with an upstream set, adds a local commit never pushed, abandons the run and asserts 'delete-branch kept: it has commits merged into neither its upstream nor HEAD' and that git log cs/<slug> still shows the commit

## statement

rhiza run abandon keeps the run's local cs/<slug> branch when it holds commits merged into neither its upstream nor HEAD, reporting delete-branch kept with that reason, whether or not the branch has an upstream

## Links

### Out

- decided-by [abandon-unpushed-branch-scope](../../../../../decision/abandon-unpushed-branch-scope.md)
- derived-from [abandon-unpushed-plan](../../../../../source/abandon-unpushed-plan.md)

### In

- verifies [abandon-keeps-unpushed-branch-verified](../../../../../evidence/abandon-keeps-unpushed-branch-verified.md)
- verifies [abandon-keeps-unpushed-branch-qa-1](../../../../../evidence/abandon-keeps-unpushed-branch-qa-1.md)
- verifies [pr-267-fix-rhiza-run-abandon-keeps-a](../../../../../evidence/pr-267-fix-rhiza-run-abandon-keeps-a.md)

## History

- `01M6JWSCZYYQVFR6ZTVWS6JWJ0` Fix rhiza run abandon keeps a local branch holding commits that were never pushed, instead of force-deleting them (replaces draft change set 01M6B9FCNTTH65, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
