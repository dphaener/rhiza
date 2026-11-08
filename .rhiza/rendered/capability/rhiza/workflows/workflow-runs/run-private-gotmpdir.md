# run-private-gotmpdir

capability · `01M6WGXA5QZV54A91F061Y3KNQ` · parent [workflow-runs](../workflow-runs.md)

## acceptance

Unit tests in internal/run (and internal/api where the step runner builds the env) show a command, a wait and an agent step each get GOTMPDIR=<git common dir>/rhiza/runs/<run id>/gotmp, the folder exists when the step's process starts, and it is gone after the run finishes, fails or is abandoned; go test -trimpath ./internal/run/... ./internal/api/... -run for those tests passes

## statement

Every process a run starts (command, wait and agent steps) runs with GOTMPDIR set to the run's own folder <git common dir>/rhiza/runs/<run id>/gotmp, created when a step starts and removed when the run ends, so no sweep of the shared TMPDIR can remove a live go build's $WORK

## Links

### Out

- decided-by [go-work-vanish-fix-shape](../../../../decision/go-work-vanish-fix-shape.md)

### In

- verifies [run-private-gotmpdir-verified](../../../../evidence/run-private-gotmpdir-verified.md)
- verifies [run-private-gotmpdir-qa-1](../../../../evidence/run-private-gotmpdir-qa-1.md)
- verifies [pr-388-fix-a-go-build-inside-a-run](../../../../evidence/pr-388-fix-a-go-build-inside-a-run.md)

## History

- `01M6WGEPSYWV8R14VGKJX9GGJM` Fix A Go build inside a run step fails when its $WORK temp folder vanishes mid-build, and the run ends: plan run 01M6W7G5J0M6 (change set 01M6W7G8YBXB) failed at step rebase on 2026-11-08 ~00:15Z with 'could not import io (open $WORK/b033/_pkg_.a: no such file or directory)' for many std packages while ~5 agents built and tested on the 4-core VM (load 28). The daemon's go build sweep (gosweep, TempMaxAge 2h) logged nothing in that window, so something else removed a live ~/tmp/go-build* folder: find what (candidates: a test's TestMain or tmpsweep removing TMPDIR entries by prefix, an agent's cleanup, go's own cleanup under a shared TMPDIR) and stop it. Also: a step whose command fails with this signature should retry once instead of ending the run. Unstuck by rhiza run resume --retry 01M6W7G5J0M6.
