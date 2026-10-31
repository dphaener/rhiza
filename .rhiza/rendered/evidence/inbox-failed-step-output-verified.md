# inbox-failed-step-output-verified

evidence · `01M69QCEQEEAYDG3CQJKMV8PC1` · actor [claude-code-0b17c20a](../actor/claude-code-0b17c20a.md)

## result

PASS

## statement

At 37da5f26, run-abandon-failed.txtar gains a loud workflow whose command step runs 'echo out1; echo boom >&2; exit 3'; the script asserts rhiza inbox prints '    exit 3' and '    boom' indented under the failed item, and inbox --json's detail is 'step "b" ended with outcome "fail"\nexit 3\nout1\nboom'. go test ./cmd/rhiza -run 'TestScripts/(run-abandon-failed|inbox-list)$' -count=1 passed; fmt-check, vet and go test ./... passed. Diff read: internal/api/inbox.go failedDetail appends the last failing visit's multi-line step detail to the run detail, and Human indents lines after the first. make self-check flagged only the unrendered actor file of this verifier session, not the leaf.

## Links

### Out

- verifies [inbox-failed-step-output](../capability/rhiza/orchestration/inbox/inbox-list/inbox-failed-step-output.md)

## History

- `01M686RRCJSQ57M82FENWT2HNM` When a run's command step fails, its journal detail keeps the tail of the step's output, shown by run status and the inbox, so an unattended failure says why
