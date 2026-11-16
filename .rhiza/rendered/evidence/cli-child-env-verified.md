# cli-child-env-verified

evidence · `01M7HFQKF0Z30EGNSYCZBDHJSH` · actor [claude-code-6a12bf26](../actor/claude-code-6a12bf26.md)

## result

PASS

## statement

At 02752b3f1, capped and with RHIZA_/SHIP_/REX_ unset, every Checks: command exits 0: go test ./cmd/rhiza -run '^Test(NewCtx|CtxFromEnv|ShellWithStdin|ChildEnv|RunWorkingDir|ParseGlobals)' (TestChildEnv PASS), go test ./internal/api -run '^Test(ShipGate|ShipPrePush|ShipGenerate|DaemonHook|UsageLive|ChildEnv)', go test ./internal/index ./internal/run ./internal/gosweep ./internal/userdir, TestScripts/(ship-|session-hook|usage-live) (35s), internal/api -skip Budget, cmd/rhiza -skip 'TestScripts|Budget', TestScripts/(workflow-|run-start|tui-|daemon-detach), go vet on all six packages, and rhiza ship generate plus rhiza render in a detached worktree, which left it clean. No fork runaway. TestChildEnv t.Setenvs the process HOME, XDG_CACHE_HOME, GIT_CONFIG_GLOBAL and a process-only marker, then calls run() with other values. It shows git reads the given GIT_CONFIG_GLOBAL, the index .db lands under the given XDG_CACHE_HOME/rhiza/index and run-tmp under the given XDG_CACHE_HOME/rhiza/run-tmp, nothing goes under the process cache, a workflow step sees the given marker, and the process-only variable stays unset (the failure case). internal/api/child_env_test.go shows a hook child and the gate and head envs see a Ctx Environ marker the process lacks. Diff read: gitIn, gitConfig, execIn, shellWithStdin, startProgram/programEnv, detach (and its RHIZA_TEST_DETACH_LOG read), tui colorprofile.Detect, and api's ship_build headChildEnv, usage_live, daemon_hook, ship_prepush, ship_generate, ship_gate, run_leaves commitOnHead and project clone all take env from the Ctx (childEnv). index.DefaultPath, Remove and IndexDir, run.TmpDir and Store.Getenv (runStore), personalWorkflowsDir and gosweep.CacheDirFrom place dirs through the new internal/userdir over getenv. Grep finds no os.Getenv/LookupEnv/Environ/UserCacheDir/UserHomeDir/TempDir and no exec without Env left on run()'s path in these packages. The only exceptions are main()'s own os.Environ/RHIZA_CPUPROFILE, nil-fallbacks for programmatic callers, run/clock.go's scale env (left to the dependent leaf run-clock-from-ctx) and daemon reexec. Bare run.NewStore sites only List. No txtar changes, so output, flags and exit codes are unchanged.

## Links

### Out

- verifies [cli-child-env](../capability/rhiza/integrity/self-check/cli-child-env.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
