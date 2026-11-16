# run-concurrent-guard

capability · `01M7H0YQN1QRDYASREWWQD9ED9` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 ./cmd/rhiza -run '^Test(RunConcurrentDirs|NoProcessGlobals)$' passes and go test -race -trimpath -count=1 ./cmd/rhiza -run '^TestRunConcurrentDirs$' is clean: TestRunConcurrentDirs calls run() from two goroutines, many times each, on two scratch repos with different working directories, RHIZA_ACTOR values and HOMEs, and shows each call's output names its own repo's nodes and actor. Failure case: TestNoProcessGlobals fails, naming file and line, on a fixture file that calls os.Getenv outside the allowlist, and on one that calls os.UserCacheDir. go vet ./cmd/rhiza is clean. Touches cmd/rhiza's new test file(s) only, plus whatever site the guard finds that the earlier leaves left (fixed, not allowlisted, unless it is the daemon re-exec).

## statement

Two run() calls proceed at once in one process, each answering for its own working directory and environment, and a test in cmd/rhiza keeps it so: it parses cmd/rhiza's and internal/api's, internal/run's, internal/daemon's, internal/gosweep's and internal/index's non-test Go files and fails on any use of os.Getwd, os.Chdir, os.Args, os.Executable, os.Getenv, os.LookupEnv, os.Environ, os.Setenv, os.Unsetenv, os.UserCacheDir, os.UserHomeDir or os.TempDir outside an allowlist in that test that names each remaining site (main() and selfPath, internal/daemon/reexec.go) with its reason

## Links

### Out

- depends-on [run-clock-from-ctx](run-clock-from-ctx.md)
- decided-by [run-takes-env-list](../../../../decision/run-takes-env-list.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)

### In

- depends-on [scripts-in-process-harness](scripts-in-process-harness.md)
- verifies [run-concurrent-guard-verified](../../../../evidence/run-concurrent-guard-verified.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
