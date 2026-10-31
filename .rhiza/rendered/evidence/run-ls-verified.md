# run-ls-verified

evidence · `01M69884E9JHSJSP447G7M6V37` · actor [claude-code-36726f6c](../actor/claude-code-36726f6c.md)

## result

PASS

## statement

At ccba290 ran go test ./cmd/rhiza -run 'TestScripts/(run-ls|describe-golden)$' -count=1: both PASS. run-ls.txtar (RHIZA_TEST_NOW/SEED set by Setup) creates three runs' journals (ok, bad, ok workflows), checks run ls --json is empty with no runs, lists titles Third, Second, First in that order, shows live, state succeeded/failed and an RFC 3339 started_at, keeps only the newest with --limit 1, rejects --limit -1 and 1001, and checks the human line '<ID> succeeded ok since <time>, change set <ID>: Third run'. describe-golden.txtar carries run.ls's request/response schema. Read internal/api/run_ls.go: it reverses run.Store.List (sorted by start time then ID), applies the limit, and takes state and live from the shared runInfo. make check passed on a rerun; the first run failed only on the known read-timing budget, and go test ./... passed on its own.

## Links

### Out

- verifies [run-ls](../capability/rhiza/workflows/workflow-runs/run-ls.md)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
