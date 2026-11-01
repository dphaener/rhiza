# ship-check-reruns-failed-once-qa-1

evidence · `01M6ANV2XXKPCGD29F1CB5WA5Y` · actor [claude-code-709ac66c](../actor/claude-code-709ac66c.md)

## result

PASS

## statement

At 4c9b8c5b built rhiza (go build ./cmd/rhiza) and ran it by hand in a scratch repo with a gh stub shaped like real gh output (read-only gh pr view 153 --json statusCheckRollup on this repo gives detailsUrl .../actions/runs/41270718492/job/...; real gh run list --json databaseId,... works; real gh run rerun --help lists --failed). ship check --settled on two failed checks of one Actions run: one 'run rerun 41270718492 --failed' (deduped), --json reran=[41270718492], settled false, exit 1, ship record rerun_head/reruns set. Next check pending: not settled. Failed again on the same head: --settled reports failed (exit 0) and plain ship check exits 1, still one rerun in gh.log. A new head reruns once more; a rerun then passing gives '3 checks passed' exit 0. A gh rerun error exits 1 with the gh message and the next check fails without a second rerun. --main reads databaseId 99, reruns once, second failure fails, success passes. A failing StatusContext and a CheckRun with a non-Actions URL fail with no rerun.

## Links

### Out

- verifies [ship-check-reruns-failed-once](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-reruns-failed-once.md)

## History

- `01M6A0JQ2WHGH7DDFKJGRXJD85` A failed CI check on a run's PR or main is rerun once per head before the run fails
