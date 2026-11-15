# evidence-commit-retries-moved-head

decision · `01M7FT14N1P6AH118XPRZPAPQP`

## reasoning

headless-run-agent-tabs failed on #593's first CI attempt: a verifier's Evidence commit was refused ('cannot lock ref HEAD: is at c8e252a but expected e351d72') because another session committed in the same run worktree between commitOnly's HEAD read and its update-ref. The swap is right to refuse a stale parent; rebuilding on the new HEAD keeps both commits. The brief asks for a reproduction by injected delay, not stress, and a fix; a hook in the window is that delay, deterministic.

## statement

commitOnly retries when its compare-and-swap of HEAD is refused because HEAD moved: it reads the new HEAD and builds the commit again on it, at most 5 times, then fails with git's error as today; any other failure fails at once. A package hook between reading HEAD and update-ref lets a test inject the racing commit.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [fast-suite-finish-research](../source/fast-suite-finish-research.md)

### In

- decided-by [evidence-commit-survives-moved-head](../capability/evidence-commit-survives-moved-head.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
