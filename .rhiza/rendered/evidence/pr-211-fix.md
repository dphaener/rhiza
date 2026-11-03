# pr-211-fix

evidence · `01M6H20MAD2SNGVXS4KEP718GR` · actor [dphaener](../actor/dphaener.md)

## head

5bb528e95a311ea5506d18b84a90577f849254a4

## pr

https://github.com/dphaener/rhiza/pull/211

## result

1 check passed: check SUCCESS

## statement

Pull request #211 (https://github.com/dphaener/rhiza/pull/211) passed CI at 5bb528e

## status

verified

## Links

### Out

- verifies [close-out-drops-log-copies](../capability/rhiza/orchestration/shipping/close-out-drops-log-copies.md)
- verifies [close-out-retry-after-forget](../capability/rhiza/orchestration/shipping/close-out-retry-after-forget.md)

## History

- `01M6GCPSMFW8G6MXXNW6Y5J0Q6` Fix: Request plan runs fail at close-out: the daemon writes the new change set's first .rhiza/log entries untracked into the main checkout, so after the plan PR merges close-out's git pull --ff-only refuses to overwrite them (run 01M6FSQBBJYK). The plan run's log entries should be written only in its worktree, or close-out should drop untracked log files identical to the incoming ones before pulling; and a retried close-out must not fail because the first visit forgot the ship state
