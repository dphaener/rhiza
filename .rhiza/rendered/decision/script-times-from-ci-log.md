# script-times-from-ci-log

decision · `01M7EZ13D3AQHEA3J3WV9MGQVA`

## reasoning

The orchestrator ruled that the timing file comes from a CI scripts job's go test -json output, not a full local run. CI's scripts jobs run go test without -json today, so no per-script time reaches any log. Raw -json in the log would make testscript log every script and bury the --- FAIL: lines ship check --test-failures and ci-fix read; an uploaded artifact needs a new action. A replay keeps the failure text as go test prints it and adds about 400 short lines per run.

## statement

CI's scripts jobs record each script's time in their own log: make test-scripts runs go test -json through a replay that prints go test's text for what did not pass and one script-time line per passing script, and make script-timings reads those lines with gh run view --log; no artifact upload and no new workflow step

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [parallel-scripts-plan-review](../source/parallel-scripts-plan-review.md)

### In

- decided-by [ci-script-timings](../capability/rhiza/integrity/self-check/ci-script-timings.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
