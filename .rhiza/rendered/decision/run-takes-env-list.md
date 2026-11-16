# run-takes-env-list

decision · `01M7H0YQ7ZYEJNZ00TESB0B5EH`

## question

What replaces os.Environ for the processes rhiza starts when run() runs in-process?

## reasoning

A getenv func cannot be enumerated, and git, gh, the gate, the pre-push check and generate are started with cmd.Env = os.Environ() or with a nil Env (which inherits it); in-process, that is the test binary's environment, not the script's, so git would lose its hermetic HOME and config. One list carried on Ctx serves both reads and children.

## statement

run() takes its environment as a list in os.Environ() form instead of only a getenv func; api.Ctx gains Environ (that list), Getenv is derived from it, and every child process rhiza starts on run()'s call path gets its environment from ctx.Environ (with the existing additions, such as GIT_INDEX_FILE or RHIZA_GATE_HEAD, appended), never os.Environ. main() passes os.Environ(). The daemon's self re-exec (internal/daemon/reexec.go) keeps os.Args and os.Environ: it replaces the daemon process, which only runs as a process.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [in-process-scripts-plan-research](../source/in-process-scripts-plan-research.md)

### In

- decided-by [cli-no-process-globals](../capability/rhiza/integrity/self-check/cli-no-process-globals.md)
- decided-by [cli-child-env](../capability/rhiza/integrity/self-check/cli-child-env.md)
- decided-by [run-concurrent-guard](../capability/rhiza/integrity/self-check/run-concurrent-guard.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
