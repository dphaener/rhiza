# cli-child-env-fix-verified

evidence · `01M7HQJYYECKTJBE6DFW440SQM` · actor [claude-code-f0614309](../actor/claude-code-f0614309.md)

## result

PASS

## statement

Verified cli-child-env at fix commit 7d0913867 in a detached scratch worktree, RHIZA_/SHIP_/REX_ unset, every go test under a capped systemd scope with -p 2. Every Checks: line exits 0 as written: go test ./cmd/rhiza -skip TestScripts; go test ./internal/api ./internal/envexec; TestScripts/^(rex-install|run-status-tokens|run-ship-needs-leaves|run-log|run-tail)$; TestScripts/(ship-|session-hook|usage-live); go vet ./cmd/rhiza ./internal/api ./internal/envexec; rhiza ship generate (built from the commit: make docs, 0 changed); rhiza render. Acceptance: ./cmd/rhiza -run '^Test(NewCtx|CtxFromEnv|ShellWithStdin|ChildEnv)' PASS (TestChildEnv, TestChildEnvPath/absolute and /relative); ./internal/api -run '^Test(ShipGate|ShipPrePush|ShipGenerate|DaemonHook|UsageLive|ChildEnv)' PASS (TestChildEnv, TestChildEnvPath); ./internal/index ./internal/run ./internal/gosweep PASS. TestChildEnvPath puts a fake git first on run()'s given PATH (absolute and relative entry) while the process's PATH lacks it, and the fake runs. Failure case: a program only on the process's PATH is ErrNotFound through LookPath, Exec, Start and Detach; the api test gets ErrNotFound for a hook whose Ctx PATH has no sh. Mutation check: making envexec.Command use exec.LookPath fails all of those assertions (reverted). TestNoProcessGlobals/fixtures PASS: it flags exec.LookPath, and exec.Command/CommandContext with a bare, relative or variable name under an aliased import, but not an absolute literal. The allowlist is main, selfPath (main-only), internal/daemon/reexec.go and envexec.Command, each with a reason that holds. A grep of the non-test code finds no exec.Command/LookPath on run()'s path left outside envexec, apart from selfPath. Not-found errors are still *exec.Error wrapping ErrNotFound, so output, flags and exit codes are unchanged. Note, outside this leaf's acceptance (env/PATH, not Dir): shellWithStdin now runs in wd, but runUserStatusLine (internal/api/usage_live.go:104) and runProjectClone's git clone (internal/api/project.go:228) still set no cmd.Dir, so they run in the process's working directory. The CLI is unaffected because its working directory is wd; an in-process caller with another wd would run them elsewhere.

## Links

### Out

- verifies [cli-child-env](../capability/rhiza/integrity/self-check/cli-child-env.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
