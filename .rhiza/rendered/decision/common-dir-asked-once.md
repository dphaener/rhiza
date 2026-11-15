# common-dir-asked-once

decision · `01M7FT14HETCZ40CS5HDTBR8KM`

## reasoning

757 of run-ship-conflict-fix-more's 3,874 git children are rev-parse --git-common-dir, about 2.6 per rhiza process (287), and 133 of bugfix-workflow-autonomous's 481. A checkout's common dir does not change during one process, so keeping it is safe. orphanLogEntries (check.go) asks for the same checkout's common dir with --path-format=absolute on every rhiza check and model commit-untracked and compares it with rev-parse --absolute-git-dir, so it takes the kept answer too. The kept answer is asked with --path-format=absolute rather than made absolute with filepath.Abs(ctx.Root) as gitCommonDir does today: filepath.Abs keeps a symlinked root's path (macOS /var -> /private/var, a symlinked checkout) while git resolves it, so a kept Abs-joined path would no longer equal --absolute-git-dir in a main checkout and orphanLogEntries would silently report nothing. Sessions and countSessionFiles' glob name the same folder either way. Only run_leaves.go's openCheckout, which asks in the build worktree through worktreeGit, answers for another directory and stays as it is.

## statement

A rhiza process asks git for a checkout's common dir once, as git rev-parse --path-format=absolute --git-common-dir (git's own absolute, symlink-resolved path, the form rev-parse --absolute-git-dir prints), and keeps the cleaned answer on its Ctx; cmd/rhiza's test-mode ID offset (countSessionFiles), internal/api's gitCommonDir and orphanLogEntries in internal/api/check.go share it

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [fast-suite-finish-research](../source/fast-suite-finish-research.md)

### In

- decided-by [common-dir-once](../capability/common-dir-once.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
