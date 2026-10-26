# tab-title

capability · `01M5TYCX7CHJ1H823GCD762206` · parent [rex-sessions](../rex-sessions.md)

## acceptance

A test builds the title for a change set with and without a title and for a build and a review step.

## statement

A session's Rex tab is named `<repo>: <change set title>`, for example `genesis: Carving keeps IDs`, with no state in it and the step added after ` · ` only when it is not the build step.

## Links

### Out

- decided-by [tab-title-repo-prefix](../../../../decision/tab-title-repo-prefix.md)
- decided-by [no-state-in-tab-title](../../../../decision/no-state-in-tab-title.md)
- decided-by [human-readable-tab-names](../../../../decision/human-readable-tab-names.md)

### In

- verifies [tab-title-test](../../../../evidence/tab-title-test.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
