# modellint-unapplied-decisions-verified

evidence · `01M7B18XBT9Z2CQJEBHHQTF8EF` · actor [claude-code-e5e5f755](../actor/claude-code-e5e5f755.md)

## result

PASS

## statement

At 931bbb81e re-ran the Checks: block: go test -trimpath ./internal/modellint/... exit 0 (TestUnappliedDecisions -v: all 7 subtests pass), rhiza ship generate exit 0 (0 changed), rhiza render exit 0 (no tracked change). go vet and gofmt clean. Read the diff: internal/modellint/decisions.go adds UnappliedDecisions over the Decisions cs writes (writtenNodes, factored out of refs.go with StaleRefs unchanged in behaviour); decisions_test.go table covers no decided-by link (reported), answers-only (not), landed capability applies (not), leaf without Leaf: commit (reported), leaf committed by handle and by ID (not), and the failure case of another change set's unapplied Decision (not). No change outside internal/modellint besides the staged actor record and rendered files; no stubs or TODOs; tests use in-memory models only.

## Links

### Out

- verifies [modellint-unapplied-decisions](../capability/rhiza/workflows/built-in-workflows/built-in-review/modellint-unapplied-decisions.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
