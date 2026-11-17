# spawn-ends-on-done-qa-1

evidence · `01M7MCHMMWTFP3PFMMHBWWP40W` · actor [claude-code-dd85ab3b](../actor/claude-code-dd85ab3b.md)

## result

PASS

## statement

QA built 8c10aeb1d to a scratch path and ran rhiza run start --foreground in a scratch root, with a fake claude whose step session calls rhiza run spawn and a fake rex that runs the spawned child in the background and kills it on rex block close. A child that ran rhiza session done and then idled: spawn returned outcome pass, detail said done, in 2 s, journaled spawn-ended pass, ran rex block close block:b1, and the child's process was gone. A child that exited 3 without a done mark: spawn waited about 6 s, returned fail with exit 3, rc 1, and closed no pane. A rex block close that errors: still pass, said done, rc 0, with a stderr line naming the pane that stays open. A child that said done and then exited 7: pass, said done. go test -trimpath ./internal/api -run 'TestRunSpawn|TestStepRunnerChild' and go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-spawn|describe-golden)$' pass; rhiza docs from the build leaves docs/reference unchanged in a detached worktree at HEAD.

## Links

### Out

- verifies [spawn-ends-on-done](../capability/rhiza/workflows/workflow-runs/spawn-ends-on-done.md)

## History

- `01M7M27J4MMS8X5PT4W3EXKDTM` Fix rhiza run spawn waits for the child agent's process to exit, and with no daemon nothing closes a child that already said rhiza session done, so spawn and its run hang until someone closes the child's pane
