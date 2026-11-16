# run-clock-from-ctx-verified

evidence · `01M7HGNMJFNM8EF58S664NTK57` · actor [claude-code-7df4e6a2](../actor/claude-code-7df4e6a2.md)

## result

PASS

## statement

Verified at 5ed97ba50 in a detached worktree (RHIZA_/SHIP_/REX_ unset, capped scope, -p 2): go test -trimpath -count=1 ./internal/run -run Clock ok; ./internal/api -run '^Test(RunClock|CtxSleep|RunStart|RunResume|SessionWait)' ok, including TestRunClockPerCtx (two Ctxs from scale 1000 and empty with the process env unset get ScaledClock and RealClock; 'abc','1.5','0' give RealClock; a set runClock overrides both); full ./internal/api ok (224s); ./cmd/rhiza TestScripts/(run-start|run-resume|session-run) 21 scripts pass; go vet ./internal/run ./internal/api ./cmd/rhiza clean; gofmt clean; ship generate and render leave the tree clean. Diff: runClock is a nil package var, clockOf(ctx) = runClock, else ctx.Clock, else RealClock; run.start/resume engines, abandonedCleanup, spawn poll (ctxSleep), retryRex/retryGH, typePointer, the usage-limit pause and ship deploy all use clockOf(ctx); newCtx (the only Ctx constructor) sets Clock from run.ClockFromEnv(getenv). No durations, output, flags or exit codes changed.

## Links

### Out

- verifies [run-clock-from-ctx](../capability/rhiza/integrity/self-check/run-clock-from-ctx.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
