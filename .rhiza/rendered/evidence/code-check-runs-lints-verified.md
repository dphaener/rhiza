# code-check-runs-lints-verified

evidence · `01M7B7G0PHRSEPYV5CTYDFHYZ6` · actor [claude-code-0e753f53](../actor/claude-code-0e753f53.md)

## result

PASS

## statement

At 3068936 (attempt 2) re-ran the Checks: go test -trimpath ./cmd/rhiza -run 'TestScripts/(review-code-check|describe-golden)$' exit 0; go test -trimpath ./internal/run/... exit 0; go test -trimpath -timeout 30m -skip Budget ./internal/api/... exit 0; rhiza ship generate exit 0 (0 changed); rhiza render exit 0. Read the diff of b01580340 and 3068936: code-check runs StaleRefs, UnappliedDecisions and RemovedCode in the run's worktree (head tree, Leaf: trailers from merge-base, deleted paths) and lists hits as lints in --json, human output and the journal entry, never as causes. review-code-check.txtar asserts the failure case (exit 1, causes name only alpha's Evidence, three lints listed) and the pass (passed true, exit 0, all three hits in --json, human output and run status code_checks for step code). TestCodeCheckLintsRoundTrip (journal Append/Read/State) and TestRunInfoCodeCheckLints (runInfo and human status) cover the round-trip. describe-golden was regenerated. Both earlier findings are addressed, there are no TODOs, and only the leaf's files plus the generated docs/reference changed.

## Links

### Out

- verifies [code-check-runs-lints](../capability/rhiza/workflows/built-in-workflows/built-in-review/code-check-runs-lints.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
