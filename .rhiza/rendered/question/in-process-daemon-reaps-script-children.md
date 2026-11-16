# in-process-daemon-reaps-script-children

question · `01M7HSGGPARVH6P184DM02E54N`

## statement

An in-process rhiza daemon line (rhiza daemon --once) runs ctx.Adopt and ctx.Reap with the test binary's pid: adoptList lists every child of that process, other scripts' exec children too, and reapList Wait4s them, so a parallel script's exec fails with 'waitid: no child processes' (seen twice in TestScripts/^[a-d]: daemon-continue-slot:31, daemon-resume-landed-plan:37). Three a-d scripts already run the daemon loop in-process (daemon-event-hook, daemon-load-settings, daemon-notify-disk-low; workflow-plan-autonomy in s-w too). Should run() leave Adopt and Reap nil for an in-process call (a change to run()/ctx.go), or should those scripts go back on processScripts with exec rhiza? This retry kept the 18 other daemon scripts on the list for this reason instead of converting them.

## Links

### Out

- relates-to [scripts-in-process-a-d](../capability/rhiza/integrity/self-check/scripts-in-process-a-d.md) — found while converting the a-d daemon scripts; makes the a-d shard flaky

### In

- answers [in-process-run-no-adopt-reap](../decision/in-process-run-no-adopt-reap.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
