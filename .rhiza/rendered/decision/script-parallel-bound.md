# script-parallel-bound

decision · `01M7EQE6332J8WHCG6GRHJXZA1`

## reasoning

The measured -parallel 12 ran with 0 failures, but the VM OOM-crashed on 2026-11-14 when parallel suites spawned about 1,600 script children, and the owner has 4 cores and at most 24 GB to give. Once scripts stop sleeping (the scaled clock) they are CPU-bound, so more than twice the cores adds memory and children, not speed; 8 keeps one suite well inside a 6 GB, 300-task scope.

## statement

cmd/rhiza's scripts run at -parallel of twice the CPU count, at most 8, unless RHIZA_TEST_SCRIPT_PARALLEL names another positive number; on the 4-core VM that is 8

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [perf-milestone-interview](../source/perf-milestone-interview.md)

### In

- decided-by [scripts-parallel-bounded](../capability/rhiza/integrity/self-check/scripts-parallel-bounded.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
