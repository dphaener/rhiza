# step-retry-vanished-go-work-verified

evidence · `01M6XQ8SSD2ZFZMFA7A75ZJ0AE` · actor [claude-code-307713d5](../actor/claude-code-307713d5.md)

## result

PASS

## statement

At 2bee5d01 ran go test -trimpath -count=1 ./internal/api -run 'TestVanishedGoWork|TestCommandRetriesVanishedGoWork' -v: PASS. TestVanishedGoWork accepts the 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' and "can't create $WORK/b094/_pkg_.a: open $WORK/b094/_pkg_.a: no such file or directory" lines and rejects a compile error, a missing file outside $WORK and a non-ENOENT $WORK error. TestCommandRetriesVanishedGoWork drives stepRunner.Command with scripted fake processes: signature-then-pass returns exit 0 after 2 runs with exactly one 'running it once more' stderr line; signature twice returns the second failure's code 2 after 2 runs; a plain failure runs once with no retry line. Read the diff: Command reruns only when err is nil, the code is non-zero and not StopExit, and the output matches open \$WORK/b\d+/...: no such file or directory; tail and buf are reset so the second result's detail stands. gofmt and go vet are clean. The run_test, run_leaves, run_checks_panes and run_gowork tests pass from a scratch worktree outside ~/code; inside ~/code only the rootguard tripped, on usage.live.json that the live root daemon keeps rewriting. A full ./internal/api run timed out under load average ~17 (TestReadBudget 290s, known contention), in tests this change does not touch.

## Links

### Out

- verifies [step-retry-vanished-go-work](../capability/rhiza/workflows/workflow-runs/step-retry-vanished-go-work.md)

## History

- `01M6WGEPSYWV8R14VGKJX9GGJM` Fix A Go build inside a run step fails when its $WORK temp folder vanishes mid-build, and the run ends: plan run 01M6W7G5J0M6 (change set 01M6W7G8YBXB) failed at step rebase on 2026-11-08 ~00:15Z with 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' for many std packages while ~5 agents built and tested on the 4-core VM (load 28). The daemon's go build sweep (gosweep, TempMaxAge 2h) logged nothing in that window, so something else removed a live ~/tmp/go-build* folder: find what (candidates: a test's TestMain or tmpsweep removing TMPDIR entries by prefix, an agent's cleanup, go's own cleanup under a shared TMPDIR) and stop it. Also: a step whose command fails with this signature should retry once instead of ending the run. Unstuck by rhiza run resume --retry 01M6W7G5J0M6.
