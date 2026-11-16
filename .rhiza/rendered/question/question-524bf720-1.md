# question-524bf720-1

question · `01M7HSM8R5VWH76XGXF593HN37`

## session

524bf720-e097-40a1-9121-01a1e2d6798e

## statement

Question in-process-daemon-reaps-script-children: an in-process `rhiza daemon --once` line runs ctx.Adopt and ctx.Reap against the test binary's pid. adoptList takes every child of the process, including other scripts' exec children, and reapList Wait4s them, so a parallel script's exec fails with 'waitid: no child processes'. This made TestScripts/^[a-d] flaky (daemon-continue-slot:31, daemon-resume-landed-plan:37). Should run() leave Adopt and Reap nil unless main() passes them (only the real process adopts and reaps), with the daemon scripts then converted in-process? Or should every script that runs the daemon loop stay on processScripts with exec rhiza?

## Links

### In

- answers [in-process-run-no-adopt-reap](../decision/in-process-run-no-adopt-reap.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
