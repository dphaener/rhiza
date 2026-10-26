# ship-state-file

decision · `01M5WHKE4GAD2E52F12RC7MAK3`

## statement

Ship-change's steps pass the PR, the checked head and the merge commit through a state file per change set under the git common dir, kept out of the log.

## wording

the building agent's design, assumed in autonomous mode

## Links

### In

- decided-by [pinned-merge](../capability/rhiza/orchestration/shipping/pinned-merge.md)

## History

- `01M5SZJ03NQS0A389JF8TTMF6M` The ship-change workflow: shipit as the agent step, then wait for a green PR, a merge pinned to the checked head, the close-out (branch, worktree, CI on main, the repo's ship command), and land with the PR and CI recorded as Evidence.
