# cli-child-env-qa-1

evidence · `01M7HZJ7KM58BKED99E1H5JPY8` · actor [claude-code-10d5f94b](../actor/claude-code-10d5f94b.md)

## result

FAIL

## statement

QA: with a fake git first on an absolute PATH and a marker variable, check / ship check / run ls ran the same git calls with the marker visible, same as the base build. With PATH='.:/usr/bin' or ':/usr/bin' and a git file in the working directory, the new build runs that planted git (12 calls) while the base build refused it (Go ErrDot): recorded as a finding.

## Links

### Out

- verifies [cli-child-env](../capability/rhiza/integrity/self-check/cli-child-env.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
