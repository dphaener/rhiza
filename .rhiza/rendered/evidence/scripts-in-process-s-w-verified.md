# scripts-in-process-s-w-verified

evidence · `01M7HRJBKYPNWRD4T90W61CFA8` · actor [claude-code-eee2aab0](../actor/claude-code-eee2aab0.md)

## result

PASS

## statement

Detached checkout at 98d460f21. Capped (MemoryMax=6G, TasksMax=300), RHIZA_/SHIP_/REX_ unset: go test -p 2 -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/^[s-w]' exit 0 (24.3s); -run 'TestProcessScripts|TestInProcess' exit 0; go vet ./cmd/rhiza exit 0; gofmt -l cmd/rhiza empty. grep -cE '"[s-w][^"]*":.*notConverted' cmd/rhiza/main_test.go = 0; 398 txtar files; git diff --word-diff=porcelain HEAD^ HEAD -- cmd/rhiza/testdata/script shows 431 '-exec' removals and nothing else; diff touches only s-w txtar files, main_test.go's processScripts and .rhiza log/rendered. Failure case: removing the set entry and running TestScripts/^set$ exit 1 with 'rhiza: script set started rhiza as a process but is not on processScripts...' (reverted). Every s-w script with exec rhiza is on the list and every listed one keeps exec rhiza. All 50 kept reasons checked mechanically against the scripts (non-comment sh -c running rhiza, rhiza ... &, stdin feeding a rhiza line, bin/ fakes); by hand: set/show/stacked/watchdog-health sh -c rhiza, session-done-captures bin/claude calls rhiza, ship-land-branch-build bin/buildhead execs rhiza, session-journal bin/rex runs its tab command, usage-live-record/session-message-busy/session-hook-log-guard stdin into exec rhiza, session-message(-permission) rhiza with &. The 6 converted scripts with fakes (gh, rex, claude) never call rhiza back.

## Links

### Out

- verifies [scripts-in-process-s-w](../capability/rhiza/integrity/self-check/scripts-in-process-s-w.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
