# cli-no-process-globals

capability · `01M7GG0GD7S0MZQ888A108TB0C` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 ./cmd/rhiza -run '^Test(NewCtx|CtxFromEnv|TUI|RunWorkingDir)' passes, with a new test TestRunWorkingDir that, with the process's working directory outside any root, calls run() with a scratch repo's folder as its working directory and shows the command answers for that repo, and that an error it reports lands in the stderr buffer it was given; failure case in the same test: a working directory that does not exist fails with the 'working directory' error and exit code 1. go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/tui-' passes and go vet ./cmd/rhiza is clean. Touches cmd/rhiza/main.go, ctx.go, tui.go, tui_open.go, main_test.go's TestMain entry and docs_examples_test.go's run() call. Must not change any command's output, flags or exit codes, nor environment handling (cli-child-env).

## statement

cmd/rhiza's run() takes its working directory and its own binary's path as parameters instead of reading os.Getwd (ctx.go:42) and os.Args[0] (selfPath, ctx.go:421); tui.go opens a project's app by passing the project's folder as the working directory instead of os.Chdir (tui.go:147-276 and wdNow), tui_open.go's openModal uses those parameters instead of os.Executable and os.Getwd, and ctx.Stderr, ctx.Terminal and ctx.Rex are the stderr run() was given, not os.Stderr; main() passes os.Getwd(), selfPath() and os.Stderr

## Links

### Out

- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)
- decided-by [run-takes-env-list](../../../../decision/run-takes-env-list.md)

### In

- depends-on [cli-child-env](cli-child-env.md)
- depends-on [scripts-in-process-harness](scripts-in-process-harness.md)
- verifies [cli-no-process-globals-verified](../../../../evidence/cli-no-process-globals-verified.md)
- verifies [cli-no-process-globals-fix-verified](../../../../evidence/cli-no-process-globals-fix-verified.md)
- verifies [cli-no-process-globals-adopt-reap-verified](../../../../evidence/cli-no-process-globals-adopt-reap-verified.md)
- verifies [cli-no-process-globals-qa-1](../../../../evidence/cli-no-process-globals-qa-1.md)
- verifies [cli-no-process-globals-qa-2](../../../../evidence/cli-no-process-globals-qa-2.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
