# cli-child-env

capability · `01M7H0YQDRC85SWH2704XH2WC9` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 ./cmd/rhiza -run '^Test(NewCtx|CtxFromEnv|ShellWithStdin|ChildEnv)' and go test -trimpath -count=1 ./internal/api -run '^Test(ShipGate|ShipPrePush|ShipGenerate|DaemonHook|UsageLive|ChildEnv)' pass, with a new test TestChildEnv in cmd/rhiza that calls run() with an environment whose HOME and GIT_CONFIG_GLOBAL differ from the process's (set with t.Setenv to other values) and shows the git child sees the given ones (a git config value only that GIT_CONFIG_GLOBAL holds is read back), and a test in internal/api that a gate or hook child started through a Ctx whose Environ holds a marker variable the process lacks sees it; go test -trimpath -count=1 ./internal/index ./internal/run ./internal/gosweep passes, and TestChildEnv also shows that run() given a HOME and XDG_CACHE_HOME other than the process's puts the index database under the given XDG_CACHE_HOME/rhiza/index and the run's temp folder under the given XDG_CACHE_HOME/rhiza/run-tmp, never under the process's cache dir; failure case: a Ctx whose Environ unsets a variable the process has (PATH-free, or a removed marker) leaves it unset in the child. go test -trimpath ./cmd/rhiza -run 'TestScripts/(ship-|session-hook|usage-live)' passes; go vet ./cmd/rhiza ./internal/api is clean. Touches cmd/rhiza/main.go, ctx.go, detach_unix.go, tui.go, tui_open.go, internal/api/api.go and the internal/api, internal/index, internal/run and internal/gosweep files named in the statement, and their tests. Must not change any command's output, flags or exit codes.

## statement

run() takes its environment as a list, api.Ctx carries it as Environ with Getenv derived from it, and nothing on run()'s call path reads os.Getenv, os.LookupEnv or os.Environ or starts a child with a nil Env (which inherits the process's): gitIn, gitConfig, execIn, shellWithStdin and startProgram's programEnv in cmd/rhiza/ctx.go, detach_unix.go:19, tui.go's colorprofile.Detect calls, and internal/api's ship_build.go:199, usage_live.go:106, daemon_hook.go:87, ship_prepush.go:247, ship_generate.go:104, ship_gate.go:209, run_leaves.go:690, run.go:1657 (XDG_CONFIG_HOME) and project.go:228 (git clone) take it from the Ctx, and gosweep.CacheDir's callers on that path use CacheDirFrom(ctx.Getenv); nor does anything on that path call os.UserCacheDir, os.UserHomeDir or os.TempDir, which read the process's XDG_CACHE_HOME, HOME and TMPDIR: the user cache, home and temp dirs come from ctx.Getenv by the same rules, for the index path (index.DefaultPath, used at cmd/rhiza/ctx.go:117 and internal/api's question.go:175, inbox.go:1415, run_model_cache.go:127, slice.go:111 and daemon_projects.go:47, and index.Remove), index.IndexDir (internal/index/prune.go:79, used at internal/api/daemon_leftovers.go:38), the run's temp folder (run.TmpDir and Store.TmpDir, internal/run/store.go:49-52), personalWorkflowsDir's home (internal/api/run.go:1660) and gosweep.CacheDirFrom's user cache dir (internal/gosweep/gosweep.go:51); the daemon's self re-exec (internal/daemon/reexec.go) is the one exception

## Links

### Out

- depends-on [cli-no-process-globals](cli-no-process-globals.md)
- decided-by [run-takes-env-list](../../../../decision/run-takes-env-list.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)

### In

- depends-on [run-clock-from-ctx](run-clock-from-ctx.md)
- verifies [cli-child-env-verified](../../../../evidence/cli-child-env-verified.md)
- verifies [cli-child-env-fix-verified](../../../../evidence/cli-child-env-fix-verified.md)
- verifies [cli-child-env-qa-1](../../../../evidence/cli-child-env-qa-1.md)
- verifies [cli-child-env-qa-2](../../../../evidence/cli-child-env-qa-2.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
