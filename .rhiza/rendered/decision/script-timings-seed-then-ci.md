# script-timings-seed-then-ci

decision · `01M7EZ13EVPEKPW95GAYJB3C0H`

## reasoning

No CI log before this change set carries per-script times, the orchestrator ruled out a full local run as the source, and builders never push, so CI times exist only once the PR's CI has run. Equal seconds make the order and the shards fall back to name order (round-robin), today's behaviour, until the refresh commit lands with the change set.

## statement

The first checked-in script-timings.tsv gives every script 1 second, and the orchestrator rewrites it with make script-timings RUN=<run> from the change set's own PR CI run before merge

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [parallel-scripts-plan-review](../source/parallel-scripts-plan-review.md)

### In

- decided-by [script-timings-file](../capability/rhiza/integrity/self-check/script-timings-file.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
