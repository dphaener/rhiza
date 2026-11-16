# run-clock-from-ctx-qa-1

evidence · `01M7HZJ7PH5SA3V5S33QT4STDZ` · actor [claude-code-10d5f94b](../actor/claude-code-10d5f94b.md)

## result

PASS

## statement

QA: a wait step (timeout 2s, interval 1s) under run start --foreground with RHIZA_TEST_CLOCK_SCALE=1000, empty and abc: 40, 3 and 3 checks in about 2.05 s each, identical to the base build.

## Links

### Out

- verifies [run-clock-from-ctx](../capability/rhiza/integrity/self-check/run-clock-from-ctx.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
