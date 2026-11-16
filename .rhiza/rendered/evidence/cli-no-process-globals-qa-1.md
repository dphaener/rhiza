# cli-no-process-globals-qa-1

evidence · `01M7HZJ7GJ2RTSSGJBS3N7FNN9` · actor [claude-code-10d5f94b](../actor/claude-code-10d5f94b.md)

## result

PASS

## statement

QA: built ./cmd/rhiza at HEAD and at the merge-base to scratch paths, ran init, change new, apply (relative file from repo/sub/deeper, ./missing.json, ../deeper/w.json --dry-run), source add with relative file and --raw (and missing ones), slice, check, render in scratch repos under a scratch HOME with RHIZA_TEST_NOW/SEED: output and exit codes identical to the base build. daemon --once (via sh -c 'sleep 1 & exec rhiza daemon --once') dropped a stale drivers entry and adopted its sleep child exactly as the base build did.

## Links

### Out

- verifies [cli-no-process-globals](../capability/rhiza/integrity/self-check/cli-no-process-globals.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
