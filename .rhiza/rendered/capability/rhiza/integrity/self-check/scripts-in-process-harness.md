# scripts-in-process-harness

capability · `01M7H1G7F352JVZ96KFBEDA1TF` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath -count=1 ./cmd/rhiza -run '^Test(ProcessScripts|ScriptEnv|InProcessCommand)$' passes, with new tests: TestProcessScripts shows that a fixture script off the list whose fake bin/ command calls rhiza fails with exit 2 and a message naming the list, and that the same script on the list runs; TestScriptEnv shows the built environment holds Setup's names, a value an `env NAME=x` line set, and `NAME=` for a name an `env NAME=` line cleared (a git child started by run() reads the script's HOME and GIT_CONFIG_GLOBAL); TestInProcessCommand shows a fixture script's bare `rhiza version` and `! rhiza nosuch` lines run with no process started (TestMain's `rhiza` entry counts none) and report the same stdout, stderr and exit status as `exec rhiza`. go test -trimpath -count=1 ./cmd/rhiza -run TestScripts passes with every script unchanged (git diff --stat cmd/rhiza/testdata is empty); go vet ./cmd/rhiza is clean. Touches cmd/rhiza/main_test.go and new test files in cmd/rhiza only. Must not change any command's output, flags or exit codes.

## statement

cmd/rhiza's TestScripts gives scripts a custom `rhiza` command that runs run() in-process with the script's working directory, the environment built per scripts-in-process-env-from-setup and the script's standard streams; testscript.RunMain's one entry is renamed `rhiza-process`, TestMain runs run() itself when started as `rhiza`, and TestScripts hard-links a `rhiza` beside `rhiza-process` on PATH and passes its path to run() as its own; Setup gives every script RHIZA_TEST_SCRIPT=<name>, and TestMain's `rhiza` entry exits 2, naming the list, when that script is not on one list of process scripts in cmd/rhiza/main_test.go, each with its reason. The list starts with every script whose `exec rhiza` lines are not yet converted, under the reason `not converted yet`; no txtar file changes

## Links

### Out

- depends-on [cli-no-process-globals](cli-no-process-globals.md)
- depends-on [run-concurrent-guard](run-concurrent-guard.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)
- decided-by [scripts-mark-process-lines-exec](../../../../decision/scripts-mark-process-lines-exec.md)
- decided-by [scripts-in-process-env-from-setup](../../../../decision/scripts-in-process-env-from-setup.md)

### In

- depends-on [scripts-in-process-a-d](scripts-in-process-a-d.md)
- depends-on [scripts-in-process-e-q](scripts-in-process-e-q.md)
- depends-on [scripts-in-process-r](scripts-in-process-r.md)
- depends-on [scripts-in-process-s-w](scripts-in-process-s-w.md)
- verifies [scripts-in-process-harness-verified](../../../../evidence/scripts-in-process-harness-verified.md)
- verifies [scripts-in-process-harness-qa-1](../../../../evidence/scripts-in-process-harness-qa-1.md)
- verifies [scripts-in-process-harness-qa-2](../../../../evidence/scripts-in-process-harness-qa-2.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
