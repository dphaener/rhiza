# run-concurrent-guard-verified

evidence · `01M7HHSB6NYRFW62SVXW37WHGF` · actor [claude-code-a184fc2b](../actor/claude-code-a184fc2b.md)

## result

PASS

## statement

In a detached checkout at 078d0c629, capped and with RHIZA_/SHIP_/REX_ unset: go test -p 2 -trimpath -count=1 ./cmd/rhiza -run '^Test(RunConcurrentDirs|NoProcessGlobals)$' exit 0 (both PASS, fixtures subtest PASS: it reports cmd/rhiza/main.go:8 os.Getenv in a non-main function and internal/index/cache.go:5 os.UserCacheDir via an aliased import, ignoring _test.go and the allowlisted reexec.go); -race on '^TestRunConcurrentDirs$' alone exit 0 (two goroutines, 25 rounds, two repos with own dirs, HOMEs and RHIZA_ACTOR; each show names its own actor and node, never the other's); go vet ./cmd/rhiza ./internal/gosweep exit 0; go test ./internal/gosweep/... ./internal/daemon/... exit 0; go test ./internal/api/... (full) exit 0; the internal/api ChildEnv|ResolveWorkflow|PersonalWorkflows|DaemonAutoUpdate run exit 0; TestScripts/^(daemon-settings|workflow-) exit 0; rhiza ship generate and rhiza render exit 0 with nothing changed. Diff read: allowlist is only cmd/rhiza/main.go main, ctx.go selfPath (called only from main) and internal/daemon/reexec.go, each with a reason; fixed sites are childEnv (nil Environ is empty, not os.Environ), personalWorkflowsDir (nil getenv is no dir), the daemon's channel auto-update (ctx.Self resolved via EvalSymlinks instead of os.Executable) and gosweep's unused CacheDir/Sweep removed. The only production api.Ctx is built in cmd/rhiza/ctx.go with Environ and Self set, so CLI output, flags and exit codes are unchanged. No stubs or TODOs.

## Links

### Out

- verifies [run-concurrent-guard](../capability/rhiza/integrity/self-check/run-concurrent-guard.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
