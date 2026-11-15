# script-timings-from-ci-then-local

decision · `01M7FT145REH8RPZ842V223TZ2`

## reasoning

The file still holds the 1-second seed for all 398 scripts (decision script-timings-seed-then-ci left the CI refresh for later and nobody ran it), so #593's weighted shards weigh nothing: scripts 3/4 ran 112 s of tests to the others' 60-80 s. Run 43721082089 is the newest code run on main's tree with every script timed. Its four slow scripts' times go stale once their dead wait is gone, and a full local scripts run is not allowed on this VM, so only those four are re-timed, alone.

## statement

cmd/rhiza/testdata/script-timings.tsv is filled from CI run 43721082089's script-time lines (make script-timings RUN=43721082089), then the four scripts fake-rex-no-dead-wait fixed are timed again locally, one go test -json run of those four alone merged with make script-timings JSON=<file>

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [fast-suite-finish-research](../source/fast-suite-finish-research.md)

### In

- decided-by [script-timings-measured](../capability/script-timings-measured.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
