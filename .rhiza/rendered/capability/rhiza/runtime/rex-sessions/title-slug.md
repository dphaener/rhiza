# title-slug

capability · `01M5TYCXMES38V9E8DBVEZMFV4` · parent [rex-sessions](../rex-sessions.md)

## acceptance

A test derives the slug, branch and worktree names from one title.

## statement

The branch `cs/<slug>` and worktree `<repo>.<slug>` are named from the same title.

## Links

### Out

- decided-by [one-slug-everywhere](../../../../decision/one-slug-everywhere.md)

### In

- verifies [title-slug-test](../../../../evidence/title-slug-test.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
