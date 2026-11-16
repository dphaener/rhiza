# scripts-in-process-r

capability · `01M7H1G7J4D54HBAJHZ73NQ84R` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/^r' passes. No process-list entry for a script starting with r keeps the reason `not converted yet` (grep in cmd/rhiza/main_test.go finds none). Failure case: such a script off the list that still has an `exec rhiza` line fails with the list's message. Touches the txtar files starting with r and the process list in cmd/rhiza/main_test.go only. Must not delete or skip a script (the count of cmd/rhiza/testdata/script/*.txtar stays 398) or weaken an assertion: each changed line differs only in its exec prefix (git diff --word-diff shows only removed `exec `).

## statement

The scripts in cmd/rhiza/testdata/script whose names start with r (117 scripts, about 1,150 lines on 2026-11-16) run their rhiza lines in-process: each `exec rhiza` line becomes `rhiza`, except in the scripts that need a process (a fake bin/ command that calls rhiza back, rhiza with &, a stdin or ttyin feed, a command that starts ctx.Self, a run or daemon lock a later line must find interrupted because the pid it recorded must have exited, or a difference between unset and empty that shows), which keep `exec rhiza` and whose process-list entry's reason `not converted yet` is replaced by the reason they need a process; the converted scripts leave the list

## Links

### Out

- depends-on [scripts-in-process-harness](scripts-in-process-harness.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)
- decided-by [scripts-mark-process-lines-exec](../../../../decision/scripts-mark-process-lines-exec.md)
- decided-by [scripts-in-process-env-from-setup](../../../../decision/scripts-in-process-env-from-setup.md)

### In

- depends-on [scripts-in-process-docs](scripts-in-process-docs.md)
- verifies [scripts-in-process-r-verified-2e88f62f5](../../../../evidence/scripts-in-process-r-verified-2e88f62f5.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
