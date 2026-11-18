# run-intervention-message-verified

evidence · `01M7R3FHJ7GVZK2DYZXW803KAB` · actor [claude-code-7100264b](../actor/claude-code-7100264b.md)

## result

PASS

## statement

Verified at 8cb6bf045 (with 1953d890a), capped scope, RHIZA_/SHIP_/REX_ unset: go test -trimpath -count=1 -p 2 ./internal/api -run 'SessionMessage|TypePointer|QueuePointer|APIError|DaemonNudge' ok; go test ./cmd/rhiza -run 'TestScripts/(session-message.*|run-intervention-message|session-run-api-error.*|daemon-nudge.*|run-status-interventions)$' ok (run-intervention-message PASS); TestProcessScripts|TestInProcess|TestProcessListConverted|TestScriptEnv|TestProcessRefusedFailsFast ok; go vet ./internal/api ./cmd/rhiza exit 0; rhiza ship generate exit 0; rhiza render exit 0. The script covers all three acceptance points: an outside message gives intervention_count 1 with actor, session and first line in run status --json; a message with RHIZA_RUN+RHIZA_SESSION of the same run (and with only RHIZA_SESSION, the record naming the run) adds none; a message to a no-run session is still sent (0001.md exists) and journals nothing. The earlier findings are fixed: daemonSend and sendAPIErrorMessage now call deliverSessionMessage (journal=false), so only the session.message op (runSessionMessage) journals. TestSessionMessageRhizaSendsJournalNothing sends through both and finds the journal empty. Delivery code is unchanged apart from the journal hook placed after a successful queue or type. A failed journal write only warns on stderr. typePointer/queuePointer callers (hand-back, amend, lost-session, API-error pointers) never journal. No TODOs. Files outside the leaf's list (main_test.go processScripts entry, daemon_nudge.go, session_run_apierror.go) are the required stdin-feed registration and the fix for the earlier findings.

## Links

### Out

- verifies [run-intervention-message](../capability/rhiza/workflows/workflow-runs/run-intervention-message.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted
