# repo-ship-command

capability · `01M5TYKQB1VWPAN3DBN01K0BB0` · parent [shipping](../shipping.md)

## acceptance

A testscript's repo ship command runs once, after CI on main is green.

## statement

Close-out ends with the repo's own ship command from .rhiza/workflows/ (make install for Rhiza).

## Links

### Out

- decided-by [task-close-out](../../../../decision/task-close-out.md)
- decided-by [ship-command-file](../../../../decision/ship-command-file.md)

### In

- verifies [repo-ship-command-tests](../../../../evidence/repo-ship-command-tests.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
