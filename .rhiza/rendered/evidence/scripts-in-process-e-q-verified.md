# scripts-in-process-e-q-verified

evidence · `01M7HM77CSMARE8SWCAHSVE16K` · actor [claude-code-f4d26330](../actor/claude-code-f4d26330.md)

## result

PASS

## statement

Detached checkout at b8da7ea13 (diff 878916791^..b8da7ea13). Capped go test -p 2 -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/^[e-q]': ok 31.4s. grep notConverted for e-q entries in cmd/rhiza/main_test.go: 0. Failure case: health (off the list) with its first rhiza line turned back to exec rhiza fails with 'script health started rhiza as a process but is not on processScripts...' (reverted). ls testdata/script/*.txtar: 398. git diff --word-diff=porcelain on testdata/script: only 809 '-exec' removals. Checks block re-run: the 7 removed entries plus TestProcessScripts|TestInProcessCommand|TestScriptEnv pass; gofmt -l cmd/rhiza empty. Every kept e-q entry was checked against its script: the 23 'sh -c' entries each have a sh -c line that runs rhiza (so the 7 earlier mislabelled entries are fixed); fake bin/ entries have the bin/ files named and keep exec rhiza lines; & reasons match exec rhiza ... & lines; init has stdin before exec rhiza init; e2e, leaf-cycle, import-retires-docs and import-unbuilt-as-drafts keep exec only on the apply / source add lines that read a relative file; quickstart uses quickstart-doc. Only e-q txtar files, their processScripts lines and .rhiza log/rendered changed.

## Links

### Out

- verifies [scripts-in-process-e-q](../capability/rhiza/integrity/self-check/scripts-in-process-e-q.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
