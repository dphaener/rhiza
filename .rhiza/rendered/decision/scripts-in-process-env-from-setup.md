# scripts-in-process-env-from-setup

decision · `01M7H0YQ5CKJQJ8845KC2NVCGJ`

## question

testscript has no accessor for a script's whole environment; how does the in-process rhiza command get the environment it passes to run()?

## reasoning

testscript v1.14.1 exposes ts.Getenv per name but not the whole environment, and the children rhiza starts (git above all, whose hermetic config comes from HOME and GIT_CONFIG_GLOBAL) need the whole environment the script would have given an exec'd rhiza. The names in play are exactly Setup's plus those the script's env lines set, and both are known: Setup's from env.Vars, the script's from its own text. Upgrading go-internal for an accessor would be a dependency change, which a plan must not add.

## statement

The in-process `rhiza` command builds the environment it passes to run() as the script's Setup environment (env.Vars, kept per script with env.Values) with, for every name that environment holds or that one of the script's `env NAME=...` lines sets, the current value from ts.Getenv; a name whose value is empty is kept as `NAME=` (set and empty), as testscript's `env NAME=` appends `NAME=` to the environment an exec'd rhiza and its children get (testscript.go:1237, 983, 1036). ts.Getenv cannot tell a name a later env line sets, still unset, from an empty one, so such a name is passed set and empty too; a script where that difference shows (a fake bin/ command or check that tells unset from empty, such as ${NAME-unset}) is put on the process list.

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [in-process-scripts-plan-research](../source/in-process-scripts-plan-research.md)

### In

- decided-by [scripts-in-process-harness](../capability/rhiza/integrity/self-check/scripts-in-process-harness.md)
- decided-by [scripts-in-process-a-d](../capability/rhiza/integrity/self-check/scripts-in-process-a-d.md)
- decided-by [scripts-in-process-e-q](../capability/rhiza/integrity/self-check/scripts-in-process-e-q.md)
- decided-by [scripts-in-process-r](../capability/rhiza/integrity/self-check/scripts-in-process-r.md)
- decided-by [scripts-in-process-s-w](../capability/rhiza/integrity/self-check/scripts-in-process-s-w.md)
- decided-by [scripts-in-process-docs](../capability/rhiza/integrity/self-check/scripts-in-process-docs.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
