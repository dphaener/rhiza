# run-clock-chosen-per-ctx

decision · `01M7H0YQ9WASMPHSJT2XGW6S0Z`

## question

How does a run's clock follow the script's RHIZA_TEST_CLOCK_SCALE once rhiza runs in-process?

## reasoning

runClock is chosen once at package init from the process environment (internal/api/run.go:100). In-process, the test binary has no RHIZA_TEST_CLOCK_SCALE, so every script's 30 s polls and backoffs would run on the real clock, and a script that clears the scale with `env RHIZA_TEST_CLOCK_SCALE=` would not get the real clock either way. 35 test sites replace runClock, so it stays as an override rather than being removed.

## statement

A run's waits go by the clock chosen from the Ctx's own environment (RHIZA_TEST_CLOCK_SCALE through ctx.Getenv, via run.ClockFromEnv taking a getenv func) when the Ctx is made, held on api.Ctx; the package variable runClock stays only as the tests' override and is nil by default.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [in-process-scripts-plan-research](../source/in-process-scripts-plan-research.md)

### In

- decided-by [run-clock-from-ctx](../capability/rhiza/integrity/self-check/run-clock-from-ctx.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
