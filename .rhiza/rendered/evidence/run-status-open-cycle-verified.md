# run-status-open-cycle-verified

evidence · `01M6B55CD9ZVEHE3E0T97JY2TT` · actor [claude-code-c3331453](../actor/claude-code-c3331453.md)

## result

PASS

## statement

Verifier at a30953d0: go test ./cmd/rhiza -run TestScripts/run-status passes; its hand-written journal (alpha cycle 1 failed, cycle 2 open) prints '- alpha (step review): open in cycle 2 (cycle 1 failed)' and --json has no check-level outcome. Rebuilt the binary and ran the same journal in a scratch repo: text matches, JSON check object has open_cycle 2 / prior_outcome failed and no outcome, cycle 1 keeps outcome failed. TestChecksSummaryCountsOpenCycleRunning (liveapp) passes: an open-cycle check with empty outcome counts '1 running'. TestOpenCycleReadsOpen (api) passes: runCheckInfo leaves outcome empty, reviewSummary heads '### alpha: open', and a judged cycle 2 restores outcome passed. Read the diff: CheckOutcome/CheckOpen in internal/run/journal.go now feed run status, run summary and checksStatus. make check exits 0.

## Links

### Out

- verifies [run-status-open-cycle](../capability/rhiza/workflows/workflow-runs/run-status/run-status-open-cycle.md)

## History

- `01M6ATASW7ZS2DPMJM2Q12GHM3` A review check whose latest cycle is open reads as open in run status, the live app and run summary
