# scripts-in-process-docs

capability · `01M7H1G7M59QWBW902Z4WN442Y` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 ./cmd/rhiza -run '^Test(AgentDocsConsistent|DevelopingDocLessons|ClaudeMD|ProcessScripts|ProcessListConverted)$' passes, with a new test TestProcessListConverted that fails, naming the script, on a list entry whose reason is `not converted yet` (shown with a fixture list) and passes on the real list; go test -trimpath -count=1 ./cmd/rhiza passes. The commit's Checks: block records, on a quiet VM, /usr/bin/time -f '%e s %S sys' go test -trimpath -count=1 ./cmd/rhiza at the commit before the harness and after, and the number of rhiza processes started (counted by TestMain's `rhiza` entry) before and after (2026-11-15: about 4,870 exec rhiza lines in 398 scripts). Touches cmd/rhiza/main_test.go, AGENTS.md and agents/developing.md only.

## statement

No process-list entry keeps the reason `not converted yet`, and a test fails on one that does; AGENTS.md's Tests bullet and agents/developing.md:21 say that a bare `rhiza` line runs in-process and an `exec rhiza` line starts a process, allowed only in a script on the process list in cmd/rhiza/main_test.go; the before/after cost of cmd/rhiza's suite is recorded

## Links

### Out

- depends-on [scripts-in-process-a-d](scripts-in-process-a-d.md)
- depends-on [scripts-in-process-e-q](scripts-in-process-e-q.md)
- depends-on [scripts-in-process-r](scripts-in-process-r.md)
- depends-on [scripts-in-process-s-w](scripts-in-process-s-w.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)
- decided-by [scripts-mark-process-lines-exec](../../../../decision/scripts-mark-process-lines-exec.md)
- decided-by [scripts-in-process-env-from-setup](../../../../decision/scripts-in-process-env-from-setup.md)

### In

- verifies [scripts-in-process-docs-verified](../../../../evidence/scripts-in-process-docs-verified.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
