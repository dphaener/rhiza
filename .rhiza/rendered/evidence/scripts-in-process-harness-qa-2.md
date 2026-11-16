# scripts-in-process-harness-qa-2

evidence · `01M7K0CWDB3CEW2JZ4HXCVQ6P8` · actor [claude-code-037cc825](../actor/claude-code-037cc825.md)

## result

PASS

## statement

QA cycle 2: session-journal, printout-counts-used-tokens, printout-gap-names-reason, run-autonomous pass in under 1 s each; TestScripts/^r passes in 63 s under TasksMax=300 with no runaway.

## Links

### Out

- verifies [scripts-in-process-harness](../capability/rhiza/integrity/self-check/scripts-in-process-harness.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
