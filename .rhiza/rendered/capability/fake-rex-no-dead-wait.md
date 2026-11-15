# fake-rex-no-dead-wait

capability · `01M7FT144PPDNCMCPTJZ31CP1Z`

## acceptance

go test -trimpath -p 2 -count=1 -v ./cmd/rhiza -run 'TestScripts/^(bugfix-workflow-autonomous|plan-workflow-autonomous-left|plan-workflow-autonomous|plan-review-workflow)$' passes, each subtest under 15 s run alone on the VM (2026-11-15: 91, 64, 59 and 58 s), with each script's stdout, stderr and file assertions unchanged (git diff shows only the fake rex sections changed). Every other script edited for the same fall-through passes the same way, run by name. Check: grep -l 'never got ready' cmd/rhiza/testdata/script/*.txtar lists only scripts whose fake rex either answers block inspect or reaches its ready loop only past a -- argument. Touches cmd/rhiza/testdata/script/bugfix-workflow-autonomous.txtar, plan-workflow-autonomous-left.txtar, plan-workflow-autonomous.txtar, plan-review-workflow.txtar, and any other script with the same fall-through. Must not touch internal/ or the scaled clock.

## statement

No fake rex in cmd/rhiza's scripts waits for an agent's ready file on a call that starts no command: bugfix-workflow-autonomous, plan-workflow-autonomous-left, plan-workflow-autonomous and plan-review-workflow answer rex block inspect with a window id, and every other script whose fake rex has the ready-file loop runs it only after starting a pane's command, so those four scripts take seconds, not a minute per run start.

## Links

### Out

- decided-by [fake-rex-answers-inspect](../decision/fake-rex-answers-inspect.md)

### In

- depends-on [script-timings-measured](script-timings-measured.md)
- verifies [fake-rex-no-dead-wait-verified](../evidence/fake-rex-no-dead-wait-verified.md)
- verifies [fake-rex-no-dead-wait-qa-1](../evidence/fake-rex-no-dead-wait-qa-1.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
