# settings-one-config-read-per-call

decision · `01M7FT14DKS9EMMANF3SYAC2YY`

## reasoning

In run-ship-conflict-fix-more, 1,650 of 3,874 git children are git config --get, about 40 per settingsOf call (one per setting), plus four per autonomy read; bugfix-workflow-autonomous shows the same share (227 of 481). Scripts are CPU-bound once the dead waits are gone (13 s alone for 15 CPU-s) and cold make test needs CPU cut (1,089 CPU-s, a 272 s floor on 4 cores). A cache across calls would go stale: the daemon re-reads settings every poll and daemon set writes them while processes run.

## statement

settingsOf reads git config once per call (one git config -z --get-regexp of the rhiza keys, or --list) and looks each setting up in that snapshot, and the autonomy settings are read the same way; there is no cache across calls. A snapshot keeps git config --get's meaning: the last value of a key wins, section and key names match case-insensitively, a subsection case-sensitively.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [fast-suite-finish-research](../source/fast-suite-finish-research.md)

### In

- decided-by [settings-one-git-read](../capability/settings-one-git-read.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
