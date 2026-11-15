# contention-packages-left-alone

decision · `01M7FT156PZZB0BEHQTHYEPS45`

## reasoning

Run one test at a time (go test -p 1 -parallel 1 -count=1, budgets skipped) on the idle VM: index 16 s (17 CPU-s), gen 5 s, web 8 s, slice 6 s, against 77, 30, 25 and 30 s inside the 2026-11-15 cold make test. They wait for CPU that cmd/rhiza's scripts and internal/api hold; cutting those is what speeds them up.

## statement

internal/index, internal/gen, internal/web and internal/slice get no work here: their slower cold times in make test are contention, not cost

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [fast-suite-finish-research](../source/fast-suite-finish-research.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
