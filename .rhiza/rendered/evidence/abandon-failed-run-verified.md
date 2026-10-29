# abandon-failed-run-verified

evidence · `01M634192ASNXD0HAKPRZZNW1V` · actor [claude-code-b580cac9](../actor/claude-code-b580cac9.md)

## result

PASS

## statement

At adb7bec I read the diff (refuseEnded in internal/api/run.go now lets a failed, not-yet-abandoned run through; inbox, resume and daemon pick already skip State.Abandoned) and ran go test ./cmd/rhiza -run 'TestScripts/(run-abandon-failed|run-abandon.*|inbox.*|describe-golden)$' -count=1 plus go test ./... -count=1, fmt-check and vet: all pass. run-abandon-failed.txtar fails a fake flaky run, abandons it with exit 0, shows 'failed' with 'abandoned by test' in run status and --json, drops it from rhiza inbox and inbox --json, gets exit 2 naming the run 'was abandoned' from resume and resume --retry, exit 2 'is finished' from a second abandon and from abandoning a succeeded run. make self-check fails only on stale .rhiza/rendered actor files and index from this run's sessions, not on the leaf's diff.

## Links

### Out

- verifies [abandon-failed-run](../capability/rhiza/workflows/workflow-runs/abandon-failed-run.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
