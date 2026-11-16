# run-clock-from-ctx

capability · `01M7H0YQHE7WQMF5VZJDWBVTGM` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 ./internal/run -run Clock and go test -trimpath -count=1 ./internal/api -run '^Test(RunClock|CtxSleep|RunStart|RunResume|SessionWait)' pass, with a new test in internal/api that two Ctxs made with different RHIZA_TEST_CLOCK_SCALE values (1000 and empty) in one process get a scaled and a real clock while the process environment has neither, and that a test which sets runClock still overrides both; failure case: an unparsable scale gives the real clock, as ClockFromEnv does today. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-start|run-resume|session-run)' passes; go vet ./internal/run ./internal/api is clean. Touches internal/run/clock.go, internal/api/run.go, api.go and the files that read runClock, cmd/rhiza/ctx.go, and their tests. Must not change a wait's length under either clock.

## statement

A run's waits (run.start and run.resume polls, the spawn poll, session waits, the default ctx.Sleep and ship deploy) go by a clock chosen from the Ctx's environment when the Ctx is made: run.ClockFromEnv takes a getenv func, api.Ctx holds the chosen Clock, and the package variable runClock (internal/api/run.go:100) is nil unless a test sets it as an override

## Links

### Out

- depends-on [cli-child-env](cli-child-env.md)
- decided-by [run-clock-chosen-per-ctx](../../../../decision/run-clock-chosen-per-ctx.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)

### In

- depends-on [run-concurrent-guard](run-concurrent-guard.md)
- verifies [run-clock-from-ctx-verified](../../../../evidence/run-clock-from-ctx-verified.md)
- verifies [run-clock-from-ctx-qa-1](../../../../evidence/run-clock-from-ctx-qa-1.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
