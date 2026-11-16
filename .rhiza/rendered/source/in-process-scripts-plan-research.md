# in-process-scripts-plan-research

source · `01M7H0WW5QJXF6E8K7R8SHJS76`

## quotes

```json
[]
```

## raw_sha256

d1fdf96445c4df9659dec75335a0d5cd5dd65022eeacfc3d3ca37f3b99fe5957

## session

d3d753d3-18bc-4528-97d1-01a1e20afb72

## summary

Autonomous plan run for change set 01M7GFTS0N (no owner; the orchestrator acts as arbiter). The arbiter's brief: design choices are the planner's, recorded as '(assumed, would have escalated)' with reasoning; keep each leaf one concern a fresh session can finish and verify; stay out of change set 01M7GFTDY5. Research findings: scripts invoke rhiza as `exec rhiza` (about 4,870 lines in 398 scripts), and go-internal v1.14.1's testscript looks up builtin commands before Params.Cmds, so `exec rhiza` cannot be intercepted: in-process lines need a custom `rhiza` command and the scripts must change. Custom commands cannot run in the background (&) and cannot read the `stdin` builtin's input (ts.stdin is private); 12 scripts use stdin or ttyin, 6 run rhiza with &. testscript exposes ts.Getenv per name but no way to list the whole environment. testscript.RunMain already puts a `rhiza` on PATH that re-executes the test binary, so no separate build is needed for the process lines. Process-wide state on run()'s call path beyond os.Getwd (ctx.go:42) and os.Args[0] (selfPath, ctx.go:421): os.Chdir in tui.go (147-276), os.Executable and os.Getwd in tui_open.go:62-66, ctx.Stderr/Terminal/Rex set to os.Stderr (ctx.go:83-94), child processes inheriting os.Environ (gitIn, gitConfig and execIn with nil env in ctx.go; programEnv(os.Environ()) in startProgram; internal/api ship_build.go:199, usage_live.go:106, daemon_hook.go:87, ship_prepush.go:247, ship_generate.go:104, ship_gate.go:209, run_leaves.go:690), os.Getenv reads (internal/api/run.go:1657 XDG_CONFIG_HOME, gosweep.CacheDir, detach_unix.go:19), and runClock (internal/api/run.go:100), chosen once at package init from RHIZA_TEST_CLOCK_SCALE in the process environment, which in-process would put every script's waits on the real clock. internal/daemon/reexec.go uses os.Args and os.Environ for the daemon's self re-exec, which only runs in a daemon process. Package-level vars in internal/api are test seams that commands do not write.

## title

Planning Rhiza in-process in scripts (autonomous)

## Links

### In

- derived-from [scripts-mark-process-lines-exec](../decision/scripts-mark-process-lines-exec.md)
- derived-from [scripts-in-process-env-from-setup](../decision/scripts-in-process-env-from-setup.md)
- derived-from [run-takes-env-list](../decision/run-takes-env-list.md)
- derived-from [run-clock-chosen-per-ctx](../decision/run-clock-chosen-per-ctx.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
