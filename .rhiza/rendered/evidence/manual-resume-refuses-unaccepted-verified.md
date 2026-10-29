# manual-resume-refuses-unaccepted-verified

evidence · `01M62ZQGY4FMY2VH4TGH5R0E7D` · actor [claude-code-45e0d90c](../actor/claude-code-45e0d90c.md)

## result

PASS

## statement

At b1480a3, ran go test ./cmd/rhiza -run TestScripts/run-resume-unaccepted -count=1 -v: it PASSes. The script creates a change set, proposes, accepts and withdraws it, writes an interrupted-run journal (dead pid) for that change set, then rhiza run resume exits 3 with stderr 'is not resumed: its change set ... is withdrawn, no longer accepted' and 'next: rhiza run abandon WITHDRAWN01', and cmp shows the journal byte-identical afterwards. Diff read: internal/api/run.go runRunResume refuses with Conflictf when the run is not ended and the change set state is Terminal (landed/rejected/withdrawn/superseded), before any journal write; draft/proposed/accepted and unknown change sets pass through. go test ./... passes (one earlier run failed on a timing-budget test, a rerun passed); fmt-check and vet clean; self-check fails only on stale rendered actor files from this live run.

## Links

### Out

- verifies [manual-resume-refuses-unaccepted](../capability/rhiza/workflows/workflow-runs/manual-resume-refuses-unaccepted.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
