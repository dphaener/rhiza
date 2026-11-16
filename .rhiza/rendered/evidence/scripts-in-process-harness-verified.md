# scripts-in-process-harness-verified

evidence · `01M7HJNTMGRK4Q16GBZQ6KGQDA` · actor [claude-code-5c1c2f38](../actor/claude-code-5c1c2f38.md)

## result

PASS

## statement

At builder commit 370ea3fd3, under a capped systemd scope with session RHIZA_/SHIP_/REX_ vars unset: go test -p 2 -trimpath -count=1 ./cmd/rhiza -run '^Test(ProcessScripts|ScriptEnv|InProcessCommand)$' exit 0 (subtests off-list-fixture, process-list-fixture, script-env-fixture, in-process-fixture pass); TestScripts ranges ^[a-d] (29s), ^[e-q] (35s), ^r (60s), ^[s-w] (26s) each exit 0, covering all 398 scripts; -skip TestScripts exit 0; go vet ./cmd/rhiza clean; git diff --stat HEAD~1 HEAD -- cmd/rhiza/testdata empty. Diff read: RunMain's one entry is rhiza-process; TestMain runs rhizaProcess when os.Args[0] base is rhiza; rhizaLink hard-links rhiza beside rhiza-process once (sync.Once) and passes it to run() as self; setup sets RHIZA_TEST_SCRIPT; rhizaProcess exits 2 naming processScripts in cmd/rhiza/main_test.go for an off-list script; every txtar with an exec rhiza line is on processScripts (395 notConverted plus daemon-no-root and the two fixtures). No runaway: the rhiza entry calls run() once and never re-execs itself; no test hit the TasksMax cap. No stubs or TODOs; only cmd/rhiza/main_test.go, a new cmd/rhiza/script_inprocess_test.go and the actor record changed.

## Links

### Out

- verifies [scripts-in-process-harness](../capability/rhiza/integrity/self-check/scripts-in-process-harness.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
