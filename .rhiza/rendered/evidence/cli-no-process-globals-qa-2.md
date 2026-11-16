# cli-no-process-globals-qa-2

evidence · `01M7K0CW889P4YXWR0R2J5JJSE` · actor [claude-code-037cc825](../actor/claude-code-037cc825.md)

## result

PASS

## statement

QA cycle 2: go test ./cmd/rhiza -skip TestScripts (TestNoProcessGlobals with userdir/envexec and MkdirTemp fixtures, required ownChildren) passes; no production api.Ctx or run store relies on a nil getenv (all run stores that call TmpDir come from runStore(ctx)).

## Links

### Out

- verifies [cli-no-process-globals](../capability/rhiza/integrity/self-check/cli-no-process-globals.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
