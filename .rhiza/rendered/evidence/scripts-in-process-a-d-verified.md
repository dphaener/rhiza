# scripts-in-process-a-d-verified

evidence · `01M7HK9V85P5QESDTJ5FA5PMRK` · actor [claude-code-f74006e2](../actor/claude-code-f74006e2.md)

## result

PASS

## statement

Detached checkout at e2c6bd895, capped (systemd-run MemoryMax=6G TasksMax=300, -p 2, RHIZA/SHIP/REX env unset). go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/^[a-d]': ok (21.1s). TestProcessScripts|TestInProcessCommand|TestScriptEnv: ok. grep for [a-d] notConverted entries in cmd/rhiza/main_test.go: 0. txtar count: 398. git diff --word-diff=porcelain HEAD^ HEAD on testdata/script: only 707 '-exec' removals. Only files changed outside .rhiza/: main_test.go and [a-d] txtar files. Failure case: carve.txtar (off the list) with its first rhiza line made exec rhiza, TestScripts/^carve$ fails with 'script carve started rhiza as a process but is not on processScripts ...' (probe reverted). Per-script audit of all 96 [a-d] scripts: every off-list script has zero exec rhiza lines; every kept script matches its reason: 47 with fake bin/ rex/claude/gh/systemctl/loginctl and PATH=$WORK/bin (confirmed in ctx.go: execIn/startProgramIn use exec.Command(name), which resolves via the test binary's own PATH, not cmd.Env), 15 with exec sh -c 'rhiza ...' lines (spot-checked check-move, check-git, change-set-ref), apply/apply-concurrent (internal/api/apply.go os.ReadFile(file) resolves relative to the process cwd), describe (stdin feed), custom-workflows-doc (quickstart-doc). The fake-bin PATH reason goes beyond the leaf's enumerated list but is a real need, and the gap is filed as draft change set 01M7HK31WRYQ54W29MJCHKGGXK.

## Links

### Out

- verifies [scripts-in-process-a-d](../capability/rhiza/integrity/self-check/scripts-in-process-a-d.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
