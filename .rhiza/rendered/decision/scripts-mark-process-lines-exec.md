# scripts-mark-process-lines-exec

decision · `01M7H0YQ3NV130SDW77EGD0NP4`

## question

How do in-process script lines and the lines that need a real rhiza process coexist, and how is the list of process scripts enforced?

## reasoning

go-internal v1.14.1's testscript looks up scriptCmds (builtins, and every RunMain command, which RunMain registers there: testscript/exe.go:93) before Params.Cmds (testscript.go:723-725), so neither `exec rhiza` nor a bare `rhiza` can be redirected in-process while RunMain has a `rhiza` entry; RunMain's entry is renamed and TestMain dispatches the `rhiza` name itself, which RunMain's own basename dispatch would otherwise do. The scripts' lines have to change, and the `exec` prefix is already the honest marker of a process start. Custom commands cannot run with & and cannot read the `stdin` builtin's input (ts.stdin is private), so those lines must stay processes. A hard link to the binary RunMain already placed on PATH gives a working `rhiza` with no separate build. Enforcing the list at the process entry, by script name, catches every way a process gets started, including ctx.Self children and fake bin/ callbacks, which a text scan of the script would miss.

## statement

In cmd/rhiza's scripts a bare `rhiza ...` line runs in-process, through a custom testscript command `rhiza` in TestScripts' Cmds that calls run() with the script's working directory (ts.MkAbs(".")), environment, ts.Stdout() and ts.Stderr(); `exec rhiza ...` stays the way to start a process. For the bare line to reach Cmds, testscript.RunMain gets no `rhiza` entry: its one entry is named `rhiza-process`, and TestMain itself, before calling RunMain, runs run() when the binary was started under the name `rhiza` (filepath.Base(os.Args[0])); TestScripts adds, once per test binary, a hard link named `rhiza` beside RunMain's `rhiza-process` in the bin directory RunMain put on PATH, so `exec rhiza`, fake bin/ callbacks and ctx.Self (run() is given that link's path as its own) start the test binary as rhiza and no second binary is built. A script may have `exec rhiza` lines, or start rhiza as a process any other way (a fake bin/ command that calls rhiza, rhiza with &, a `stdin` or `ttyin` feed, a command that starts ctx.Self such as run start or daemon run, or a line whose run or daemon lock a later line must find interrupted: a run records the pid of the process that ran it (run.ThisProcess, internal/run/process.go:70) and the daemon lock its holder's (internal/daemon/lock.go:83), and in-process that pid is the test binary, which stays alive, so such a line, for example run-abandon.txtar:17's `rhiza run start --foreground --no-wait`, needs a process that has exited), only when it is named in one list in cmd/rhiza/main_test.go with its reason; Setup gives every script RHIZA_TEST_SCRIPT=<name>, and TestMain's `rhiza` entry exits 2 with a message naming the list when that script is not on it.

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
