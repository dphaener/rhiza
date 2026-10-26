# change-set-title

capability · `01M5TYCZA314H413P69D60GNVC` · parent [change-sets](../change-sets.md)

## acceptance

A testscript sets a title, shows it in `change show --json`, refuses an over-long one, and shows the fallback.

## statement

A change set has an optional short title, set with `genesis change new --title` (under about 30 characters), and without one its intent's first clause stands in.

## Links

### Out

- decided-by [change-set-title-field](../../../decision/change-set-title-field.md)

### In

- verifies [change-set-title-test](../../../evidence/change-set-title-test.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
