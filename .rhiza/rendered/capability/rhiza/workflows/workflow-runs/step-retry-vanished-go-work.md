# step-retry-vanished-go-work

capability · `01M6WGXA821QDFHM5T4K3G3P2Q` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A table test of the pure signature match accepts the 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' and "can't create $WORK/b094/_pkg_.a: open $WORK/b094/_pkg_.a: no such file or directory" lines and rejects other failures (a compile error, a missing file outside $WORK); a stepRunner.Command test in internal/api with a command that fails with the signature once then passes returns exit 0 after one rerun and one stderr line, one that fails with it twice returns the second failure, and one that fails otherwise is not rerun

## statement

A command step that exits non-zero with go's vanished-$WORK signature (open $WORK/<bNNN>/...: no such file or directory) is rerun once by the step runner, which writes one stderr line saying so, and the second result stands

## Links

### Out

- decided-by [go-work-vanish-fix-shape](../../../../decision/go-work-vanish-fix-shape.md)

### In

- verifies [step-retry-vanished-go-work-verified](../../../../evidence/step-retry-vanished-go-work-verified.md)
- verifies [step-retry-vanished-go-work-qa-1](../../../../evidence/step-retry-vanished-go-work-qa-1.md)
- verifies [pr-388-fix-a-go-build-inside-a-run](../../../../evidence/pr-388-fix-a-go-build-inside-a-run.md)

## History

- `01M6WGEPSYWV8R14VGKJX9GGJM` Fix A Go build inside a run step fails when its $WORK temp folder vanishes mid-build, and the run ends: plan run 01M6W7G5J0M6 (change set 01M6W7G8YBXB) failed at step rebase on 2026-11-08 ~00:15Z with 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' for many std packages while ~5 agents built and tested on the 4-core VM (load 28). The daemon's go build sweep (gosweep, TempMaxAge 2h) logged nothing in that window, so something else removed a live ~/tmp/go-build* folder: find what (candidates: a test's TestMain or tmpsweep removing TMPDIR entries by prefix, an agent's cleanup, go's own cleanup under a shared TMPDIR) and stop it. Also: a step whose command fails with this signature should retry once instead of ending the run. Unstuck by rhiza run resume --retry 01M6W7G5J0M6.
