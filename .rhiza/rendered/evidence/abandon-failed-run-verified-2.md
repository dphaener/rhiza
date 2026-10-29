# abandon-failed-run-verified-2

evidence · `01M634X6YCGHN1KV19524EWNN9` · actor [claude-code-6780cd52](../actor/claude-code-6780cd52.md)

## result

PASS

## statement

Attempt 2 verifier at a6c4d3a. Read the leaf diff (adb7bec: refuseEnded now lets an ended run through only when it failed and is not already abandoned) and the fix a6c4d3a (run-abandon revised to refuse only a succeeded or already abandoned run; abandon-failed-run depends-on run-abandon). Ran go test ./cmd/rhiza -run 'TestScripts/(run-abandon|describe-golden|inbox)' -count=1: run-abandon-failed passes, failing a fake run, abandoning it with exit 0, run status and --json showing failed with 'abandoned by test', rhiza inbox and inbox --json no longer listing it, resume and resume --retry exiting 2 with '<run> was abandoned', a second abandon and abandoning a succeeded run exiting 2 with 'is finished'; run-abandon, describe-golden and inbox-* pass too. go test ./... , fmt-check and vet pass. make self-check fails only on a stale .rhiza/rendered from the QA check's actor record (b87ce24), not on this leaf.

## Links

### Out

- verifies [abandon-failed-run](../capability/rhiza/workflows/workflow-runs/abandon-failed-run.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
