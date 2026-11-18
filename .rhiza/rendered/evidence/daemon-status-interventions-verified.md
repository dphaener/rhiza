# daemon-status-interventions-verified

evidence · `01M7R7M2MXKQJD21JPPHM1KRXN` · actor [claude-code-6642bdf8](../actor/claude-code-6642bdf8.md)

## result

PASS

## statement

Verifier at 441ef77d7 re-ran every Checks: command, capped, exit 0 each: go test ./cmd/rhiza -run 'TestScripts/(daemon-.*|health-projects|inbox-projects|watchdog-health|hook-commands|usage-hold|stage-autonomy)$' (incl. new daemon-status-interventions.txtar: no finished runs gives 0/0/0 with no error in --json and text; live run row intervention_count 2 and ', 2 interventions'; three finished runs intervened/clean/clean by finish time with IDs in reverse order give finished 3, finished_clean 2, clean_streak 2; a newer intervened run resets the streak to 0); describe-golden -update leaves no diff and describe-golden passes; go test ./internal/api -run 'Daemon|Health|Web|ReadOnly|Describe|Schema|Docs'; go vet ./internal/api; rhiza ship generate (0 changed); rhiza render. Diff read: internal/api/daemon.go only adds fields (finished, finished_clean, clean_streak at top level and per project, intervention_count per run) and the text line; the streak sorts ended runs by EndedAt newest first (later ID on a tie); existing fields and scheduling untouched (two existing scripts only gain the new intervention_count key). A scratch multi-project script (root with repo and second, removed afterwards) showed per-project finished 2/1/1 and 1/1/1, a live row intervention_count 1, top-level 0s, and the text lines '  Finished runs: N, M with no interventions, clean streak K' under each project.

## Links

### Out

- verifies [daemon-status-interventions](../capability/rhiza/workflows/workflow-runs/daemon-status-interventions.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted
