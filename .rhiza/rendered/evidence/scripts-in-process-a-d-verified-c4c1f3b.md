# scripts-in-process-a-d-verified-c4c1f3b

evidence · `01M7HTR70S8XCB6QN1RJPY9P3Y` · actor [claude-code-ae4bbd36](../actor/claude-code-ae4bbd36.md)

## result

PASS

## statement

Verified scripts-in-process-a-d at c4c1f3b72 (fixes 0199055c2, c4c1f3b72) in a detached checkout, with RHIZA_/SHIP_/REX_ unset, stand-ins for rex, gh and claude first on PATH, under a capped systemd scope. go test -p 2 -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/^[a-d]' ran twice back to back: both exit 0 (22s), no 'waitid: no child processes'. main_test.go has 0 a-d entries with notConverted and 0 with 'adopts and reaps its process'. 398 txtar files. git diff --word-diff=porcelain on the txtar files shows only removed 'exec ' (22 in 0199055c2, 386 in c4c1f3b72). Failure case: describe taken off processScripts gives exit 1 with 'script describe started rhiza as a process but is not on processScripts...' (reverted). Checked all 46 kept a-d entries against their scripts: every 'sh -c' reason has a non-comment sh -c line that runs rhiza (autonomy-accept-rule and daemon-resume-landed included); '&' reasons have &name& lines; apply and describe have a stdin feed; the 'fake rex runs its tab's command' reasons match rex bodies that exec the args after -- (autonomous-permissions, change-request-kind, daemon-start-stop and others); the claude reasons match fakes that call rhiza; custom-workflows-doc runs doc commands through quickstart-doc. daemon-headless-runs probe (converted and off the list) fails 3 of 3 runs, so its ctx.Self reason holds. The leaf's commits change only a-d txtar files and a-d processScripts entries.

## Links

### Out

- verifies [scripts-in-process-a-d](../capability/rhiza/integrity/self-check/scripts-in-process-a-d.md)

## History

- `01M7GFTS0NVWBXF3CQZD916J1C` cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.
