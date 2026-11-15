# m38-targets-in-reach

decision · `01M7FT15889DZJVS7YHJJ5FN82`

## reasoning

Cold make test used 1,089 CPU-s at 258% of 4 cores (422 s), a 272 s floor at full use, so waits alone cannot reach 240 s. The leaves cut CPU and the long poles: about 1,550 of every 3,900 git children in a busy script (settings and common-dir reads), TestLinkEndpointKinds' 81 s, 23 s of serial budgets and the four scripts' dead waits (about 280 script-seconds and the 99 s tail). That should bring CPU near 850 CPU-s and wall near 220-240 s. CI: with no recompile (about 50 s a job) and balanced shards the longest job should be about 110-130 s, a run about 150 s. If the orchestrator's measurements after landing still miss, the numbers go to the owner then.

## statement

m38-cold and m38-ci stay as written; this change set's leaves are expected to meet them, so no question goes to the owner about the hardware

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [fast-suite-finish-research](../source/fast-suite-finish-research.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
