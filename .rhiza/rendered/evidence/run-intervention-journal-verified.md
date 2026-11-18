# run-intervention-journal-verified

evidence · `01M7R2V8MD8YFTVBDBDXKSJNWC` · actor [claude-code-4c822ffa](../actor/claude-code-4c822ffa.md)

## result

PASS

## statement

Verified b2b166379 at HEAD. Re-ran its Checks: go test -trimpath -p 2 ./internal/run/... (exit 0); go test ./internal/api -run 'RunStatus|RunInfo|Describe|Schema|Docs|Human' (exit 0); go test ./cmd/rhiza -run 'TestScripts/describe-golden$' -update (exit 0, golden unchanged in git); go test -count=1 ./cmd/rhiza -run 'TestScripts/(describe-golden|run-status.*|run-note)$' (exit 0, run-status-interventions PASS); go vet ./internal/run ./internal/api (exit 0); rhiza ship generate (exit 0, 0 changed); rhiza render (exit 0, only this session's actor file). Acceptance: go test -count=1 ./internal/run -run Intervention -v passes TestInterventionFoldsInOrder (two entries in order, count 2, nothing else in the fold changed), TestInterventionNoneCountsZero (count 0), TestInterventionUnknownKind (AppendIntervention refuses unknown kind/empty kind/blank actor and journals nothing; a journalled unknown kind is not counted and the fold warns). run-status-interventions.txtar checks --json interventions list in order and intervention_count 2/0/1, text '- Interventions: 2 (message 1, hand-commit 1)', '- Interventions: 0', and the unknown-kind warning. Diff: must_not held - AppendIntervention has no non-test caller, no other entry type's fold or state/result changed (the new case only appends Interventions or a warning); no TODOs or stubs. store.go and docs/reference/cli.md (generated) are outside the listed files but needed for AppendIntervention and the generated docs.

## Links

### Out

- verifies [run-intervention-journal](../capability/rhiza/workflows/workflow-runs/run-intervention-journal.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted
