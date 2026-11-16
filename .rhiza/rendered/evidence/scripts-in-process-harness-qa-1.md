# scripts-in-process-harness-qa-1

evidence · `01M7HZJ7SJVKJ75MV3FJS9SW2Z` · actor [claude-code-10d5f94b](../actor/claude-code-10d5f94b.md)

## result

PASS

## statement

QA: in a scratch worktree, scratch txtars: bare rhiza lines (apply and source add with relative files from a subdir, ! rhiza nosuch, env RHIZA_ACTOR= and env GIT_CONFIG_GLOBAL= lines) passed; exec rhiza off the list and a fake bin/ command calling rhiza both failed in about 0.02 s with exit 2 and the processScripts message. The leaves' named unit tests in cmd/rhiza, internal/envexec, internal/userdir and internal/run pass.

## Links

### Out

- verifies [scripts-in-process-harness](../capability/rhiza/integrity/self-check/scripts-in-process-harness.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
