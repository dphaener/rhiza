# plan-run-resumes-landed-record-verified

evidence · `01M6H5N4TFS4GNX14ZGCY7CBK9` · actor [claude-code-616127c3](../actor/claude-code-616127c3.md)

## result

PASS

## statement

At 92222406, go test ./cmd/rhiza -run 'TestScripts/(run-resume-landed-plan|run-resume-unaccepted)$' passes: run-resume-landed-plan.txtar starts a plan run that stops at its wait step, lands its record change set (state landed), then rhiza run resume continues it to succeeded; in the same script resume of a ship-change run whose change set is landed and of a plan run whose change set is withdrawn each exits non-zero naming the state and 'next: rhiza run abandon <id>', with the journal cmp-identical; run-resume-unaccepted.txtar still passes. Negative check: with HEAD~1's internal/api/run.go in a scratch copy the new script fails at the resume with 'is landed, no longer accepted'. Diff read: only landed is exempted and only for workflow plan; rejected/withdrawn/superseded still refused. fmt-check, vet and go test ./... pass except TestReadBudget/TestBudget timing tests under full-suite load, which pass rerun alone; self-check fails only on this verifier's own unrendered actor record.

## Links

### Out

- verifies [plan-run-resumes-landed-record](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-gate/plan-run-resumes-landed-record.md)

## History

- `01M6GCQACNNTMW5B57HW8N7K50` Fix: rhiza run resume refuses a plan run whose record change set landed (internal/api/run.go ~792, the Terminal() guard). Every plan run lands its record during the interview, so a plan run stopped at the gate on open questions can never resume. Exempt plan-workflow runs from that guard
