# cli-no-process-globals-adopt-reap-verified

evidence · `01M7HT29A7XHZ6GMKDN12ZR8ST` · actor [claude-code-0ddb25f8](../actor/claude-code-0ddb25f8.md)

## result

PASS

## statement

Re-verified cli-no-process-globals at fix 90b43b9be (Decision in-process-run-no-adopt-reap) in a detached worktree, every go test capped (MemoryMax=6G, TasksMax=300, -p 2, RHIZA/SHIP/REX vars unset). Checks: TestInProcessDaemonLeavesSiblings exit 0 (in-process subtest: sibling child not on the driver list and its Wait gets exit 3; process subtest with processChildren: sibling adopted). Old wiring (ctx.go Adopt/Reap set back to adoptChildren/reapDrivers, scratch only, reverted) exit 1: in-process sibling on the driver list, so the test bites. -skip TestScripts exit 0 (27s); TestScripts/^(daemon-event-hook|daemon-load-settings|daemon-notify-disk-low|daemon-continue-slot|daemon-resume-landed-plan)$ exit 0; ship generate and render exit 0 with a scratch build (only drift is other sessions' actors from the shared root). Acceptance: -run '^Test(NewCtx|CtxFromEnv|TUI|RunWorkingDir)' exit 0, TestScripts/tui- exit 0, go vet ./cmd/rhiza clean, gofmt clean. Diff: only main() and main_test.go rhizaProcess pass processChildren; every other run() caller (testscript rhiza command main_test.go:515, ctx_test, process_globals_test, docs_examples_test) passes none; tui.go's newCtx copies g so own carries; internal/api checks Adopt/Reap for nil at daemon.go:171, 498 and daemon_projects.go:224, the only uses. No output, flag or exit-code change. Variadic own: sound for today's callers, but it hides two mistakes: more than one value is accepted with the last silently winning, and a future real-process entry point that forgets it compiles and silently stops reaping drivers (zombies) with no test to catch it; a required parameter would make each caller choose.

## Links

### Out

- verifies [cli-no-process-globals](../capability/rhiza/integrity/self-check/cli-no-process-globals.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
