# abandon-failed-run

capability · `01M62S5CM2T21E6M0JQ28RY5Z7` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript fails a fake run, abandons it with exit 0, and shows failed with the abandoned-by detail in run status and --json; rhiza inbox no longer lists the run; resume and resume --retry exit 2 naming it; a second abandon exits 2; abandoning a succeeded run exits 2 as finished

## statement

rhiza run abandon also takes a failed run: it journals the run as abandoned by the actor, so the inbox stops listing it and rhiza run resume refuses it, naming the run; a succeeded run is still refused as finished

## Links

### Out

- depends-on [run-abandon](run-abandon.md)
- derived-from [inbox-stale-interview](../../../../source/inbox-stale-interview.md)

### In

- verifies [abandon-failed-run-verified](../../../../evidence/abandon-failed-run-verified.md)
- verifies [abandon-failed-run-qa-1](../../../../evidence/abandon-failed-run-qa-1.md)
- verifies [abandon-failed-run-verified-2](../../../../evidence/abandon-failed-run-verified-2.md)
- verifies [pr-66-inbox-clears-stale-failures](../../../../evidence/pr-66-inbox-clears-stale-failures.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
