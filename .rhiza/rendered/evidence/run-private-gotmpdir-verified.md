# run-private-gotmpdir-verified

evidence · `01M6XP2MTP58DJF1BRVX1X89B9` · actor [claude-code-77708b30](../actor/claude-code-77708b30.md)

## result

PASS

## statement

Verified at 613bc62b (commits b2b5a2be and 613bc62b). Read the diffs: run.TmpDir is <git common dir>/rhiza/runs/<id>/gotmp, separate from the <id>.jsonl journal and lock files. The engine's stepEnv creates it and adds GOTMPDIR to every command, wait and agent step, and session.run creates it and adds GOTMPDIR to the agent's env. The engine.run defer removes <runs>/<id> under the run lock when the run finishes or fails, and runRunAbandon removes it under the lock after journaling the end, which covers attempt 1's finding f1. Ran go vet ./internal/run ./internal/api and go test -trimpath -count=1 ./internal/run/... ./internal/api/... -run 'GoTmp|TestSessionRun$|TestCommandPassFollowsOnPass|TestArbitratePasses|Abandon'. All passed: TestStepsGetRunPrivateGoTmp checks GOTMPDIR on command, wait and agent steps, that the folder exists when each step runs, and that it is gone after the run succeeds. TestGoTmpGoneAfterFailedRun covers a failed run, TestAbandonRemovesRunGoTmp covers abandon (the journal stays readable and the run is abandoned), and TestSessionRun checks the agent env and that the folder exists. Ran the abandon test in a throwaway worktree at 2bee5d01 (without the fix) and saw it fail with 'run folder ... still there after abandon', so the test does check the fix.

## Links

### Out

- verifies [run-private-gotmpdir](../capability/rhiza/workflows/workflow-runs/run-private-gotmpdir.md)

## History

- `01M6WGEPSYWV8R14VGKJX9GGJM` Fix A Go build inside a run step fails when its $WORK temp folder vanishes mid-build, and the run ends: plan run 01M6W7G5J0M6 (change set 01M6W7G8YBXB) failed at step rebase on 2026-11-08 ~00:15Z with 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' for many std packages while ~5 agents built and tested on the 4-core VM (load 28). The daemon's go build sweep (gosweep, TempMaxAge 2h) logged nothing in that window, so something else removed a live ~/tmp/go-build* folder: find what (candidates: a test's TestMain or tmpsweep removing TMPDIR entries by prefix, an agent's cleanup, go's own cleanup under a shared TMPDIR) and stop it. Also: a step whose command fails with this signature should retry once instead of ending the run. Unstuck by rhiza run resume --retry 01M6W7G5J0M6.
