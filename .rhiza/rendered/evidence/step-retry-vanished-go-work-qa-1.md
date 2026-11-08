# step-retry-vanished-go-work-qa-1

evidence · `01M6XQM7XPH98W84VE9P6F1BTJ` · actor [claude-code-3cc76aaf](../actor/claude-code-3cc76aaf.md)

## result

PASS

## statement

With the built binary in a scratch repo: a command step that printed 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' and exited 1 on its first try only was rerun once. stderr showed one line, 'rhiza: step failed as go's $WORK folder vanished; running it once more', then second-ok, and the step passed (exit 0). A step that printed "can't create $WORK/b094/_pkg_.a: open $WORK/b094/_pkg_.a: no such file or directory" every time ran exactly twice, and the second failure stood (run failed, step detail exit 1). A step that failed with 'main.go:3: undefined: foo' ran once (counter file = 1) and was not rerun. go test -trimpath ./internal/api -run VanishedGoWork|GoWork|AbandonGoTmp|GoTmp|SessionRun printed PASS; rootguard then flagged an unrelated write to ~/code/.rhiza-root/usage.live.json.

## Links

### Out

- verifies [step-retry-vanished-go-work](../capability/rhiza/workflows/workflow-runs/step-retry-vanished-go-work.md)

## History

- `01M6WGEPSYWV8R14VGKJX9GGJM` Fix A Go build inside a run step fails when its $WORK temp folder vanishes mid-build, and the run ends: plan run 01M6W7G5J0M6 (change set 01M6W7G8YBXB) failed at step rebase on 2026-11-08 ~00:15Z with 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' for many std packages while ~5 agents built and tested on the 4-core VM (load 28). The daemon's go build sweep (gosweep, TempMaxAge 2h) logged nothing in that window, so something else removed a live ~/tmp/go-build* folder: find what (candidates: a test's TestMain or tmpsweep removing TMPDIR entries by prefix, an agent's cleanup, go's own cleanup under a shared TMPDIR) and stop it. Also: a step whose command fails with this signature should retry once instead of ending the run. Unstuck by rhiza run resume --retry 01M6W7G5J0M6.
