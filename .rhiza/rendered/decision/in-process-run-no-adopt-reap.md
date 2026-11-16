# in-process-run-no-adopt-reap

decision · `01M7HSNAZKWVQXZT9DRBGRAGCA`

## reasoning

Adopting and reaping every child of the current pid is a process-wide act, like os.Getwd, os.Args and the daemon's self re-exec, which this change set moves out of run() into main() or the parameters main() passes (cli-no-process-globals, cli-child-env, run-concurrent-guard). In a shared test binary it is wrong by construction: adoptList lists other scripts' exec children and reapList Wait4s them, giving 'waitid: no child processes' flakes (daemon-continue-slot:31, daemon-resume-landed-plan:37). internal/api already treats nil Adopt and Reap as a no-op (daemon.go:171 and 498, daemon_projects.go:223), so the change is confined to cmd/rhiza's ctx construction and main(); the daemon's real behaviour is unchanged because main() still passes both, and TestReap/TestReapAdopt keep covering them directly. Keeping 28 daemon scripts on processScripts instead would forgo the change set's intent (no process start per script line) for the daemon scripts, and leaves the three a-d scripts and workflow-plan-autonomy that already run the loop in-process flaky unless they too were reverted.

## statement

run() leaves api.Ctx.Adopt and Reap nil unless its caller passes them: main() (the real rhiza process) passes adoptChildren and reapDrivers, an in-process call (the testscript rhiza command, TestRunConcurrentDirs, any programmatic caller) passes none, so only a real process adopts and reaps its own children. The scripts kept on processScripts only because 'rhiza daemon adopts and reaps its process's children' are then converted to in-process lines unless another listed reason keeps them there.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-524bf720-1](../question/question-524bf720-1.md)
- answers [in-process-daemon-reaps-script-children](../question/in-process-daemon-reaps-script-children.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
