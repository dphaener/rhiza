# docs-reference-workflow-verified

evidence · `01M6ZB4GC9PSHBFY8TBQN71E9K` · actor [claude-code-e15cd5b0](../actor/claude-code-e15cd5b0.md)

## result

PASS

## statement

Verified at ae176775. internal/workflow defines 8 step kinds (agent, wait, command, arbiter, leaves, checks, signoff, workflow: the Kind* constants and kindOutcomes). The new testscript cmd/rhiza/testdata/script/docs-reference-workflow.txtar greps a '## <kind> step' heading for all 8 in the rhiza docs output, plus outcomes and a key, and checks a second run gives the same bytes. TestKindDocsCoverKinds fails when Kinds and kindOutcomes disagree, so a new kind cannot go undocumented. Ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(docs-reference-workflow|docs-reference-cli|describe-golden)$': all PASS. go test ./internal/workflow/... ./internal/docs/...: ok. Built a scratch binary and ran rhiza docs --out: workflow-file.md has file keys, shared keys, a section per kind with keys, tables and outcomes (taken from kindOutcomes), role keys and prompt placeholders. Every toml key the parser declares appears in the doc. gofmt and go vet are clean. The full internal/api package did not finish here: TestReadBudget waits for machine quiet at load 22 on 4 cores. That is contention, not this change; the docs op it adds is covered by the scripts above.

## Links

### Out

- verifies [docs-reference-workflow](../capability/rhiza/documentation/docs-reference-workflow.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
