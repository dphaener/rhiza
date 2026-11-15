# parallel-scripts-plan-review

source · `01M7EYZQVMZTCGPSEZ671SETVX`

## quotes

```json
[]
```

## raw_sha256

e5fcccd122cdcb33c25f62f16fab881215ccd255cc06d97d185fe744b770256a

## session

2109b166-ad65-4eaf-80f5-01a1ddec9b1e

## summary

Autonomous plan run (no owner present) completing the four drafted leaves of the Parallel weighted scripts change set after the virtual clock merged as PR #588. The orchestrator's message set the rules: keep script parallelism bounded for the 4-core VM that OOMed on 2026-11-14, take the timing file from a CI scripts job's go test -json output rather than a full local run, do not widen scope, record choices as assumed Decisions. Research found: CI's scripts jobs run go test without -json, so no per-script time reaches any CI log today; the post-clock CI run 43685514513 had script shards at 106, 172, 197 and 151 s (round-robin); testscript keeps Params.Files in the order given; go test -json makes testscript log every script, so a CI log of raw JSON would bury the --- FAIL: lines ship check --test-failures reads; the Makefile has no CPU-count variable a test can set. The plan adds a leaf that makes CI's scripts jobs print one script-time line per passing script from go test -json while keeping go test's failure text, reworks the timings tool to read those CI log lines, seeds the checked-in file and has the orchestrator refresh it from the PR's own CI run before merge, and makes the CPU count injectable (NPROC) so make -n tests can check the -parallel arithmetic.

## title

Autonomous plan review of Parallel weighted scripts (2026-11-15)

## Links

### In

- derived-from [script-times-from-ci-log](../decision/script-times-from-ci-log.md)
- derived-from [script-timings-seed-then-ci](../decision/script-timings-seed-then-ci.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
